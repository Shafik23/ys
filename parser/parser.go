package parser

import (
	"fmt"

	"github.com/shafik23/ys/ast"
	"github.com/shafik23/ys/lexer"
	"github.com/shafik23/ys/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

type prefixParseFn func() ast.Expression
type infixParseFn func(ast.Expression) ast.Expression

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn // map of prefix parse functions
	infixParseFns  map[token.TokenType]infixParseFn  // map of infix parse functions

	errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// Register prefix parse functions
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn) // initialize the map
	p.registerPrefix(token.IDENT, p.parseIdentifier)           // register the parseIdentifier function

	// Read two tokens, so curToken and peekToken are both set.
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}              // create a new AST root node
	program.Statements = []ast.Statement{} // initialize the Statements field to an empty slice

	for !p.curTokenIs(token.EOF) { // loop until we reach the end of the input
		stmt := p.parseStatement() // parse a statement

		if stmt != nil {
			program.Statements = append(program.Statements, stmt) // append it to the Statements field
		}

		p.nextToken() // advance the tokens
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type { // check the type of the current token
	case token.LET: // if it is a let statement
		return p.parseLetStatement() // parse it
	case token.RETURN: // if it is a return statement
		return p.parseReturnStatement() // parse it
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken} // create a new let statement node and set its token field

	if !p.expectPeek(token.IDENT) { // if the next token is not an identifier
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal} // set the name field to an identifier node

	if !p.expectPeek(token.ASSIGN) { // if the next token is not an assignment operator
		return nil
	}

	// Skip the expressions until we encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) { // loop until we reach the end of the statement
		p.nextToken() // advance the tokens
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t // check the type of the current token
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t // check the type of the next token
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) { // check the type of the next token
		p.nextToken() // if it matches, advance the tokens
		return true
	} else {
		p.peekErrors(t) // otherwise, add an error to the errors slice
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekErrors(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type) // create an error message
	p.errors = append(p.errors, msg)                                                        // append it to the errors slice
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken} // create a new return statement node and set its token field

	p.nextToken() // advance the tokens

	// TODO!
	// Skip the expressions until we encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) { // loop until we reach the end of the statement
		p.nextToken() // advance the tokens
	}

	return stmt
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn // register a prefix parse function for a given token type
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn // register an infix parse function for a given token type
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken} // create a new expression statement node and set its token field

	stmt.Expression = p.parseExpression(LOWEST) // parse the expression

	if p.peekTokenIs(token.SEMICOLON) { // if the next token is a semicolon
		p.nextToken() // advance the tokens
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type] // get the prefix parse function for the current token type

	if prefix == nil { // if there is no prefix parse function
		return nil // return nil
	}

	leftExp := prefix() // parse the prefix expression

	return leftExp
}
