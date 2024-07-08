package scanning

import (
	"fmt"
	"strconv"

	"github.com/Drumstickz64/golox/assert"
	"github.com/Drumstickz64/golox/errors"
	"github.com/Drumstickz64/golox/token"
)

func scanningError(line int, msg any) error {
	return errors.NewBuildtimeError(line, "", msg)
}

type Scanner struct {
	source         string
	tokens         []token.Token
	start, current int
	line           int
}

func NewScanner(source string) Scanner {
	return Scanner{
		source: source,
		line:   1,
	}
}

func (s *Scanner) ScanTokens() ([]token.Token, []error) {
	errs := []error{}
	for !s.isAtEnd() {
		s.start = s.current
		if err := s.scanToken(); err != nil {
			errs = append(errs, err)
		}
	}

	s.tokens = append(s.tokens, token.Token{
		Kind:    token.EOF,
		Lexeme:  "",
		Literal: nil,
		Line:    s.line,
	})

	return s.tokens, errs
}

func (s *Scanner) scanToken() error {
	char := s.advance()
	switch char {
	case '(':
		s.addToken(token.LEFT_PAREN)
	case ')':
		s.addToken(token.RIGHT_PAREN)
	case '{':
		s.addToken(token.LEFT_BRACE)
	case '}':
		s.addToken(token.RIGHT_BRACE)
	case ',':
		s.addToken(token.COMMA)
	case '.':
		s.addToken(token.DOT)
	case '-':
		s.addToken(token.MINUS)
	case '+':
		s.addToken(token.PLUS)
	case ';':
		s.addToken(token.SEMICOLON)
	case '*':
		s.addToken(token.STAR)
	case '!':
		s.addCompoundToken('=', token.BANG_EQUAL, token.BANG)
	case '=':
		s.addCompoundToken('=', token.EQUAL_EQUAL, token.EQUAL)
	case '<':
		s.addCompoundToken('=', token.LESS_EQUAL, token.LESS)
	case '>':
		s.addCompoundToken('=', token.GREATER_EQUAL, token.GREATER)
	case '/':
		if s.match('/') {
			s.dropLine()
		} else if s.match('*') {
			if err := s.ignoreMultilineComment(); err != nil {
				return nil
			}
		} else {
			s.addToken(token.SLASH)
		}
	case '"':
		if err := s.addStringToken(); err != nil {
			return err
		}
	case ' ', '\t', '\r':
	case '\n':
		s.line++
	default:
		if isDigit(char) {
			s.addNumberToken()
		} else if isAlpha(char) {
			s.addIdentifierToken()
		} else {
			return scanningError(s.line, fmt.Sprintf("Found unexpected character '%v'", string(char)))
		}
	}

	return nil
}

func (s *Scanner) advance() rune {
	char := rune(s.source[s.current])
	s.current++
	return char
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\x00'
	}

	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\x00'
	}

	return rune(s.source[s.current+1])

}

func (s *Scanner) addToken(kind token.Kind) {
	s.addLiteralToken(kind, nil)
}

func (s *Scanner) addLiteralToken(kind token.Kind, literal any) {
	s.tokens = append(s.tokens, token.Token{
		Kind:    kind,
		Lexeme:  s.source[s.start:s.current],
		Literal: literal,
		Line:    s.line,
	})
}

// used for adding "compound" tokens if the correct matching token comes
// after the original token to complete the compound token.
// Otherwise adds the simple token
//
// For example:
//
// != is compound
//
// ! is simple
func (s *Scanner) addCompoundToken(completingChar rune, compound, simple token.Kind) {
	if s.match(completingChar) {
		s.addToken(compound)
	} else {
		s.addToken(simple)
	}

}

func (s *Scanner) addStringToken() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}

		s.advance()
	}

	if s.isAtEnd() {
		return scanningError(s.line, "Unterminated string")
	}

	// consume last "
	s.advance()

	// + 1 and - 1 to trim off the " characters
	literal := s.source[s.start+1 : s.current-1]
	s.addLiteralToken(token.STRING, literal)

	return nil
}

func (s *Scanner) addNumberToken() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance() // consume .
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	lexeme := s.source[s.start:s.current]
	literal, err := strconv.ParseFloat(lexeme, 64)
	assert.That(err == nil, fmt.Sprintf("'%v' is a valid number literal: %v", lexeme, err))
	s.addLiteralToken(token.NUMBER, literal)
}

func (s *Scanner) addIdentifierToken() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	lexeme := s.source[s.start:s.current]
	kind, isKeyword := token.Keywords[lexeme]
	if !isKeyword {
		kind = token.IDENTIFIER
	}

	s.addToken(kind)
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != byte(expected) {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) dropLine() {
	for s.peek() != '\n' && !s.isAtEnd() {
		s.advance()
	}
}

func (s *Scanner) ignoreMultilineComment() error {
	for !(s.peek() == '*' && s.peekNext() == '/') && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}

		s.advance()
	}

	if s.isAtEnd() {
		return scanningError(s.line, "Unterminated multi-line comment, expected '*/'")
	}

	// consume the closing */
	s.advance() // consume *
	s.advance() // consume /

	return nil
}

func isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

func isAlpha(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '_'
}

func isAlphaNumeric(char rune) bool {
	return isDigit(char) || isAlpha(char)
}
