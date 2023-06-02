package Scanner

import (
	"fmt"
	"strconv"

	"github.com/AnshVM/golox/Error"
	"github.com/AnshVM/golox/Tokens"
)

type Scanner struct {
	source  string
	tokens  []*Tokens.Token
	start   uint
	current uint
	line    uint
}

func NewScanner(source string) Scanner {
	return Scanner{source: source, tokens: []*Tokens.Token{}}
}

func (scanner *Scanner) ScanTokens() []*Tokens.Token {
	for !scanner.isAtEnd() {
		scanner.start = scanner.current
		scanner.scanToken()
	}
	//start = current for case where the input ends with comment
	scanner.start = scanner.current
	scanner.addToken(Tokens.EOF, nil)
	return scanner.tokens
}

func (scanner *Scanner) isAtEnd() bool {
	return scanner.current >= uint(len(scanner.source))
}

func (scanner *Scanner) scanToken() {
	c := scanner.advance()
	switch c {
	case '(':
		scanner.addToken(Tokens.LEFT_PAREN, nil)
		break
	case ')':
		scanner.addToken(Tokens.RIGHT_PAREN, nil)
		break
	case '{':
		scanner.addToken(Tokens.LEFT_BRACE, nil)
		break
	case '}':
		scanner.addToken(Tokens.RIGHT_BRACE, nil)
		break
	case ',':
		scanner.addToken(Tokens.COMMA, nil)
		break
	case '.':
		scanner.addToken(Tokens.DOT, nil)
		break
	case '-':
		scanner.addToken(Tokens.MINUS, nil)
		break
	case '+':
		scanner.addToken(Tokens.PLUS, nil)
		break
	case ';':
		scanner.addToken(Tokens.SEMICOLON, nil)
		break
	case '*':
		scanner.addToken(Tokens.STAR, nil)
		break
	case '?':
		scanner.addToken(Tokens.QUESTION_MARK, nil)
		break
	case ':':
		scanner.addToken(Tokens.COLON, nil)
		break
	case '!':
		scanner.matchAddToken('=', Tokens.BANG_EQUAL, Tokens.BANG)
		break
	case '=':
		scanner.matchAddToken('=', Tokens.EQUAL_EQUAL, Tokens.EQUAL)
		break
	case '>':
		scanner.matchAddToken('=', Tokens.GREATER_EQUAL, Tokens.GREATER)
		break
	case '<':
		scanner.matchAddToken('=', Tokens.LESS_EQAUL, Tokens.LESS)
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
			scanner.addToken(Tokens.SLASH, nil)
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
		Error.ReportScanError(scanner.line, fmt.Sprintf("Unexpected token: %c", c))
	}
}

func (scanner *Scanner) identifier() {
	for isAlphaNumeirc(scanner.peek()) {
		scanner.advance()
	}
	text := scanner.source[scanner.start:scanner.current]
	if Tokens.Keywords[text] == "" {
		scanner.addToken(Tokens.IDENTIFIER, nil)
	} else {
		scanner.addToken(Tokens.Keywords[text], nil)
	}

}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeirc(c byte) bool {
	return isAlpha(c) || isDigit(c)
}

func (scanner *Scanner) number() {
	if scanner.peek() == '.' && isDigit(scanner.peekNext()) {
		scanner.advance()
	}
	for isDigit(scanner.peek()) {
		scanner.advance()
		if scanner.peek() == '.' && isDigit(scanner.peekNext()) {
			scanner.advance()
		}
	}
	value, _ := strconv.ParseFloat(scanner.source[scanner.start:scanner.current], 32)
	scanner.addToken(Tokens.NUMBER, float32(value))
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
		Error.ReportScanError(scanner.line, "Unterminated string")
		return
	}
	scanner.advance()
	value := scanner.source[scanner.start+1 : scanner.current-1]
	scanner.addToken(Tokens.STRING, value)
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
		&Tokens.Token{Type: tokenType, Lexeme: lexeme, Literal: literal, Line: scanner.line},
	)
}

func (scanner *Scanner) advance() byte {
	scanner.current++
	return scanner.source[scanner.current-1]
}
