package parser

import (
	"fmt"
	"strconv"

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

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

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
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)

	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)

	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

	// Register infix parse functions
	p.infixParseFns = make(map[token.TokenType]infixParseFn) // initialize the map
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	p.registerInfix(token.LPAREN, p.parseCallExpression)

	// Read two tokens, so curToken and peekToken are both set.
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}          // create a new array literal node and set its token field
	array.Elements = p.parseExpressionList(token.RBRACKET) // parse the array elements
	return array
}

func (p *Parser) parseExpressionList(t token.TokenType) []ast.Expression {
	list := []ast.Expression{} // initialize the list to an empty slice

	if p.peekTokenIs(t) { // if the next token is a right bracket
		p.nextToken() // advance the tokens
		return list   // return the empty slice
	}

	p.nextToken()                                  // advance the tokens
	list = append(list, p.parseExpression(LOWEST)) // parse the first element

	for p.peekTokenIs(token.COMMA) { // loop until we reach the end of the elements
		p.nextToken()                                  // advance the tokens
		p.nextToken()                                  // advance the tokens
		list = append(list, p.parseExpression(LOWEST)) // parse the next element
	}

	if !p.expectPeek(t) { // if the next token is not a right bracket
		return nil
	}

	return list
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	// create a new call expression node and set its token and function fields
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments() // parse the call arguments
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{} // initialize the arguments slice to an empty slice

	if p.peekTokenIs(token.RPAREN) { // if the next token is a right parenthesis
		p.nextToken() // advance the tokens
		return args   // return the empty slice
	}

	p.nextToken()                                  // advance the tokens
	args = append(args, p.parseExpression(LOWEST)) // parse the first argument

	for p.peekTokenIs(token.COMMA) { // loop until we reach the end of the arguments
		p.nextToken()                                  // advance the tokens
		p.nextToken()                                  // advance the tokens
		args = append(args, p.parseExpression(LOWEST)) // parse the next argument
	}

	if !p.expectPeek(token.RPAREN) { // if the next token is not a right parenthesis
		return nil
	}

	return args
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken} // create a new function literal node and set its token field

	if !p.expectPeek(token.LPAREN) { // if the next token is not a left parenthesis
		return nil
	}

	lit.Parameters = p.parseFunctionParameters() // parse the function parameters

	if !p.expectPeek(token.LBRACE) { // if the next token is not a left brace
		return nil
	}

	lit.Body = p.parseBlockStatement() // parse the function body

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{} // initialize the identifiers slice to an empty slice

	if p.peekTokenIs(token.RPAREN) { // if the next token is a right parenthesis
		p.nextToken() // advance the tokens
		return identifiers
	}

	p.nextToken() // advance the tokens

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal} // create a new identifier node and set its token and value fields

	identifiers = append(identifiers, ident) // append it to the identifiers slice

	for p.peekTokenIs(token.COMMA) { // loop until we reach the end of the parameters
		p.nextToken() // advance the tokens
		p.nextToken() // advance the tokens

		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal} // create a new identifier node and set its token and value fields

		identifiers = append(identifiers, ident) // append it to the identifiers slice
	}

	if !p.expectPeek(token.RPAREN) { // if the next token is not a right parenthesis
		return nil
	}

	return identifiers
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
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

		program.Statements = append(program.Statements, stmt) // append it to the Statements field

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

	p.nextToken()                          // advance the tokens
	stmt.Value = p.parseExpression(LOWEST) // parse the value

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

	stmt.ReturnValue = p.parseExpression(LOWEST) // parse the return value

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
		p.noPrefixParseFnError(p.curToken.Type) // add an error to the errors slice
		return nil                              // return nil
	}

	leftExp := prefix() // parse the prefix expression

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type] // get the infix parse function for the next token type

		if infix == nil { // if there is no infix parse function
			return leftExp // return the left expression
		}

		p.nextToken() // advance the tokens

		leftExp = infix(leftExp) // parse the infix expression
	}

	return leftExp
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken} // create a new integer literal node and set its token field

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64) // parse the integer literal

	if err != nil { // if there was an error
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal) // create an error message
		p.errors = append(p.errors, msg)                                        // append it to the errors slice
		return nil                                                              // return nil
	}

	lit.Value = value // set the value field

	return lit
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t) // create an error message
	p.errors = append(p.errors, msg)                               // append it to the errors slice
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal} // create a new prefix expression node and set its token and operator fields

	p.nextToken() // advance the tokens

	expression.Right = p.parseExpression(PREFIX) // parse the right expression

	return expression
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok { // if the peek token is in the precedences map
		return p // return its precedence
	}

	return LOWEST // otherwise, return the lowest precedence
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok { // if the current token is in the precedences map
		return p // return its precedence
	}

	return LOWEST // otherwise, return the lowest precedence
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{Token: p.curToken, Operator: p.curToken.Literal, Left: left} // create a new infix expression node and set its token, operator, and left fields

	precedence := p.curPrecedence() // get the precedence of the current token

	p.nextToken() // advance the tokens

	expression.Right = p.parseExpression(precedence) // parse the right expression

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken() // advance the tokens

	exp := p.parseExpression(LOWEST) // parse the expression

	if !p.expectPeek(token.RPAREN) { // if the next token is not a right parenthesis
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken} // create a new if expression node and set its token field

	if !p.expectPeek(token.LPAREN) { // if the next token is not a left parenthesis
		return nil
	}

	p.nextToken() // advance the tokens

	expression.Condition = p.parseExpression(LOWEST) // parse the condition

	if !p.expectPeek(token.RPAREN) { // if the next token is not a right parenthesis
		return nil
	}

	if !p.expectPeek(token.LBRACE) { // if the next token is not a left brace
		return nil
	}

	expression.Consequence = p.parseBlockStatement() // parse the consequence

	if p.peekTokenIs(token.ELSE) { // if the next token is an else
		p.nextToken() // advance the tokens

		if !p.expectPeek(token.LBRACE) { // if the next token is not a left brace
			return nil
		}

		expression.Alternative = p.parseBlockStatement() // parse the alternative
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken} // create a new block statement node and set its token field

	block.Statements = []ast.Statement{} // initialize the Statements field to an empty slice

	p.nextToken() // advance the tokens

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) { // loop until we reach the end of the block
		stmt := p.parseStatement() // parse a statement

		if stmt != nil { // if the statement is not nil
			block.Statements = append(block.Statements, stmt) // append it to the Statements field
		}

		p.nextToken() // advance the tokens
	}

	return block
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}
