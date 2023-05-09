package Scanner

import (
	"strconv"

	"github.com/AnshVM/golox/Error"
)

type Scanner struct {
	source  string
	tokens  []Token
	start   uint
	current uint
	line    uint
}

func NewScanner(source string) Scanner {
	return Scanner{source: source, tokens: []Token{}}
}

func (scanner *Scanner) ScanTokens() []Token {
	for !scanner.isAtEnd() {
		scanner.start = scanner.current
		scanner.scanToken()
	}
	//start = current for case where the input ends with comment
	scanner.start = scanner.current
	scanner.addToken(EOF, nil)
	return scanner.tokens
}

func (scanner *Scanner) isAtEnd() bool {
	return scanner.current >= uint(len(scanner.source))
}

func (scanner *Scanner) scanToken() {
	c := scanner.advance()
	switch c {
	case '(':
		scanner.addToken(LEFT_PAREN, nil)
		break
	case ')':
		scanner.addToken(RIGHT_PAREN, nil)
		break
	case '{':
		scanner.addToken(LEFT_BRACE, nil)
		break
	case '}':
		scanner.addToken(RIGHT_BRACE, nil)
		break
	case ',':
		scanner.addToken(COMMA, nil)
		break
	case '.':
		scanner.addToken(DOT, nil)
		break
	case '-':
		scanner.addToken(MINUS, nil)
		break
	case '+':
		scanner.addToken(PLUS, nil)
		break
	case ';':
		scanner.addToken(SEMICOLON, nil)
		break
	case '*':
		scanner.addToken(STAR, nil)
		break
	case '!':
		scanner.matchAddToken('=', BANG_EQUAL, BANG)
		break
	case '=':
		scanner.matchAddToken('=', EQUAL_EQUAL, EQUAL)
		break
	case '>':
		scanner.matchAddToken('=', GREATER_EQUAL, GREATER)
		break
	case '<':
		scanner.matchAddToken('=', LESS_EQAUL, LESS)
		break
	case '/':
		if scanner.match('/') {
			for scanner.peek() != '\n' && !scanner.isAtEnd() {
				scanner.advance()
			}
		} else if scanner.match('*') {
			for !scanner.isAtEnd() {
				if scanner.peek() == '*' && scanner.peekNext() == '/' {
					scanner.advance()
					scanner.advance()
					break
				}
				scanner.advance()
			}
		} else {
			scanner.addToken(SLASH, nil)
		}
		break
	case '"':
		scanner.string()
		break

	case ' ':
	case '\r':
	case '\t':
		break

	case '\n':
		scanner.line++
		break

	default:
		if isDigit(c) {
			scanner.number()
			break
		}
		if isAlpha(c) {
			scanner.identifier()
			break
		}

		Error.Error(scanner.line, "Unexpected character.")
	}
}

func (scanner *Scanner) identifier() {
	for isAlphaNumeirc(scanner.peek()) {
		scanner.advance()
	}
	text := scanner.source[scanner.start:scanner.current]
	if keywords[text] == "" {
		scanner.addToken(IDENTIFIER, nil)
	} else {
		scanner.addToken(keywords[text], nil)
	}

}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'B') || c == '_'
}

func isAlphaNumeirc(c byte) bool {
	return isAlpha(c) || isDigit(c)
}

func (scanner *Scanner) number() {
	for isDigit(scanner.peek()) {
		scanner.advance()
		if scanner.peek() == '.' && isDigit(scanner.peekNext()) {
			scanner.advance()
		}
	}
	value, _ := strconv.ParseFloat(scanner.source[scanner.start:scanner.current], 32)
	scanner.addToken(NUMBER, value)
}

func (scanner *Scanner) peekNext() byte {
	if scanner.current == uint(len(scanner.source)-1) {
		return byte(0)
	}
	return scanner.source[scanner.current+1]
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (scanner *Scanner) string() {
	for scanner.peek() != '"' && !scanner.isAtEnd() {
		if scanner.peek() == '\n' {
			scanner.line++
		}
		scanner.advance()
	}
	if scanner.isAtEnd() {
		Error.Error(scanner.line, "Unterminated string")
	}
	scanner.advance()
	value := scanner.source[scanner.start+1 : scanner.current-1]
	scanner.addToken(STRING, value)
}

func (scanner *Scanner) peek() byte {
	if scanner.isAtEnd() {
		return 0
	}
	return scanner.source[scanner.current]
}

// because go does not support ternary
func (scanner *Scanner) matchAddToken(match_char byte, token_matched string, token_unmatched string) {
	if scanner.match(match_char) {
		scanner.addToken(token_matched, nil)
	} else {
		scanner.addToken(token_unmatched, nil)
	}
}

func (scanner *Scanner) match(c byte) bool {
	if scanner.isAtEnd() {
		return false
	}
	if c != byte(scanner.source[scanner.current]) {
		return false
	}
	scanner.current++
	return true
}

func (scanner *Scanner) addToken(tokenType string, literal any) {
	lexeme := scanner.source[scanner.start:scanner.current]
	scanner.tokens = append(
		scanner.tokens,
		Token{Type: tokenType, Lexeme: lexeme, Literal: literal, Line: scanner.line},
	)
}

func (scanner *Scanner) advance() byte {
	scanner.current++
	return scanner.source[scanner.current-1]
}
