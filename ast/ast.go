package ast

import (
	"bytes"

	"github.com/shafik23/ys/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of every AST our parser produces.
// Every valid YS program is a series of statements.
type Program struct {
	Statements []Statement
}

////////////////////////////////////////////////////////////////

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 { // if the program has any statements
		return p.Statements[0].TokenLiteral() // return the literal of the first statement
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements { // iterate over the statements
		out.WriteString(s.String()) // append the string representation of each statement to the buffer
	}

	return out.String()
}

////////////////////////////////////////////////////////////////

// LetStatement represents a let statement.
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier // the name of the variable
	Value Expression  // the value the variable is bound to
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ") // append the literal of the token.LET token
	out.WriteString(ls.Name.String())        // append the string representation of the name of the variable
	out.WriteString(" = ")                   // append the assignment operator

	if ls.Value != nil { // if the variable is bound to a value
		out.WriteString(ls.Value.String()) // append the string representation of the value
	}

	out.WriteString(";") // append the statement terminator

	return out.String()
}

////////////////////////////////////////////////////////////////

// Identifier represents an identifier.
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string      // the identifier's value
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

////////////////////////////////////////////////////////////////

type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression  // the value the return statement returns
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ") // append the literal of the token.RETURN token

	if rs.ReturnValue != nil { // if the return statement returns a value
		out.WriteString(rs.ReturnValue.String()) // append the string representation of the value
	}

	out.WriteString(";") // append the statement terminator

	return out.String()
}

////////////////////////////////////////////////////////////////

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil { // if the expression is not nil
		return es.Expression.String() // return its string representation
	}

	return ""
}

////////////////////////////////////////////////////////////////

type IntegerLiteral struct {
	Token token.Token // the token.INT token
	Value int64       // the integer literal's value
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

////////////////////////////////////////////////////////////////
