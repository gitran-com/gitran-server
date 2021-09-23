package model

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gitran-com/gitran-server/config"
	"gorm.io/gorm"
)

type ProjFile struct {
	ID      int64  `json:"id" gorm:"primaryKey"`
	ProjID  int64  `json:"proj_id" gorm:"index"`
	Path    string `json:"path" gorm:"index"`
	Valid   bool   `json:"-" gorm:"index"`
	SentCnt int    `json:"sent_cnt" gorm:""`
	WordCnt int    `json:"word_cnt" gorm:""`
	Content string `json:"-" gorm:"type:text"`
}

//TableName return table name
func (*ProjFile) TableName() string {
	return config.DB.TablePrefix + "project_files"
}

func (file *ProjFile) Write() {
	db.Save(file)
}

func (file *ProjFile) ReadContent() []byte {
	proj := GetProjByID(file.ProjID)
	if proj == nil {
		return nil
	}
	return readFile(proj.Path, file.Path)
}

//NewProjFile creates a new file for project
func NewProjFile(file *ProjFile) (*ProjFile, error) {
	if result := db.Create(file); result.Error != nil {
		return nil, result.Error
	}
	return file, nil
}

func GetProjFileByPath(proj_id int64, file string) *ProjFile {
	var pf ProjFile
	res := db.Where("proj_id=? AND path=?", proj_id, file).First(&pf)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return &pf
}

func MustGetValidProjFile(id int64, file string) *ProjFile {
	pf := &ProjFile{}
	res := db.Where("proj_id=? AND path=?", id, file).First(pf)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		pf, _ = NewProjFile(&ProjFile{
			ProjID: id,
			Path:   file,
			Valid:  true,
		})
		return pf
	}
	pf.Valid = true
	pf.Write()
	return pf
}

func (file *ProjFile) TxnProcess(wg *sync.WaitGroup, tx *gorm.DB, cfg *ProjCfg) {
	defer wg.Done()
	data := []byte(file.Content)
	ext := filepath.Ext(file.Path)
	var (
		sens []string
		offs []int
	)
	switch ext {
	case ".xml":
		sens, offs = ProcessXML(cfg, data)
	default:
		sens, offs = ProcessTXT(data)
	}
	// SetAllSentsInvalid(file.ID)
	tx.Model(&ProjFile{}).Where("proj_id=?", cfg.ID).Updates(map[string]interface{}{"valid": false})
	file.SentCnt = len(sens)
	file.WordCnt = 0
	for i, str := range sens {
		file.WordCnt += len(strings.Fields(str))
		hash := fmt.Sprintf("%x", md5.Sum([]byte(str)))
		sent := &Sentence{
			ProjID:  cfg.ID,
			FileID:  file.ID,
			Offset:  offs[i],
			Valid:   true,
			Content: str,
			MD5:     hash,
		}
		res := tx.Where("file_id=? AND md5=?", file.ID, hash).First(sent)
		if res.Error != nil {
			tx.Create(sent)
		} else {
			sent.Valid = true
			tx.Save(sent)
		}
	}
	tx.Save(file)
}

func SetAllFilesInvalid(proj_id int64) {
	db.Model(&ProjFile{}).Where("proj_id=?", proj_id).Updates(map[string]interface{}{"valid": false})
}

func readFile(root string, rel string) []byte {
	abs := path.Join(root, rel)
	data, _ := ioutil.ReadFile(abs)
	return data
}

func ListProjFiles(proj_id int64) []ProjFile {
	var pfs []ProjFile
	db.Where("proj_id=?", proj_id).Find(&pfs)
	return pfs
}

func GetProjFileByID(proj_id int64, file_id int64) *ProjFile {
	var pf ProjFile
	res := db.Where("proj_id=? AND id=?", proj_id, file_id).First(&pf)
	if res.Error != nil {
		return nil
	}
	return &pf
}
