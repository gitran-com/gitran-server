package main

import (
	"flag"

	"github.com/WangZhengru/gitran-be/config"
)

func main() {
	flag.Parse()
	err := config.Init()
	if err != nil {
		return
	}

}
