package lexer

import (
	"fmt"
	"unicode"
)

// Lexer converts U2C source text into a sequence of tokens.
type Lexer struct {
	source      []rune
	start       int
	current     int
	line        int
	column      int
	startLine   int
	startColumn int
	tokens      []Token
}

func New(source string) *Lexer {
	return &Lexer{
		source:      []rune(source),
		line:        1,
		column:      1,
		startLine:   1,
		startColumn: 1,
	}
}

// Scan returns all tokens, ending with TokenEOF. It reports unexpected characters.
func (l *Lexer) Scan() ([]Token, error) {
	for !l.isAtEnd() {
		l.start = l.current
		l.startLine = l.line
		l.startColumn = l.column
		if err := l.scanToken(); err != nil {
			return nil, err
		}
	}
	l.tokens = append(l.tokens, Token{Kind: TokenEOF, Line: l.line, Column: l.column})
	return l.tokens, nil
}

func (l *Lexer) scanToken() error {
	c := l.advance()
	switch c {
	case ' ', '\t', '\r':
		// Ignore horizontal whitespace.
	case '\n':
		l.add(TokenNewline)
		l.line++
		l.column = 1
	case '+':
		l.add(TokenPlus)
	case '-':
		if l.match('>') {
			l.tokens = append(l.tokens, Token{Kind: TokenArrow, Lexeme: "->", Line: l.startLine, Column: l.startColumn})
		} else {
			l.add(TokenMinus)
		}
	case '#':
		l.add(TokenHash)
	case ':':
		l.add(TokenColon)
	case ',':
		l.add(TokenComma)
	case '(':
		l.add(TokenLParen)
	case ')':
		l.add(TokenRParen)
	case '{':
		l.add(TokenLBrace)
	case '}':
		l.add(TokenRBrace)
	case '/':
		if l.match('/') {
			for l.peek() != '\n' && !l.isAtEnd() {
				l.advance()
			}
		} else {
			return l.unexpected(c)
		}
	default:
		if isIdentifierStart(c) {
			for isIdentifierPart(l.peek()) {
				l.advance()
			}
			l.addIdentifierOrKeyword()
		} else {
			return l.unexpected(c)
		}
	}
	return nil
}

func (l *Lexer) add(kind Kind) {
	lexeme := string(l.source[l.start:l.current])
	l.tokens = append(l.tokens, Token{Kind: kind, Lexeme: lexeme, Line: l.startLine, Column: l.startColumn})
}

func (l *Lexer) addIdentifierOrKeyword() {
	lexeme := string(l.source[l.start:l.current])
	kind := TokenIdentifier
	switch lexeme {
	case "class":
		kind = TokenClass
	case "int":
		kind = TokenInt
	case "float":
		kind = TokenFloat
	case "string":
		kind = TokenString
	case "bool":
		kind = TokenBool
	case "void":
		kind = TokenVoid
	}
	l.tokens = append(l.tokens, Token{Kind: kind, Lexeme: lexeme, Line: l.startLine, Column: l.startColumn})
}

func (l *Lexer) unexpected(c rune) error {
	return fmt.Errorf("unexpected character %q at line %d, column %d", c, l.startLine, l.startColumn)
}

func (l *Lexer) isAtEnd() bool { return l.current >= len(l.source) }
func (l *Lexer) advance() rune {
	c := l.source[l.current]
	l.current++
	l.column++
	return c
}
func (l *Lexer) peek() rune {
	if l.isAtEnd() {
		return 0
	}
	return l.source[l.current]
}
func (l *Lexer) match(expected rune) bool {
	if l.isAtEnd() || l.source[l.current] != expected {
		return false
	}
	l.current++
	l.column++
	return true
}
func isIdentifierStart(c rune) bool { return c == '_' || unicode.IsLetter(c) }
func isIdentifierPart(c rune) bool  { return isIdentifierStart(c) || unicode.IsDigit(c) }
