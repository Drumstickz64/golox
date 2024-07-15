package parsing

import (
	"fmt"

	"github.com/Drumstickz64/golox/assert"
	"github.com/Drumstickz64/golox/ast"
	"github.com/Drumstickz64/golox/errors"
	"github.com/Drumstickz64/golox/token"
)

type ParseError struct {
	Token       token.Token
	Message     any
	ShouldPanic bool
}

func (e ParseError) Error() string {
	var where string
	if e.Token.Kind == token.EOF {
		where = "at end"
	} else {
		where = "at '" + e.Token.Lexeme + "'"
	}

	return fmt.Sprintf("[on %d:%d] error %v: %v", e.Token.Line, e.Token.Column, where, e.Message)
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

func (p *Parser) Parse() ([]ast.Stmt, []errors.BuildError) {
	statements := []ast.Stmt{}
	errs := []errors.BuildError{}
	for !p.isAtEnd() {
		statement, err := p.declaration()
		if err != nil {
			errs = append(errs, err)
			if err.ShouldPanic {
				continue
			}
		} else {
			statements = append(statements, statement)
		}
	}

	return statements, errs
}

func (p *Parser) declaration() (ast.Stmt, *ParseError) {
	if p.match(token.FUN) {
		fun, err := p.function("function")
		if err != nil {
			p.synchronize()
			return nil, err
		}

		return fun, nil
	}

	if p.match(token.VAR) {
		decl, err := p.varDeclaration()
		if err != nil {
			p.synchronize()
			return nil, err
		}

		return decl, nil
	}

	statement, err := p.statement()
	if err != nil {
		p.synchronize()
		return nil, err
	}

	return statement, nil
}

func (p *Parser) function(kind string) (*ast.FunctionStmt, *ParseError) {
	name, err := p.consume(token.IDENTIFIER, fmt.Sprintf("expected %s name", kind))
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(token.LEFT_PAREN, fmt.Sprintf("expected '(' after %s name", kind)); err != nil {
		return nil, err
	}

	parameters := []token.Token{}
	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				return nil, &ParseError{
					Token:       p.peek(),
					Message:     "can't have more than 255 parameters",
					ShouldPanic: false,
				}
			}

			parameter, err := p.consume(token.IDENTIFIER, "expected parameter name")
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, parameter)
			if !p.match(token.COMMA) {
				break
			}
		}
	}

	if _, err := p.consume(token.RIGHT_PAREN, fmt.Sprintf("expected ')' after %s parameters", kind)); err != nil {
		return nil, err
	}

	if _, err := p.consume(token.LEFT_BRACE, fmt.Sprintf("expected '{' before %s name", kind)); err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return &ast.FunctionStmt{
		Name:       name,
		Parameters: parameters,
		Body:       body,
	}, nil
}

