package scanner

import (
	"fmt"

	"github.com/singurty/lox/parseerror"
	"github.com/singurty/lox/token"
)

// Scanner transforms the source into tokens
type Scanner struct {
	source string
	tokens []token.Token
	start int
	current int
	line int
}

func New(source string) Scanner {
	scanner := Scanner{source: source, tokens: make([]token.Token, 0), start: 0, current: 0, line: 1}
	return scanner
}

func (sc *Scanner) ScanTokens() ([]token.Token) {
	for !sc.isAtEnd() {
		// beginning of the next lexeme
		sc.start = sc.current
		sc.scanToken()
	}
	sc.addToken(token.EOF)
	return sc.tokens
}

func (sc *Scanner) scanToken() (byte) {
	c := sc.advance()
	switch (c) {
		case '(':
			sc.addToken(token.LEFT_PAREN)
			break
		case ')':
			sc.addToken(token.RIGHT_PAREN)
			break
		case '{':
			sc.addToken(token.LEFT_BRACE)
			break
		case '}':
			sc.addToken(token.RIGHT_BRACE)
			break
		case ',':
			sc.addToken(token.COMMA)
			break
		case '.':
			sc.addToken(token.DOT)
			break
		case '-':
			sc.addToken(token.MINUS)
			break
		case '+':
			sc.addToken(token.PLUS)
			break
		case ';':
			sc.addToken(token.SEMICOLON)
			break
		case '*':
			sc.addToken(token.STAR)
			break
		case '!':
			if sc.match('=') {
				sc.addToken(token.BANG_EQUAL)
			} else {
				sc.addToken(token.BANG)
			}
			break
		case '=':
			if sc.match('=') {
				sc.addToken(token.EQUAL_EQUAL)
			} else {
				sc.addToken(token.EQUAL)
			}
			break
		case '<':
			if sc.match('=') {
				sc.addToken(token.LESS_EQUAL)
			} else {
				sc.addToken(token.LESS)
			}
			break
		case '>':
			if sc.match('=') {
				sc.addToken(token.GREATER_EQUAL)
			} else {
				sc.addToken(token.GREATER)
			}
			break
		case '/':
			if sc.match('/') {
				// a commment goes until the end of the line
				for (sc.peek() != '\n' && !sc.isAtEnd()) {
					sc.advance()
				}
			}
			break
		case ' ':
		case '\r':
		case '\t':
			// Ignore whitespace
			break;
		case '\n':
			sc.line++
			break
		default:
			parseerror.HadError = true
			parseerror.Error(sc.line, fmt.Sprintf("Unexpected character: %c", c))
			break
	}
	return c
}

func (sc * Scanner) scanString() {
	for (sc.peek() != '"' && !sc.isAtEnd()) {
		if (sc.peek() == '\n') {
			sc.line++
		}
		sc.advance()
	}
	if sc.isAtEnd() {
		parseerror.Error(sc.line, "Unterminated string")
		return
	}
	// The closing "
	sc.advance()
	// Trim surrounding quotes
	value := sc.source[sc.start+1 : sc.current-1]
	sc.addTokenWithLiteral(token.STRING, value)
}

func (sc *Scanner) advance() byte {
	c := sc.source[sc.current]
	sc.current++
	return c
}

func (sc *Scanner) addToken(tokenType token.Type) {
	sc.addTokenWithLiteral(tokenType, nil)
}

func (sc *Scanner) addTokenWithLiteral(tokenType token.Type, literal interface{}) {
	text := sc.source[sc.start:sc.current]
	sc.tokens = append(sc.tokens, token.Token{Type:tokenType, Lexeme: text, Literal: literal, Line: sc.line})
}

func (sc *Scanner) isAtEnd() bool {
	return sc.current >= len(sc.source)
}

func (sc *Scanner) match(expected byte) bool {
	if sc.isAtEnd() {
		return false
	}
	if sc.source[sc.current] != expected {
		return false
	}
	sc.current++
	return true
}

func (sc *Scanner) peek() byte {
	if sc.isAtEnd() {
		return 0x00
	}
	return sc.source[sc.current]
}
