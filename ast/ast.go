package ast

import "github.com/shafik23/ys/token"

type Node interface {
	TokenLiteral() string
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

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 { // if the program has any statements
		return p.Statements[0].TokenLiteral() // return the literal of the first statement
	} else {
		return ""
	}
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

////////////////////////////////////////////////////////////////

type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression  // the value the return statement returns
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

////////////////////////////////////////////////////////////////
