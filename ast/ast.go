package ast

import (
	"bytes"
	"strings"

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

type PrefixExpression struct {
	Token    token.Token // the prefix operator, e.g. !
	Operator string      // the prefix operator, e.g. !
	Right    Expression  // the right-hand side expression
}

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(") // append the opening parenthesis
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")") // append the closing parenthesis

	return out.String()
}

////////////////////////////////////////////////////////////////

type InfixExpression struct {
	Token    token.Token // the operator, e.g. +
	Left     Expression  // the left-hand side expression
	Operator string      // the operator, e.g. +
	Right    Expression  // the right-hand side expression
}

func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(") // append the opening parenthesis
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")") // append the closing parenthesis

	return out.String()
}

////////////////////////////////////////////////////////////////

type Boolean struct {
	Token token.Token // the token.TRUE or token.FALSE token
	Value bool        // the boolean literal's value
}

func (b *Boolean) expressionNode() {}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}

////////////////////////////////////////////////////////////////

type IfExpression struct {
	Token       token.Token // the token.IF token
	Condition   Expression  // the condition expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if") // append the if keyword
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil { // if the if expression has an alternative block
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

////////////////////////////////////////////////////////////////

type BlockStatement struct {
	Token      token.Token // the token.LBRACE token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements { // iterate over the statements
		out.WriteString(s.String()) // append the string representation of each statement to the buffer
	}

	return out.String()
}

////////////////////////////////////////////////////////////////

type FunctionLiteral struct {
	Token      token.Token // the token.FUNCTION token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{} // create a slice of strings

	for _, p := range fl.Parameters { // iterate over the parameters
		params = append(params, p.String()) // append the string representation of each parameter to the slice
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", ")) // join the parameters with a comma and a space
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

////////////////////////////////////////////////////////////////

type CallExpression struct {
	Token     token.Token // the token.LPAREN token
	Function  Expression  // the function expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{} // create a slice of strings

	for _, a := range ce.Arguments { // iterate over the arguments
		args = append(args, a.String()) // append the string representation of each argument to the slice
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", ")) // join the arguments with a comma and a space
	out.WriteString(")")

	return out.String()
}

////////////////////////////////////////////////////////////////

type StringLiteral struct {
	Token token.Token // the token.STRING token
	Value string      // the string literal's value
}

func (sl *StringLiteral) expressionNode() {}

func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}
