package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"

	"github.com/jasonmf/wowlua"
	"github.com/jasonmf/wowlua/cmd"
)

var (
	fInPath = flag.String("in", "", "Input file")
	fPath   = flag.String("path", "", "Nav path")
)

func main() {
	flag.Parse()
	b, err := ioutil.ReadFile(*fInPath)
	cmd.FatalIfError(err, "reading input file")

	table, err := wowlua.ParseLua(string(b))
	cmd.FatalIfError(err, "parsing")

	pathElem := strings.Split(*fPath, "/")
	_, node, err := table.GetStringPath(pathElem...)
	cmd.FatalIfError(err, "getting element")
	log.Println(node)
}
