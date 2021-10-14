package model

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/gitran-com/gitran-server/config"
	"github.com/gitran-com/gitran-server/util"
)

type Sentence struct {
	ID           int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjID       int64  `json:"proj_id" gorm:"index"`
	FileID       int64  `json:"file_id" gorm:"index"`
	PinnedTranID int64  `json:"pinned_tran_id"`
	Offset       int    `json:"offset" gorm:"index"`
	Valid        bool   `json:"-" gorm:"index"`
	Locked       bool   `json:"locked" gorm:""`
	MD5          string `json:"-" gorm:"index;type:char(32)"`
	Content      string `json:"content" gorm:"type:text"`
}

//TableName return table name
func (*Sentence) TableName() string {
	return config.DB.TablePrefix + "sentences"
}

func (sent *Sentence) Write() {
	db.Save(sent)
}

//NewSentence creates a new file for project
func NewSentence(sent *Sentence) (*Sentence, error) {
	if result := db.Create(sent); result.Error != nil {
		return nil, result.Error
	}
	return sent, nil
}

func SetAllSentsInvalid(file_id int64) {
	db.Model(&Sentence{}).Where("file_id=?", file_id).Updates(map[string]interface{}{"valid": false})
}

var (
	xmlTagReg = regexp.MustCompile(`(<.[^(><.)]+>)`)
)

func commonProcess(str string) string {
	return strings.TrimSpace(str)
}

func ProcessXML(cfg *ProjCfg, data []byte) ([]string, []int) {
	strs := []string{}
	offs := []int{}
	tags := xmlTagReg.FindAll(data, -1)
	allTags := make(map[string]bool)
	ignoredTags := make(map[string]bool)
	for _, tag := range cfg.Extra.XML.IgnoredTags {
		ignoredTags[tag] = true
		ignoredTags[tag+"/"] = true
	}
	for _, tag := range tags {
		if bytes.HasPrefix(tag, []byte("</")) {
			continue
		}
		allTags[string(tag[1:len(tag)-1])] = true
	}
	for tag := range allTags {
		if ignoredTags[tag] {
			continue
		}
		flag := true
		reg := regexp.MustCompile(`<` + tag + `>[\s\S]*?<\/` + tag + `>`)
		res := reg.FindAllIndex(data, -1)
		for _, pair := range res {
			start, end := pair[0]+len(tag)+2, pair[1]-len(tag)-3
			content := string(data[start:end])
			childs := xmlTagReg.FindAll([]byte(content), -1)
			for _, child := range childs {
				if bytes.HasPrefix(child, []byte("</")) {
					continue
				}
				end := len(child) - 1
				if bytes.HasSuffix(child, []byte("/>")) {
					end--
				}
				if !ignoredTags[string(child[1:end])] {
					flag = false
					break
				}
			}
			if flag {
				sens := util.Tokenize(content)
				for _, s := range sens {
					sen := commonProcess(s.Text)
					if sen == "" {
						continue
					}
					strs = append(strs, sen)
					offs = append(offs, start+s.Start)
				}
			}
		}
	}
	return strs, offs
}

func ProcessTXT(data []byte) ([]string, []int) {
	strs := []string{}
	offs := []int{}
	sens := util.Tokenize(string(data))
	for _, s := range sens {
		sen := commonProcess(s.Text)
		strs = append(strs, sen)
		offs = append(offs, s.Start)
	}
	return strs, offs
}

func ListValidSents(file_id int64) []Sentence {
	var sens []Sentence
	db.Where("file_id=? AND valid=?", file_id, true).Find(&sens)
	return sens
}

func GetSentByID(sent_id int64) *Sentence {
	var sent Sentence
	res := db.Where("id=?", sent_id).First(&sent)
	if res.Error != nil {
		return nil
	}
	return &sent
}

func (sent *Sentence) PinTran(tran *Translation) {
	sent.PinnedTranID = tran.ID
	db.Save(sent)
}

func (sent *Sentence) UnpinTran() {
	sent.PinnedTranID = 0
	db.Save(sent)
}
