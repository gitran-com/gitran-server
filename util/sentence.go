package util

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/antchfx/xmlquery"
)

var space = regexp.MustCompile(`\s+`)

func commonProcess(str string) string {
	return strings.TrimSpace(str)
}

func ProcessXML(data []byte) ([]string, []int) {
	doc, _ := xmlquery.Parse(bytes.NewReader(data))
	res := xmlquery.Find(doc, "//*")
	strs := []string{}
	offs := []int{}
	start := 0
	node := res[0]
	text := strings.TrimSpace(node.InnerText())
	sens := tokenizer.Tokenize(text)
	for _, s := range sens {
		sen := commonProcess(s.Text)
		off := bytes.Index(data[start:], []byte(sen))
		strs = append(strs, sen)
		offs = append(offs, off)
		start = off + 1
	}
	return strs, offs
}

func ProcessTXT(data []byte) ([]string, []int) {
	strs := []string{}
	offs := []int{}
	start := 0
	sens := tokenizer.Tokenize(string(data))
	for _, s := range sens {
		sen := commonProcess(s.Text)
		off := bytes.Index(data[start:], []byte(sen))
		strs = append(strs, sen)
		offs = append(offs, off)
		start = off + 1
	}
	return strs, offs
}
