package object

import (
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/shafik23/ys/ast"
)

const (
	NULL_OBJ         = "NULL"
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

//////////////////////////////////////////////////

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

//////////////////////////////////////////////////

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

//////////////////////////////////////////////////

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

//////////////////////////////////////////////////

type Null struct{}

func (n *Null) Inspect() string { return "null" }

func (n *Null) Type() ObjectType { return NULL_OBJ }

//////////////////////////////////////////////////

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

//////////////////////////////////////////////////

type Error struct {
	Message string
}

func (e *Error) Inspect() string { return "ERROR: " + e.Message }

func (e *Error) Type() ObjectType { return ERROR_OBJ }

//////////////////////////////////////////////////

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Inspect() string {
	return fmt.Sprintf("fn(%s) {\n%s\n}", f.Parameters, f.Body.String())
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

//////////////////////////////////////////////////

type String struct {
	Value string
}

func (s *String) Inspect() string { return s.Value }

func (s *String) Type() ObjectType { return STRING_OBJ }

//////////////////////////////////////////////////

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Inspect() string { return "builtin function" }

func (b *Builtin) Type() ObjectType { return FUNCTION_OBJ }

//////////////////////////////////////////////////

type Array struct {
	Elements []Object
}

func (ao *Array) Inspect() string {
	var out string
	for i, e := range ao.Elements {
		if i == 0 {
			out += "["
		}
		out += e.Inspect()
		if i != len(ao.Elements)-1 {
			out += ", "
		} else {
			out += "]"
		}
	}
	return out
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }

//////////////////////////////////////////////////

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

//////////////////////////////////////////////////

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Inspect() string {
	parts := make([]string, 0, len(h.Pairs))

	for _, pair := range h.Pairs {
		parts = append(parts, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	return fmt.Sprintf("{%s}", strings.Join(parts, ", "))
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }

//////////////////////////////////////////////////

type Hashable interface {
	HashKey() HashKey
}
