package main

import (
    "bytes"
    "flag"
    "io"
    "log"
    "os"
    "strings"

    "jlog"
    "wowlua"
)

var (
    fInPath = flag.String("in", "", "Input file")
    fDebug = flag.Bool("debug", false, "Debug logging")
    fPath = flag.String("path", "", "Nav path")
)

func logFatal(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    flag.Parse()
    if *fDebug {
        jlog.SetLevel(jlog.LogLevelDebug)
    }
    infh, err := os.Open(*fInPath)
    logFatal(err)
    defer infh.Close()
    buf := &bytes.Buffer{}
    io.Copy(buf, infh)
    p := wowlua.NewParser()
    t := wowlua.NewTokenizer(string(buf.Bytes()), p.Next)
    t.Tokenize()
    end, err := p.Finish()
    jlog.FatalIfError(err)
    pathElem := strings.Split(*fPath, "/")
    _, node, err := end.GetStringPath(pathElem...)
    jlog.FatalIfError(err)
    jlog.Println(node)
}
