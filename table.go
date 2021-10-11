package wowlua

import (
	"errors"
	"strings"
)

var (
	// ErrNotFound indicates the requested node wasn't found
	ErrNotFound = errors.New("Node not found")
	// ErrNotTable indicates the node wasn't a table
	ErrNotTable = errors.New("Node not a table")
	// ErrWrongType indicates that the node was the wrong type
	ErrWrongType = errors.New("Node is wrong type")
)

// Table is the top-level data structure returned by parsing. The table
// consists of table entries. Each entry has a key and a value, each of type
// Node.
type Table struct {
	entries []*tableEntry
}

// NewTable creates a new, empty table
func NewTable() *Table {
	return &Table{}
}

// String returns a string representation of the table
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

// HasKeyByString checks whether the table has an entry with the provided key
// string. This does not parse numeric strings to compare to numeric keys.
func (t *Table) HasKeyByString(s string) bool {
	return t.HasKey(&Node{nType: NodeTypeString, value: s})
}

// HasKey checks whether the table has an entry with a key equal to the
// provided key node.
func (t *Table) HasKey(k *Node) bool {
	return t.getEntry(k) != nil
}

// Set an entry in table with the provided key-value pair. If an entry exists
// with that key it is overwritten. If not it is added.
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

// GetStringByString looks for an entry in the table with a string key equal to
// the provided string and a value of type string. It returns an error if no
// entry has a matching key or the matching entry node is not a string.
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

// GetFloatByString looks for an entry in the table with a string key equal to
// the provided string and a value of type Number. It returns an error if no
// entry has a matching key or the matching entry node is not a Number.
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

// GetByString looks for an entry with a key node of type string matching the
// provided string value. If no matching node is found nil is returned.
func (t *Table) GetByString(s string) *Node {
	return t.Get(NewNode(NodeTypeString, s))
}

// Get retrieves a node from the table with a key equal to the provided key
func (t *Table) Get(k *Node) *Node {
	e := t.getEntry(k)
	if e == nil {
		return nil
	}
	return e.value
}

// GetStringByPath walks through nested tables to find a node matching the
// path. All keys in the path must be strings.
func (t *Table) GetStringPath(path ...string) (*Table, *Node, error) {
	pathNodes := make([]*Node, len(path))
	for i, elem := range path {
		pathNodes[i] = NewNode(NodeTypeString, elem)
	}
	return t.GetPath(pathNodes...)
}

// GetPath walks through nested tables to find a node matching the
// path.
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

// Keys returns all the keys in the table as a slice.
func (t *Table) Keys() []*Node {
	keys := make([]*Node, t.Len())
	for i, e := range t.entries {
		keys[i] = e.key
	}
	return keys
}

// AddIndexed appends the node to the table. The key for the new node is the
// current number of entries. This value is returned.
func (t *Table) AddIndexed(n *Node) int {
	i := t.Len()
	k := &Node{nType: NodeTypeNumber, value: i}
	t.Set(k, n)
	return i
}

// Equals returns whether this table is equivalent to another table.
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

// Len returns the number of entries in the table
func (t *Table) Len() int {
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
