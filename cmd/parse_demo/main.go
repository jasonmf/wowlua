package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/jasonmf/wowlua"
	"github.com/jasonmf/wowlua/cmd"
)

var (
	fInPath = flag.String("in", "", "Input file")
)

func main() {
	flag.Parse()
	b, err := ioutil.ReadFile(*fInPath)
	cmd.FatalIfError(err, "reading file")

	table, err := wowlua.ParseLua(string(b))
	cmd.FatalIfError(err, "parsing")
	log.Println("=================================================================")
	log.Println(table)
}
