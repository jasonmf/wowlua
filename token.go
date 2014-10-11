package wowlua

import (
    "bufio"
    "fmt"
    "jlog"
    "strings"
    "unicode"
)

const (
    StateTokenNone = iota
    StateTokenBareHyphen
    StateTokenFindNewline
    StateTokenString
    StateTokenEscapedChar
    StateTokenNumber
    StateTokenIdentifier
)

const (
    TokenTypeStartTable = iota
    TokenTypeEndTable
    TokenTypeStartKey
    TokenTypeEndKey
    TokenTypeEquals
    TokenTypeComma
    TokenTypeIgnore
    TokenTypeString
    TokenTypeNumber
    TokenTypeIdentifier
)

var (
    tokenTypeStrings = map[int]string {
        TokenTypeStartTable: "Start Table",
        TokenTypeEndTable:   "End Table",
        TokenTypeStartKey:   "Start Key",
        TokenTypeEndKey:     "End Key",
        TokenTypeEquals:     "Equals",
        TokenTypeComma:      "Comma",
        TokenTypeIgnore:     "Ignore",
        TokenTypeString:     "String",
        TokenTypeNumber:     "Number",
        TokenTypeIdentifier: "Identifier",
    }
    tokenStartTable = &Token{ Type: TokenTypeStartTable }
    tokenEndTable = &Token{ Type: TokenTypeEndTable }
    tokenStartKey = &Token{ Type: TokenTypeStartKey }
    tokenEndKey = &Token{ Type: TokenTypeEndKey }
    tokenEquals = &Token{ Type: TokenTypeEquals }
    tokenComma = &Token{ Type: TokenTypeComma }
)

type Token struct {
    Type    int
    Value   string
}

func NewToken(nType int, v string) *Token {
    t := &Token{
        Type:   nType,
        Value:  v,
    }
    return t
}

func (t Token) String() string {
    s := tokenTypeStrings[t.Type]
    if s == "" {
        s = "INVALID TOKEN"
    }
    if t.Value != "" {
        s += " (" + t.Value + ")"
    }
    return s
}

type Tokenizer struct {
    stack       []rune
    state       int
    scanner     *bufio.Scanner
    callback    func(*Token)
}

func NewTokenizer(s string, c func(*Token)) (*Tokenizer) {
    if c == nil {
        c = func(tok *Token) { fmt.Println(*tok) }
    }
    t := &Tokenizer{
        stack:  []rune{},
        state:  StateTokenNone,
        scanner: bufio.NewScanner(strings.NewReader(s)),
        callback: c,
    }
    t.scanner.Split(bufio.ScanRunes)
    return t
}

func (t *Tokenizer) Push(r rune) {
    t.stack = append(t.stack, r)
}

func (t *Tokenizer) Pop(pType int) {
    t.Emit(NewToken(pType, string(t.stack)))
    t.stack = t.stack[:0]
}

func (t *Tokenizer) Emit(tok *Token) {
    t.callback(tok)
}

func (t *Tokenizer) SetStateToken(state int) {
    t.state = state
}

func (t *Tokenizer) Tokenize() {
    for t.scanner.Scan() {
        r := rune(t.scanner.Text()[0])
        switch t.state {
            case StateTokenNone:
                switch {
                    case r == '{':
                        t.Emit(tokenStartTable)
                    case r == '}':
                        t.Emit(tokenEndTable)
                    case r == '[':
                        t.Emit(tokenStartKey)
                    case r == ']':
                        t.Emit(tokenEndKey)
                    case r == '-':
                        t.Push(r)
                        t.SetStateToken(StateTokenBareHyphen)
                    case r == '=':
                        t.Emit(tokenEquals)
                    case r == '"':
                        t.SetStateToken(StateTokenString)
                    case r == ',':
                        t.Emit(tokenComma)
                    case unicode.IsDigit(r):
                        t.Push(r)
                        t.SetStateToken(StateTokenNumber)
                    case unicode.IsSpace(r):
                        /* Do Nothing */
                    case unicode.IsLetter(r):
                        t.Push(r)
                        t.SetStateToken(StateTokenIdentifier)
                    default:
                        jlog.Fatalf("FECK None")
                }
            case StateTokenBareHyphen:
                switch {
                    case r == '-':
                        t.Pop(TokenTypeIgnore)
                        t.SetStateToken(StateTokenFindNewline)
                    case unicode.IsDigit(r):
                        t.Push(r)
                        t.SetStateToken(StateTokenNumber)
                    default:
                        jlog.Fatalf("FECK")
                }
            case StateTokenFindNewline:
                if r == '\n' {
                    t.SetStateToken(StateTokenNone)
                }
            case StateTokenString:
                switch r {
                    case '\\':
                        t.SetStateToken(StateTokenEscapedChar)
                    case '"':
                        t.Pop(TokenTypeString)
                        t.SetStateToken(StateTokenNone)
                    default:
                        t.Push(r)

                }
            case StateTokenEscapedChar:
                t.Push(r)
                t.SetStateToken(StateTokenString)
            case StateTokenNumber:
                switch {
                case r == ',':
                    t.Pop(TokenTypeNumber)
                    t.SetStateToken(StateTokenNone)
                    t.Emit(tokenComma)
                case unicode.IsDigit(r):
                    t.Push(r)
                case r == '.':
                    t.Push(r)
                case r == ']':
                    t.Pop(TokenTypeNumber)
                    t.SetStateToken(StateTokenNone)
                    t.Emit(tokenEndKey)
                default:
                    t.Pop(TokenTypeNumber)
                    t.SetStateToken(StateTokenNone)
                }
            case StateTokenIdentifier:
                switch {
                case r == ',':
                    t.Pop(TokenTypeIdentifier)
                    t.SetStateToken(StateTokenNone)
                    t.Emit(tokenComma)
                case unicode.IsSpace(r):
                    t.Pop(TokenTypeIdentifier)
                    t.SetStateToken(StateTokenNone)
                default:
                    t.Push(r)
                }
        }
    }
}
