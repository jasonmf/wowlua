package wowlua

import (
	"errors"
	"fmt"
	"jlog"
)

type Parser struct {
	stack []*Node
}

func NewParser() *Parser {
	p := &Parser{
		stack: make([]*Node, 0),
	}
	p.startTable() // Top level table
	return p
}

func (p *Parser) Push(n *Node) {
	p.stack = append(p.stack, n)
}

func (p *Parser) Pop() *Node {
	sLen := len(p.stack)
	if sLen == 0 {
		return nil
	}
	n := p.stack[sLen-1]
	p.stack = p.stack[0 : sLen-1]
	return n
}

func (p *Parser) Peek() *Node {
	sLen := len(p.stack)
	if sLen == 0 {
		return nil
	}
	return p.stack[sLen-1]
}

func (p *Parser) startTable() {
	t := NewTable()
	n := &Node{nType: NodeTypeTable, value: t}
	p.Push(n)
}

func (p *Parser) startKey() {
	n := &Node{nType: NodeTypeTableEntry, value: &tableEntry{}}
	p.Push(n)
}

func tokenToNode(t *Token) (*Node, error) {
	var n *Node
	switch t.Type {
	case TokenTypeIdentifier:
		switch t.Value {
		case "true":
			n = NewNode(NodeTypeBoolean, true)
		case "false":
			n = NewNode(NodeTypeBoolean, false)
		default:
			n = NewNode(NodeTypeString, t.Value)
		}
	case TokenTypeString:
		n = NewNode(NodeTypeString, t.Value)
	case TokenTypeNumber:
		jlog.Debugf("TokenToNode: %q", t)
		n = NewNode(NodeTypeNumber, t.Value)
	default:
		return nil, fmt.Errorf("Can't convert this token: %q", t)
	}
	return n, nil
}

func (p *Parser) Next(t *Token) error {
	jlog.Debugf("Parsing Token: %v", t)
	if t.Type == TokenTypeIgnore {
		return nil
	}
	switch t.Type {
	case TokenTypeIdentifier:
		top := p.Peek()
		switch top.nType {
		case NodeTypeTableEntry:
			n, err := tokenToNode(t)
			if err != nil {
				return err
			}
			p.Push(n)
		default:
			p.Next(tokenStartKey)
			p.Next(NewToken(TokenTypeString, t.Value))
			p.Next(tokenEndKey)
		}
	case TokenTypeEquals:
		top := p.Peek()
		if top.nType != NodeTypeTableEntry {
			return p.bailout("Found equals with non table entry.")
		}
		if top.value == nil {
			return p.bailout("Found equals with nil table entry.")
		}
		if _, ok := top.value.(*tableEntry); !ok {
			return p.bailout("Key node contains non-table entry")
		}
	case TokenTypeStartTable:
		p.startTable()
	case TokenTypeEndTable:
		top := p.Peek()
		if top.nType != NodeTypeTable {
			return p.bailout("Ending table while not in table!")
		}
	case TokenTypeStartKey:
		if p.Peek().nType != NodeTypeTable {
			return p.bailout(fmt.Sprintf("Found start of key under %q", p.Peek()))
		}
		p.startKey()
	case TokenTypeEndKey:
		key := p.Pop()
		jlog.Debugf("Popped Key: %v", key)
		top := p.Peek()
		if top.nType != NodeTypeTableEntry {
			return p.bailout("Found end key without table entry.")
		}
		if top.value == nil {
			return p.bailout("Found end key on nil entry.")
		}
		if e, ok := top.value.(*tableEntry); ok {
			if e.key != nil {
				return p.bailout("Found end key with key already set.")
			}
			e.key = key
		} else {
			return p.bailout("Found end key on non-table-entry")
		}
	case TokenTypeComma:
		v := p.Pop()
		top := p.Peek()
		switch top.nType {
		case NodeTypeTableEntry:
			if e, ok := top.value.(*tableEntry); ok {
				p.Pop()
				top = p.Peek()
				if top.nType != NodeTypeTable {
					return p.bailout("Key not in table!")
				}
				if table, ok := top.value.(*Table); ok {
					table.Set(e.key, v)
				} else {
					return p.bailout("Expected table node on top of stack.")
				}
			} else {
				return p.bailout("TableEntryValue not tableEntry!")
			}
		case NodeTypeTable:
			top.GetTable().AddIndexed(v)
		default:
			return p.bailout("Comma found outside table, table key")
		}
	case TokenTypeString:
		return p.handleValueToken(t)
	case TokenTypeNumber:
		return p.handleValueToken(t)
	default:
		jlog.Debugf("Unhandled token type: %q", t)
	}
	return nil
}

func (p *Parser) handleValueToken(t *Token) error {
	top := p.Peek()
	if top.nType != NodeTypeTable && top.nType != NodeTypeTableEntry {
		return p.bailout(fmt.Sprintf("Found value %q outside of table/key!", t))
	}
	n, err := tokenToNode(t)
	if err != nil {
		return err
	}
	p.Push(n)
	return nil
}

func (p *Parser) unwind() {
	for n := p.Pop(); n != nil; n = p.Pop() {
		jlog.Debugf("Stack: %v", n)
	}
}

func (p *Parser) dumpstack() {
	for i := len(p.stack) - 1; i > -1; i-- {
		jlog.Debugf("Stack Dump: %v", p.stack[i])
	}
}

func (p *Parser) bailout(msg string) error {
	jlog.Errorf("BAILING OUT!")
	p.unwind()
	return errors.New(msg)
}

func (p *Parser) Finish() (*Table, error) {
	for n := p.Pop(); n != nil; n = p.Pop() {
		top := p.Peek()
		switch n.nType {
		case NodeTypeTableEntry:
			if e, ok := n.value.(*tableEntry); ok {
				if top.nType != NodeTypeTable {
					return nil, p.bailout("Key not in table!")
				}
				if table, ok := top.value.(*Table); ok {
					table.Set(e.key, e.value)
				} else {
					return nil, p.bailout("Expected table node on top of stack.")
				}
			} else {
				jlog.Debugf("Stack Popped: ", n)
				return nil, p.bailout("TableEntryValue not tableEntry!")
			}
		default:
			if top != nil {
				switch top.nType {
				case NodeTypeTableEntry:
					if e, ok := top.value.(*tableEntry); ok {
						e.value = n
					} else {
						return nil, p.bailout("Top is marked tableEntry but isn't a *tableEntry")
					}
				default:
					return nil, p.bailout("SHIT BROKE SON")
				}
			} else {
				jlog.Debugf("--> TOP is NIL <--")
				if t, ok := n.value.(*Table); ok {
					return t, nil
				}
			}
		}

	}

	err := errors.New("Final result not a table!")
	jlog.Errorf(err.Error())
	return nil, err
}

func ParseLua(data string) (*Table, error) {
	p := NewParser()
	t := NewTokenizer(data, p.Next)
	if err := t.Tokenize(); err != nil {
		return nil, err
	}
	return p.Finish()
}
