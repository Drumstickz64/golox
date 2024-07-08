package parsing

import (
	"github.com/Drumstickz64/golox/assert"
	"github.com/Drumstickz64/golox/ast"
	"github.com/Drumstickz64/golox/errors"
	"github.com/Drumstickz64/golox/token"
)

func parsingError(tok token.Token, msg any) error {
	if tok.Kind == token.EOF {
		return errors.NewBuildtimeError(tok.Line, " at end", msg)
	}

	return errors.NewBuildtimeError(tok.Line, " at '"+tok.Lexeme+"'", msg)
}

type Parser struct {
	tokens  []token.Token
	current int
}

func NewParser(tokens []token.Token) Parser {
	return Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() (ast.Expr, error) {
	return p.expression()
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (ast.Expr, error) {
	exp, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		exp = &ast.BinaryExpr{
			Left:     exp,
			Operator: operator,
			Right:    right,
		}
	}

	return exp, nil
}

func (p *Parser) comparison() (ast.Expr, error) {
	exp, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		exp = &ast.BinaryExpr{
			Left:     exp,
			Operator: operator,
			Right:    right,
		}
	}

	return exp, nil
}

func (p *Parser) term() (ast.Expr, error) {
	exp, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.PLUS, token.MINUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		exp = &ast.BinaryExpr{
			Left:     exp,
			Operator: operator,
			Right:    right,
		}
	}

	return exp, nil
}

func (p *Parser) factor() (ast.Expr, error) {
	exp, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		exp = &ast.BinaryExpr{
			Left:     exp,
			Operator: operator,
			Right:    right,
		}
	}

	return exp, nil
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(token.MINUS, token.BANG) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		exp := &ast.UnaryExpr{
			Operator: operator,
			Right:    right,
		}

		return exp, nil
	}

	return p.primary()
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(token.TRUE) {
		return &ast.LiteralExpr{Value: true}, nil
	}

	if p.match(token.FALSE) {
		return &ast.LiteralExpr{Value: false}, nil
	}

	if p.match(token.NIL) {
		return &ast.LiteralExpr{Value: nil}, nil
	}

	if p.match(token.STRING, token.NUMBER) {
		return &ast.LiteralExpr{Value: p.previous().Literal}, nil
	}

	if p.match(token.LEFT_PAREN) {
		exp, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(token.RIGHT_PAREN, "Expected ')' after expression")
		if err != nil {
			return nil, err
		}

		return &ast.GroupingExpr{Expression: exp}, nil
	}

	return nil, parsingError(p.peek(), "Failed to parse expression")
}

func (p *Parser) match(kinds ...token.Kind) bool {
	for _, kind := range kinds {
		if p.check(kind) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(kind token.Kind) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Kind == kind
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) peek() token.Token {
	// p.isAtEnd() not needed because EOF is already at the end of tokens
	// so it would be returned here
	return p.tokens[p.current]
}

func (p *Parser) previous() token.Token {
	assert.That(p.current > 0, "parser is not at the start of token list")
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Kind == token.EOF
}

func (p *Parser) consume(kind token.Kind, msg string) (token.Token, error) {
	if !p.check(kind) {
		return token.Token{}, parsingError(p.peek(), msg)
	}

	return p.advance(), nil
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Kind == token.SEMICOLON {
			return
		}

		switch p.peek().Kind {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}
