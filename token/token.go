package token

import "fmt"

var Keywords = map[string]Kind{
	AND.String():    AND,
	CLASS.String():  CLASS,
	ELSE.String():   ELSE,
	FALSE.String():  FALSE,
	FUN.String():    FUN,
	FOR.String():    FOR,
	IF.String():     IF,
	NIL.String():    NIL,
	OR.String():     OR,
	PRINT.String():  PRINT,
	RETURN.String(): RETURN,
	SUPER.String():  SUPER,
	THIS.String():   THIS,
	TRUE.String():   TRUE,
	VAR.String():    VAR,
	WHILE.String():  WHILE,
}

type Kind int

const (
	// Single-character tokens.
	// (){},.-+;*/

	LEFT_PAREN Kind = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	STAR
	SLASH
	QUESTION
	COLON

	// One or two character tokens.
	// ! != = == > >= < <=

	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals.

	IDENTIFIER
	STRING
	NUMBER

	// Keywords.

	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

func (k Kind) String() string {
	switch k {
	case AND:
		return "and"
	case BANG:
		return "bang"
	case BANG_EQUAL:
		return "bang_equal"
	case CLASS:
		return "class"
	case COMMA:
		return "comma"
	case DOT:
		return "dot"
	case ELSE:
		return "else"
	case EOF:
		return "EOF"
	case EQUAL:
		return "equal"
	case EQUAL_EQUAL:
		return "equal_equal"
	case FALSE:
		return "false"
	case FOR:
		return "for"
	case FUN:
		return "fun"
	case GREATER:
		return "greater"
	case GREATER_EQUAL:
		return "greater_equal"
	case IDENTIFIER:
		return "identifier"
	case IF:
		return "if"
	case LEFT_BRACE:
		return "left_brace"
	case LEFT_PAREN:
		return "left_paren"
	case LESS:
		return "less"
	case LESS_EQUAL:
		return "less_equal"
	case MINUS:
		return "minus"
	case NIL:
		return "nil"
	case NUMBER:
		return "number"
	case OR:
		return "or"
	case PLUS:
		return "plus"
	case PRINT:
		return "print"
	case RETURN:
		return "return"
	case RIGHT_BRACE:
		return "right_brace"
	case RIGHT_PAREN:
		return "right_paren"
	case SEMICOLON:
		return "semicolon"
	case SLASH:
		return "slash"
	case STAR:
		return "star"
	case STRING:
		return "string"
	case SUPER:
		return "super"
	case THIS:
		return "this"
	case TRUE:
		return "true"
	case VAR:
		return "var"
	case WHILE:
		return "while"
	case QUESTION:
		return "question_mark"
	case COLON:
		return "colon"
	default:
		panic(fmt.Sprintf("unexpected token.Kind: %#v", k))
	}
}

type Token struct {
	Kind    Kind
	Lexeme  string
	Literal any
	Line    int
}

func (t Token) String() string {
	return fmt.Sprintf("token of kind '%v' scanned from '%s' with literal '%v'", t.Kind, t.Lexeme, t.Literal)
}
