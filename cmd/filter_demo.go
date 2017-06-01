package main

import (
    "bytes"
    "flag"
    "io"
    "log"
    "os"

    "wowlua"
)

var (
    fInPath = flag.String("in", "", "Input file")
)

func logFatal(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    flag.Parse()
    infh, err := os.Open(*fInPath)
    logFatal(err)
    defer infh.Close()
    buf := &bytes.Buffer{}
    io.Copy(buf, infh)
    /*
    f := wowlua.NewFilterStream(string(buf.Bytes()))
    for s := f.Next(); s != wowlua.NoData; s = f.Next() {
        fmt.Print(s)
    }
    */
    t := wowlua.NewTokenizer(string(buf.Bytes()), nil)
    t.Tokenize()
}
