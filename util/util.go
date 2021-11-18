package util

import (
	"github.com/neurosnap/sentences"
	"github.com/neurosnap/sentences/data"
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
