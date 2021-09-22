package model

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gitran-com/gitran-server/config"
	"github.com/gitran-com/gitran-server/constant"
	"github.com/gitran-com/gitran-server/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"gorm.io/gorm"
)

//ProjCfg means project config
type ProjCfg struct {
	ID           int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Status       int        `json:"status" gorm:"type:tinyint;"`
	SrcBr        string     `json:"src_branch" gorm:"type:varchar(32);notNull"`
	TrnBr        string     `json:"trn_branch" gorm:"type:varchar(32);notNull"`
	PullGap      uint16     `json:"pull_gap" gorm:"index;notNull"`
	PushGap      uint16     `json:"push_gap" gorm:"index;notNull"`
	PullStatus   int        `json:"pull_status" gorm:"notNull"`
	PushStatus   int        `json:"push_status" gorm:"notNull"`
	LastPullAt   *time.Time `json:"last_pull_at"`
	LastPushAt   *time.Time `json:"last_push_at"`
	SrcRegs      []string   `json:"src_files" gorm:"-"`
	SrcRegsBytes []byte     `json:"-" gorm:"column:src_regs"`
	TrnReg       string     `json:"trn_file" gorm:"column:trn_reg"`
	IgnRegs      []string   `json:"ign_regs" gorm:"-"`
	IgnRegsBytes []byte     `json:"-" gorm:"column:ign_regs"`
	Extra        ExtraCfg   `json:"extra" gorm:"-"`
	ExtraBytes   []byte     `json:"-" gorm:"column:extra"`
}

type ExtraCfg struct {
	XML *XMLCfg `json:"xml"`
}

type XMLCfg struct {
	OmitTags []string `json:"omit_tags" gorm:"-"`
}

//TableName return table name
func (*ProjCfg) TableName() string {
	return config.DB.TablePrefix + "project_configs"
}

//Write writes project config to DB
func (cfg *ProjCfg) Write() {
	db.Save(cfg)
}

func (cfg *ProjCfg) AfterFind(tx *gorm.DB) error {
	json.Unmarshal(cfg.SrcRegsBytes, &cfg.SrcRegs)
	json.Unmarshal(cfg.IgnRegsBytes, &cfg.IgnRegs)
	json.Unmarshal(cfg.ExtraBytes, &cfg.Extra)
	return nil
}

//NewProjCfg creates a new cfg for project
func NewProjCfg(cfg *ProjCfg) (*ProjCfg, error) {
	if result := db.Create(cfg); result.Error != nil {
		return nil, result.Error
	}
	return cfg, nil
}

func (cfg *ProjCfg) UpdateProjCfg(req *UpdateProjCfgRequest) error {
	if err := db.Model(cfg).Updates(req.Map()).Error; err != nil {
		return err
	}
	go cfg.Process()
	return nil
}

//SSOT <== DB
func (cfg *ProjCfg) UpdateStatus(stat int) {
	cfg.Status = stat
	cfg.Write()
}

//GetProjCfgByID get a project config by config id
func GetProjCfgByID(id int64) *ProjCfg {
	var pc ProjCfg
	res := db.Where("id=?", id).First(&pc)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return &pc
}

//UpdateProjCfgPullStatus update a project cfg pull status
func UpdateProjCfgPullStatus(cfg *ProjCfg, stat int) {
	db.Model(cfg).Select("pull_status").Updates(map[string]interface{}{"pull_status": stat})
}

//UpdateProjCfgPushStatus update a project cfg push status
func UpdateProjCfgPushStatus(cfg *ProjCfg, stat int) {
	db.Model(cfg).Select("push_status").Updates(map[string]interface{}{"push_status": stat})
}

//ListSyncProjCfg list all project cfg that should be sync
func ListSyncProjCfg() []ProjCfg {
	var cfg []ProjCfg
	db.Where("pull_gap!=0 OR push_gap!=0").Find(&cfg)
	return cfg
}

func (cfg *ProjCfg) Process() {
	var (
		proj = GetProjByID(cfg.ID)
		stat = constant.ProjStatReady
	)
	lk := ProjMutexMap.Lock(proj.ID)
	defer lk.Unlock()
	defer cfg.UpdateStatus(stat)
	cfg.UpdateStatus(constant.ProjStatProcessing)
	repo, err := git.PlainOpen(proj.Path)
	if err != nil {
		log.Errorf("ProjCfg.Process error when git.PlainOpen(%s): %+v", proj.Path, err.Error())
		return
	}
	wt, _ := repo.Worktree()
	srcBr := "refs/heads/" + cfg.SrcBr
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(srcBr),
	})
	if err != nil {
		log.Errorf("ProjCfg.Process() error when git.Checkout(%s): %+v", srcBr, err.Error())
		stat = constant.ProjStatFailed
		return
	}
	files := util.ListMultiMatchFiles(proj.Path, cfg.SrcRegs, cfg.IgnRegs)
	TxnUpdateFilesSents(proj, files)
}

func GenTrnFilesFromSrcFiles(src []string, trn string, lang *Language, proj *Project) (res []string) {
	reg := regexp.MustCompile(`\$.*?\$`)
	for _, str := range src {
		tmp := reg.ReplaceAllStringFunc(trn, func(s string) string {
			switch s {
			case `$dir_name$`, `$dir$`:
				return filepath.Dir(str)
			case `$file_name$`, `$name$`:
				return filepath.Base(str)
			case `$base_name$`, `$base$`:
				return util.FilenameNoExt(filepath.Base(str))
			case `$ext_name$`, `$ext$`:
				return filepath.Ext(str)
			case `$language$`, `$lang$`:
				return lang.ISO
			case `$code$`:
				return lang.Code
			case `$code2$`:
				return lang.Code2
			default:
				return ""
			}
		})
		res = append(res, tmp)
	}
	return res
}

func GenMultiTrnFilesFromSrcFiles(srcRegs []string, trnReg string, ignores []string, proj *Project) ([]string, map[string][]string) {
	var (
		multiTrnFiles = make(map[string][]string)
	)
	srcFiles := util.ListMultiMatchFiles(proj.Path, srcRegs, ignores)
	for _, lang := range proj.TranslateLanguages {
		trn := GenTrnFilesFromSrcFiles(srcFiles, trnReg, &lang, proj)
		multiTrnFiles[lang.Code] = trn
	}
	return srcFiles, multiTrnFiles
}

func TxnUpdateFilesSents(proj *Project, files []string) {
	db.Transaction(func(tx *gorm.DB) error {
		tx.Model(&ProjFile{}).Where("proj_id=?", proj.ID).Updates(map[string]interface{}{"valid": false})
		wg := sync.WaitGroup{}
		wg.Add(len(files))
		for _, file := range files {
			pf := &ProjFile{
				ProjID: proj.ID,
				Path:   file,
				Valid:  true,
			}
			res := tx.Where("proj_id=? AND path=?", proj.ID, file).First(pf)
			if res.Error != nil {
				pf.Content = string(readFile(proj.Path, file))
				tx.Create(pf)
			} else {
				content := string(readFile(proj.Path, file))
				if pf.Content != content {
					pf.Content = content
				}
				pf.Valid = true
				tx.Save(pf)
			}
			go pf.TxnProcess(&wg, tx, proj)
		}
		wg.Wait()
		return nil
	})
}
