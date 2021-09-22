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

func ProcessXML(data []byte) []string {
	doc, _ := xmlquery.Parse(bytes.NewReader(data))
	res := xmlquery.Find(doc, "//*")
	strs := []string{}
	for _, node := range res {
		text := strings.TrimSpace(node.InnerText())
		sens := tokenizer.Tokenize(text)
		for _, s := range sens {
			strs = append(strs, commonProcess(s.Text))
		}
	}
	return strs
}

func ProcessTXT(data []byte) []string {
	strs := []string{}
	sens := tokenizer.Tokenize(string(data))
	for _, s := range sens {
		strs = append(strs, commonProcess(s.Text))
	}
	return strs
}
