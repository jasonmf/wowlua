package wowlua

import (
    "fmt"
    "jlog"
    "math"
    "reflect"
    "strconv"
)

const (
    NodeTypeString = iota
    NodeTypeIdentifier
    NodeTypeNumber
    NodeTypeBoolean
    NodeTypeTable
    NodeTypeTableEntry
)

var (
    NaN = math.NaN()
)

type Node struct {
    nType   int
    value   interface{}
}

func NewNode(nType int, v interface{}) *Node {
    n := &Node{
        nType:  nType,
        value:  v,
    }
    return n
}

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
    case NodeTypeBoolean:
        return fmt.Sprint("BOOLEAN: ", n.GetBoolean())
    case NodeTypeTableEntry:
        if v, ok := n.value.(*tableEntry); ok {
            return fmt.Sprint("TABLEENTRY: ", v)
        }
        return ">> NodeTypeTableEntry failed type assertion to *tableEntry <<" + reflect.ValueOf(n.value).String()
    }
    return "NOPE (" + fmt.Sprint(n.nType) + ")"
}

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
            jlog.Errorf("Unhandled type comparison: %v", n.nType)
    }
    return false
}

func (n *Node) GetType() int {
    return n.nType
}

func (n *Node) GetString() string {
    if n.nType != NodeTypeString && n.nType != NodeTypeIdentifier {
        return ""
    }
    if s, ok := n.value.(string); ok {
        return s
    }
    return ""
}

func (n *Node) GetFloat64() float64 {
    if n.nType != NodeTypeNumber {
        jlog.Debugf("Node.GetFloat64 called on wrong node type: %v", n)
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
            jlog.Errorf("Error parsing float64 %q", s)
            return NaN
        }
        return f
    }
    return NaN
}

func (n *Node) GetBoolean() bool {
    if n.nType != NodeTypeBoolean {
        return false
    }
    if b, ok := n.value.(bool); ok {
        return b
    }
    return false
}

func (n *Node) GetTable() *Table {
    if n.nType != NodeTypeTable {
        return nil
    }
    if t, ok := n.value.(*Table); ok {
        return t
    }
    return nil
}
