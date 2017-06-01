package wowlua

import (
	"bufio"
	"errors"
	"fmt"
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
	StateTokenInvalid
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
	tokenTypeStrings = map[int]string{
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
	tokenStartTable = &Token{Type: TokenTypeStartTable}
	tokenEndTable   = &Token{Type: TokenTypeEndTable}
	tokenStartKey   = &Token{Type: TokenTypeStartKey}
	tokenEndKey     = &Token{Type: TokenTypeEndKey}
	tokenEquals     = &Token{Type: TokenTypeEquals}
	tokenComma      = &Token{Type: TokenTypeComma}
)

// A Token is a symbol identified by the tokenizer.
type Token struct {
	Type  int
	Value string
}

// Create a new token.
func NewToken(nType int, v string) *Token {
	t := &Token{
		Type:  nType,
		Value: v,
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

// A Tokenizer processes a string, yielding Tokens.
type Tokenizer struct {
	buffer   []rune
	state    int
	scanner  *bufio.Scanner
	callback func(*Token) error
	err      error
}

// NewTokenizer creates a new Tokenizer to process the supplied string. It will
// provide each token to the supplied callback function. If the callback
// returns an error, tokenization will stop.
func NewTokenizer(s string, c func(*Token) error) *Tokenizer {
	if c == nil {
		c = func(tok *Token) error { fmt.Println(*tok); return nil }
	}
	t := &Tokenizer{
		buffer:   []rune{},
		state:    StateTokenNone,
		scanner:  bufio.NewScanner(strings.NewReader(s)),
		callback: c,
	}
	t.scanner.Split(bufio.ScanRunes)
	return t
}

// Add a rune to the buffer.
func (t *Tokenizer) Buffer(r rune) {
	t.buffer = append(t.buffer, r)
}

// Create a new token from the buffer of the specified type, Emit() the token,
// then clear the buffer.
func (t *Tokenizer) Send(pType int) {
	t.Emit(NewToken(pType, string(t.buffer)))
	t.buffer = t.buffer[:0]
}

// If there's been no previous error, send a token to the callback and capture
// any returned error.
func (t *Tokenizer) Emit(tok *Token) {
	if t.err != nil {
		return
	}
	t.err = t.callback(tok)
}

// Set the current state. If the specified state is invalid, nothing happens.
func (t *Tokenizer) SetStateToken(state int) {
	if state < StateTokenNone || state >= StateTokenInvalid {
		return
	}
	t.state = state
}

// Process in the input stream until it's finished or an error is encountered.
func (t *Tokenizer) Tokenize() error {
	for t.scanner.Scan() {
		if t.err != nil {
			return t.err
		}
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
				t.Buffer(r)
				t.SetStateToken(StateTokenBareHyphen)
			case r == '=':
				t.Emit(tokenEquals)
			case r == '"':
				t.SetStateToken(StateTokenString)
			case r == ',':
				t.Emit(tokenComma)
			case unicode.IsDigit(r):
				t.Buffer(r)
				t.SetStateToken(StateTokenNumber)
			case unicode.IsSpace(r):
				/* Do Nothing */
			case unicode.IsLetter(r):
				t.Buffer(r)
				t.SetStateToken(StateTokenIdentifier)
			default:
				return errors.New("FECK None")
			}
		case StateTokenBareHyphen:
			switch {
			case r == '-':
				t.Send(TokenTypeIgnore)
				t.SetStateToken(StateTokenFindNewline)
			case unicode.IsDigit(r):
				t.Buffer(r)
				t.SetStateToken(StateTokenNumber)
			default:
				return errors.New("FECK")
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
				t.Send(TokenTypeString)
				t.SetStateToken(StateTokenNone)
			default:
				t.Buffer(r)

			}
		case StateTokenEscapedChar:
			t.Buffer(r)
			t.SetStateToken(StateTokenString)
		case StateTokenNumber:
			switch {
			case r == ',':
				t.Send(TokenTypeNumber)
				t.SetStateToken(StateTokenNone)
				t.Emit(tokenComma)
			case unicode.IsDigit(r):
				t.Buffer(r)
			case r == '.':
				t.Buffer(r)
			case r == ']':
				t.Send(TokenTypeNumber)
				t.SetStateToken(StateTokenNone)
				t.Emit(tokenEndKey)
			default:
				t.Send(TokenTypeNumber)
				t.SetStateToken(StateTokenNone)
			}
		case StateTokenIdentifier:
			switch {
			case r == ',':
				t.Send(TokenTypeIdentifier)
				t.SetStateToken(StateTokenNone)
				t.Emit(tokenComma)
			case unicode.IsSpace(r):
				t.Send(TokenTypeIdentifier)
				t.SetStateToken(StateTokenNone)
			default:
				t.Buffer(r)
			}
		}
	}
	return nil
}
