package wowlua

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

const (
	// NodeTypeString is a node containing a string
	NodeTypeString = iota
	// NodeTypeIdentifier is a node containing an identifier
	NodeTypeIdentifier
	// NodeTypeNumber is a node containing a number
	NodeTypeNumber
	// NodeTypeBool is a node containing a bool
	NodeTypeBool
	// NodeTypeTable is a node containing a table
	NodeTypeTable
	// NodeTypeTableEntry is a node containing a table entry
	NodeTypeTableEntry
)

var (
	// NaN is not a number
	NaN = math.NaN()
)

// Node contains a parsed value. This value can be a scalar or a table
type Node struct {
	nType int
	value interface{}
}

// NewNode creates a node of the given type with the provided value. It assumes
// the specified type is appropriate for the value.
func NewNode(nType int, v interface{}) *Node {
	n := &Node{
		nType: nType,
		value: v,
	}
	return n
}

// String returns a string representation of the value
func (n *Node) String() string {
	switch n.nType {
	case NodeTypeString:
		return "STRING: " + n.GetString()
	case NodeTypeIdentifier:
		return "IDENTIFIER: " + n.GetString()
	case NodeTypeTable:
		return fmt.Sprint("TABLE: ", n.GetTable())
	case NodeTypeNumber:
		return fmt.Sprint("NUMBER: ", n.GetFloat64())
	case NodeTypeBool:
		return fmt.Sprint("BOOL: ", n.GetBool())
	case NodeTypeTableEntry:
		if v, ok := n.value.(*tableEntry); ok {
			return fmt.Sprint("TABLEENTRY: ", v)
		}
		return ">> NodeTypeTableEntry failed type assertion to *tableEntry <<" + reflect.ValueOf(n.value).String()
	}
	return "NOPE (" + fmt.Sprint(n.nType) + ")"
}

// Equals returns whether this node is equal to another. For this to be true
// the types and values must match. Not all types support equality testing.
func (n *Node) Equals(o *Node) bool {
	if n == o {
		return true
	}
	if n.nType != o.nType {
		return false
	}
	switch n.nType {
	case NodeTypeString:
		return n.GetString() == o.GetString()
	case NodeTypeIdentifier:
		return n.GetString() == o.GetString()
	case NodeTypeNumber:
		return n.GetFloat64() == o.GetFloat64()
	default:
		logger.Errorf("Unhandled type comparison: %v", n.nType)
	}
	return false
}

// GetType returns the type the node holds
func (n *Node) GetType() int {
	return n.nType
}

// GetString returns the underlying string value of the node if it is a string
// or identifier type and empty string if not
func (n *Node) GetString() string {
	if n.nType != NodeTypeString && n.nType != NodeTypeIdentifier {
		return ""
	}
	if s, ok := n.value.(string); ok {
		return s
	}
	return ""
}

// GetFloat64 returns the underlying value of the node if it is numeric and NaN
// if not
func (n *Node) GetFloat64() float64 {
	if n.nType != NodeTypeNumber {
		logger.Debugf("Node.GetFloat64 called on wrong node type: %v", n)
		return NaN
	}
	if nInt, ok := n.value.(int); ok {
		return float64(nInt)
	}
	if nf64, ok := n.value.(float64); ok {
		return nf64
	}
	if s, ok := n.value.(string); ok {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			logger.Errorf("Error parsing float64 %q", s)
			return NaN
		}
		return f
	}
	return NaN
}

// GetBool returns the value of the node if it's a bool and false if it is not.
// You should check that the node type is a bool first.
func (n *Node) GetBool() bool {
	if n.nType != NodeTypeBool {
		return false
	}
	if b, ok := n.value.(bool); ok {
		return b
	}
	return false
}

// GetTable returns the underlying table if the node is of type table and nil
// if not
func (n *Node) GetTable() *Table {
	if n.nType != NodeTypeTable {
		return nil
	}
	if t, ok := n.value.(*Table); ok {
		return t
	}
	return nil
}
