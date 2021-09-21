package model

import (
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gitran-com/gitran-server/config"
	"github.com/gitran-com/gitran-server/util"
	"gorm.io/gorm"
)

type ProjFile struct {
	ID      int64  `gorm:"primaryKey"`
	ProjID  int64  `gorm:"index"`
	Path    string `gorm:"index"`
	Valid   bool   `gorm:"index"`
	SentCnt int    `gorm:""`
	WordCnt int    `gorm:""`
}

//TableName return table name
func (*ProjFile) TableName() string {
	return config.DB.TablePrefix + "project_files"
}

func (file *ProjFile) Write() {
	db.Save(file)
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
	res := db.Where("proj_id=? AND file=?", proj_id, file).First(&pf)
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

func (file *ProjFile) Process(wg *sync.WaitGroup) {
	defer wg.Done()
	proj := GetProjByID(file.ProjID)
	abs := path.Join(proj.Path, file.Path)
	data, err := ioutil.ReadFile(abs)
	if err != nil {
		return
	}
	ext := filepath.Ext(file.Path)
	var res []string
	switch ext {
	case ".xml":
		res = util.ProcessXML(data)
	default:
		res = util.ProcessTXT(data)
	}
	SetAllSentsInvalid(file.ID)
	file.SentCnt = len(res)
	file.WordCnt = 0
	for _, str := range res {
		file.WordCnt += len(strings.Fields(str))
		MustGetValidSent(file.ProjID, file.ID, str)
	}
	file.Write()
}

func SetAllFilesInvalid(proj_id int64) {
	db.Model(&ProjFile{}).Where("proj_id=?", proj_id).Updates(map[string]interface{}{"valid": false})
}
