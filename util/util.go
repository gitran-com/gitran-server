package util

import (
	"gopkg.in/wzru/sentences.v1"
	"gopkg.in/wzru/sentences.v1/data"
)

var (
	tokenizer *sentences.DefaultSentenceTokenizer
)

func init() {
	b, _ := data.Asset("data/english.json")
	training, _ := sentences.LoadTraining(b)
	tokenizer = sentences.NewSentenceTokenizer(training)
}

func Tokenize(text string) []*sentences.Sentence {
	return tokenizer.Tokenize(text)
}
