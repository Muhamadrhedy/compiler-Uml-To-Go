package lexer

import "fmt"

// Kind identifies the type of lexical token.
type Kind string

const (
	TokenEOF        Kind = "EOF"
	TokenNewline    Kind = "NEWLINE"
	TokenIdentifier Kind = "IDENTIFIER"
	TokenClass      Kind = "CLASS"
	TokenInt        Kind = "INT"
	TokenFloat      Kind = "FLOAT"
	TokenString     Kind = "STRING"
	TokenBool       Kind = "BOOL"
	TokenVoid       Kind = "VOID"

	TokenPlus   Kind = "+"
	TokenMinus  Kind = "-"
	TokenHash   Kind = "#"
	TokenColon  Kind = ":"
	TokenComma  Kind = ","
	TokenLParen Kind = "("
	TokenRParen Kind = ")"
	TokenLBrace Kind = "{"
	TokenRBrace Kind = "}"
)

// Token represents one lexical unit and its source location.
type Token struct {
	Kind   Kind
	Lexeme string
	Line   int
	Column int
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%q) at %d:%d", t.Kind, t.Lexeme, t.Line, t.Column)
}