func (p *Parser) varDeclaration() (ast.Stmt, *ParseError) {
	name, err := p.consume(token.IDENTIFIER, "expected variable name after 'var'")
	if err != nil {
		return nil, err
	}

	var initializer ast.Expr = nil
	if p.match(token.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(token.SEMICOLON, "expected ';' after variable declaration"); err != nil {
		return nil, err
	}

	return &ast.VarStmt{
		Name:        name,
		Initializer: initializer,
	}, nil
}

func (p *Parser) statement() (ast.Stmt, *ParseError) {
	if p.match(token.PRINT) {
		return p.printStatement()
	}

	if p.match(token.RETURN) {
		return p.returnStatement()
	}

	if p.match(token.FOR) {
		return p.forStatement()
	}

	if p.match(token.WHILE) {
		return p.whileStatement()
	}

	if p.match(token.IF) {
		return p.ifStatement()
	}

	if p.match(token.LEFT_BRACE) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}

		return &ast.BlockStmt{
			Statements: statements,
		}, nil
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() (ast.Stmt, *ParseError) {
	expression, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(token.SEMICOLON, "expected ';' after value"); err != nil {
		return nil, err
	}

	return &ast.PrintStmt{
		Expression: expression,
	}, err
}

func (p *Parser) returnStatement() (ast.Stmt, *ParseError) {
	keyword := p.previous()
	var valueExpr ast.Expr
	if !p.check(token.SEMICOLON) {
		var err *ParseError
		valueExpr, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(token.SEMICOLON, "expected ';' after return value"); err != nil {
		return nil, err
	}

	return &ast.ReturnStmt{
		Keyword: keyword,
		Value:   valueExpr,
	}, nil
}

func (p *Parser) forStatement() (ast.Stmt, *ParseError) {
	var err *ParseError
	if _, err = p.consume(token.LEFT_PAREN, "expected '(' after 'for'"); err != nil {
		return nil, err
	}

	var initializer ast.Stmt
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition ast.Expr
	if !p.check(token.SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(token.SEMICOLON, "expected ';' after condition in for loop"); err != nil {
		return nil, err
	}

	var increment ast.Expr
	if !p.check(token.RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(token.RIGHT_PAREN, "expected ')' after increment in for loop"); err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = &ast.BlockStmt{
			Statements: []ast.Stmt{body, &ast.ExpressionStmt{Expression: increment}},
		}
	}

	if condition != nil {
		body = &ast.WhileStmt{
			Condition: condition,
			Body:      body,
		}
	} else {
		body = &ast.WhileStmt{
			Condition: &ast.LiteralExpr{Value: true},
			Body:      body,
		}

	}

	if initializer != nil {
		body = &ast.BlockStmt{
			Statements: []ast.Stmt{initializer, body},
		}
	}

	return body, nil
}

func (p *Parser) whileStatement() (ast.Stmt, *ParseError) {
	if _, err := p.consume(token.LEFT_PAREN, "expected '(' after 'while'"); err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(token.RIGHT_PAREN, "expected ')' after condition in while statment"); err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &ast.WhileStmt{
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *Parser) ifStatement() (ast.Stmt, *ParseError) {
	if _, err := p.consume(token.LEFT_PAREN, "expected '(' after 'if'"); err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(token.RIGHT_PAREN, "expected ')' after if condition"); err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch ast.Stmt = nil
	if p.match(token.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &ast.IfStmt{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

func (p *Parser) expressionStatement() (ast.Stmt, *ParseError) {
	expression, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(token.SEMICOLON, "expected ';' after value"); err != nil {
		return nil, err
	}

	return &ast.ExpressionStmt{
		Expression: expression,
	}, err
}

func (p *Parser) block() ([]ast.Stmt, *ParseError) {
	statements := []ast.Stmt{}
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		statement, err := p.declaration()
		if err != nil {
			return nil, err
		}

		statements = append(statements, statement)
	}

	if _, err := p.consume(token.RIGHT_BRACE, "expected '}' after block"); err != nil {
		return nil, err
	}

	return statements, nil
}

func (p *Parser) expression() (ast.Expr, *ParseError) {
	return p.assignment()
}

func (p *Parser) assignment() (ast.Expr, *ParseError) {
	expr, err := p.logical_or()
	if err != nil {
		return nil, err
	}

	if p.match(token.EQUAL) {
		varExpr, isVarExpr := expr.(*ast.VariableExpr)
		if !isVarExpr {
			equals := p.previous()
			return nil, &ParseError{
				Token:       equals,
				Message:     "invalid assignment target",
				ShouldPanic: false,
			}
		}

		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		return &ast.AssignmentExpr{
			Name:  varExpr.Name,
			Value: value,
		}, nil
	}

	return expr, nil
}

func (p *Parser) logical_or() (ast.Expr, *ParseError) {
	expr, err := p.logical_and()
	if err != nil {
		return nil, err
	}

	for p.match(token.OR) {
		operator := p.previous()
		right, err := p.logical_and()
		if err != nil {
			return nil, err
		}

		expr = &ast.LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) logical_and() (ast.Expr, *ParseError) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(token.AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		expr = &ast.LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) equality() (ast.Expr, *ParseError) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) comparison() (ast.Expr, *ParseError) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) term() (ast.Expr, *ParseError) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.PLUS, token.MINUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) factor() (ast.Expr, *ParseError) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) unary() (ast.Expr, *ParseError) {
	if p.match(token.MINUS, token.BANG) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr := &ast.UnaryExpr{
			Operator: operator,
			Right:    right,
		}

		return expr, nil
	}

	return p.call()
}

func (p *Parser) call() (ast.Expr, *ParseError) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(token.LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee ast.Expr) (ast.Expr, *ParseError) {
	arguments := []ast.Expr{}
	if !p.check(token.RIGHT_PAREN) {
		for {
			argument, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, argument)

			if !p.match(token.COMMA) {
				break
			}
		}
	}

	paren, err := p.consume(token.RIGHT_PAREN, "expected ')' after arguments")
	if err != nil {
		return nil, err
	}

	return &ast.CallExpr{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}, nil
}

func (p *Parser) primary() (ast.Expr, *ParseError) {
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

	if p.match(token.IDENTIFIER) {
		return &ast.VariableExpr{
			Name: p.previous(),
		}, nil
	}

	if p.match(token.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(token.RIGHT_PAREN, "Expected ')' after expression")
		if err != nil {
			return nil, err
		}

		return &ast.GroupingExpr{Expression: expr}, nil
	}

	return nil, &ParseError{
		Token:       p.peek(),
		Message:     "failed to parse expression",
		ShouldPanic: true,
	}
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

func (p *Parser) consume(kind token.Kind, msg string) (token.Token, *ParseError) {
	if !p.check(kind) {
		return token.Token{}, &ParseError{
			Token:       p.peek(),
			Message:     msg,
			ShouldPanic: true,
		}
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
