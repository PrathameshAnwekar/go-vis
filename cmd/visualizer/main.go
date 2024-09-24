package main

import (
	"fmt"
	"os"

	"github.com/PrathameshAnwekar/go-vis/internal/fsutil"
	"github.com/PrathameshAnwekar/go-vis/internal/log"
	"github.com/PrathameshAnwekar/go-vis/internal/parser"
)

func main() {
	log.Init()

	args := os.Args
	if len(args) < 2 {
		log.I("Please provide a project path.")
		return
	}
	projectRoot := args[1]
	log.I("Parsing project:", args[1])

	fileList, err := fsutil.GetGoFiles(projectRoot)
	if err != nil {
		fmt.Println(err)
	}

	log.I(fileList)
	parser.ParseGoProject(fileList)
}
