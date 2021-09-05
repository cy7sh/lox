package scanner

import (
	"fmt"
	"strconv"

	"github.com/singurty/lox/token"
)

// Scanner transforms the source into tokens
type Scanner struct {
	source string
	tokens []token.Token
	start int
	current int
	line int
	HadError bool
}

// Map keywords to indentifiers
var keywords = map[string]token.Type {
	"and":		token.AND,
	"class":	token.CLASS,
	"else":		token.ELSE,
	"false":	token.FALSE,
	"for":		token.FOR,
	"fun":		token.FUN,
	"if":		token.IF,
	"null":		token.NULL,
	"or":		token.OR,
	"print":	token.PRINT,
	"return":	token.RETURN,
	"super":	token.SUPER,
	"this":		token.THIS,
	"true":		token.TRUE,
	"var":		token.VAR,
	"while":	token.WHILE,
	"break":	token.BREAK,
}

func New(source string) Scanner {
	scanner := Scanner{source: source, tokens: make([]token.Token, 0), line: 1}
	return scanner
}

func (sc *Scanner) ScanTokens() ([]token.Token) {
	for !sc.isAtEnd() {
		// beginning of the next lexeme
		sc.start = sc.current
		sc.scanToken()
	}
	sc.tokens = append(sc.tokens, token.Token{Type: token.EOF})
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
		case ':':
			sc.addToken(token.COLON)
		case '?':
			sc.addToken(token.QUESTION_MARK)
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
			} else if sc.match('*') {
				// block comment goes ultil */
				for (sc.peek() != '*' && sc.peekNext() != '/' && !sc.isAtEnd()) {
					if sc.peek() == '\n' {
						sc.line++
					}
					sc.advance()
				}
				// unterminated comment
				if sc.isAtEnd() {
					sc.handleError(sc.line, "Unterminated block comment")
				} else {
					// consume * and /
					sc.advance()
					sc.advance()
				}
			} else {
				sc.addToken(token.SLASH)
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
		case '"':
			sc.scanString()
			break
		default:
			if isDigit(c) {
				sc.scanNumber()
			} else if isAlpha(c) {
				sc.scanIdentifier()
			} else {
				sc.handleError(sc.line, fmt.Sprintf("Unexpected character: %c", c))
			}
			break
	}
	return c
}

func (sc *Scanner) scanIdentifier() {
	for isAlphaNumeric(sc.peek()) {
		sc.advance()
	}

	text := sc.source[sc.start:sc.current]
	tokenType, found := keywords[text]
	if found {
		sc.addToken(tokenType)
	} else {
		sc.addToken(token.IDENTIFIER)
	}
}

func (sc *Scanner) scanNumber() {
	for isDigit(sc.peek()) {
		sc.advance()
	}

	// Look for fractional part
	if (sc.peek() == '.' && isDigit(sc.peekNext())) {
		// Consume the .
		sc.advance()

		for isDigit(sc.peek()) {
			sc.advance()
		}
	}
	number, err := strconv.ParseFloat(sc.source[sc.start:sc.current], 64)
	if err != nil {
		sc.handleError(sc.line, "Invalid number")
	}
	// Parse as float
	sc.addTokenWithLiteral(token.NUMBER, number)
}

func (sc * Scanner) scanString() {
	for (sc.peek() != '"' && !sc.isAtEnd()) {
		if (sc.peek() == '\n') {
			sc.line++
		}
		sc.advance()
	}
	if sc.isAtEnd() {
		sc.handleError(sc.line, "Unterminated string")
		return
	}
	// The closing "
	sc.advance()
	// Trim surrounding quotes
	value := sc.source[sc.start+1 : sc.current-1]
	sc.addTokenWithLiteral(token.STRING, value)
}

func (sc *Scanner) handleError(line int, message string) {
	sc.HadError = true
	fmt.Printf("[Line %v] Error: %v\n", line, message)
}

func isDigit(c byte) bool {
	if _, err := strconv.Atoi(string(c)); err == nil {
		return true
	}
	return false
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '-'
}

func isAlphaNumeric(c byte) bool {
	return isDigit(c) || isAlpha(c)
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

func (sc *Scanner) peekNext() byte {
	if sc.current + 1 >= len(sc.source) {
		return 0x00
	}
	return sc.source[sc.current + 1]
}
