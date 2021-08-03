package util

import (
	"encoding/xml"
	"fmt"
)

func ProcessXML(file string) []string {
	var out map[string]interface{}
	if err := xml.Unmarshal([]byte(file), out); err != nil {
		fmt.Printf("ProcessXML error: %+v\n", err.Error())
		return nil
	}
	fmt.Printf("out=%+v\n", out)
	sentences := tokenizer.Tokenize(file)
	for _, s := range sentences {
		fmt.Printf("sen='%s'\n", s.Text)
	}
	return nil
}
