package wowlua

import (
	"errors"
	"strings"
)

var (
	ErrNotFound  = errors.New("Node not found")
	ErrNotTable  = errors.New("Node not a table")
	ErrWrongType = errors.New("Node is wrong type")
)

type Table struct {
	entries []*tableEntry
}

func NewTable() *Table {
	return &Table{}
}

func (t *Table) String() string {
	entryStrs := make([]string, len(t.entries)+1)
	entryStrs[0] = "{"
	for i, entry := range t.entries {
		entryStrs[i+1] = entry.String()
	}
	if len(entryStrs) > 3 {
		entryStrs[3] = "..."
		entryStrs = entryStrs[:4]
	}
	return strings.Join(entryStrs, "\n") + "}"
}

func (t *Table) HasKeyByString(s string) bool {
	return t.HasKey(&Node{nType: NodeTypeString, value: s})
}

func (t *Table) HasKey(k *Node) bool {
	return t.getEntry(k) != nil
}

func (t *Table) Set(k, v *Node) {
	e := t.getEntry(k)
	if e == nil {
		e = &tableEntry{key: k}
		t.entries = append(t.entries, e)
	}
	e.value = v
}

func (t *Table) getEntry(k *Node) *tableEntry {
	for _, e := range t.entries {
		if e.key.Equals(k) {
			return e
		}
	}
	return nil
}

func (t *Table) GetStringByString(s string) (string, error) {
	n := t.GetByString(s)
	if n == nil {
		return "", ErrNotFound
	}
	if n.GetType() != NodeTypeString {
		return "", ErrWrongType
	}
	return n.GetString(), nil
}

func (t *Table) GetFloat64ByString(s string) (float64, error) {
	n := t.GetByString(s)
	if n == nil {
		return NaN, ErrNotFound
	}
	if n.GetType() != NodeTypeNumber {
		return NaN, ErrWrongType
	}
	return n.GetFloat64(), nil
}

func (t *Table) GetByString(s string) *Node {
	return t.Get(NewNode(NodeTypeString, s))
}

func (t *Table) Get(k *Node) *Node {
	e := t.getEntry(k)
	if e == nil {
		return nil
	}
	return e.value
}

func (t *Table) GetStringPath(path ...string) (*Table, *Node, error) {
	pathNodes := make([]*Node, len(path))
	for i, elem := range path {
		pathNodes[i] = NewNode(NodeTypeString, elem)
	}
	return t.GetPath(pathNodes...)
}

func (t *Table) GetPath(path ...*Node) (*Table, *Node, error) {
	if !t.HasKey(path[0]) {
		return t, nil, ErrNotFound
	}
	sub := t.Get(path[0])
	if len(path) == 1 {
		return t, sub, nil
	}
	subT := sub.GetTable()
	if subT == nil {
		return t, sub, ErrNotTable
	}
	return subT.GetPath(path[1:]...)
}

func (t *Table) Keys() []*Node {
	keys := make([]*Node, t.Length())
	for i, e := range t.entries {
		keys[i] = e.key
	}
	return keys
}

func (t *Table) AddIndexed(n *Node) int {
	i := t.Length()
	k := &Node{nType: NodeTypeNumber, value: i}
	t.Set(k, n)
	return i
}

func (t *Table) Equals(o *Table) bool {
	if t == o {
		return true
	}
	if len(t.entries) != len(o.entries) {
		return false
	}
	// TODO(jason): Make this faster than O(n^2), tho for small cases it's fine
	for _, te := range t.entries {
		oe := o.getEntry(te.key)
		if oe == nil {
			return false
		}
		if !te.value.Equals(oe.value) {
			return false
		}
	}
	return true
}

func (t *Table) Length() int {
	return len(t.entries)
}

type tableEntry struct {
	key   *Node
	value *Node
}

func (e *tableEntry) String() string {
	keyStr := "(nil key)"
	valueStr := "(nil value)"
	if e.key != nil {
		keyStr = e.key.String()
	}
	if e.value != nil {
		valueStr = e.value.String()
	}
	return "{" + keyStr + ", " + valueStr + "}"
}
