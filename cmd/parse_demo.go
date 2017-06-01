package main

import (
    "bytes"
    "flag"
    "io"
    "log"
    "os"

    "jlog"
    "wowlua"
)

var (
    fInPath = flag.String("in", "", "Input file")
    fDebug = flag.Bool("debug", false, "Debug logging")
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
    end := p.Finish()
    log.Println("=================================================================")
    log.Println(end)
}
