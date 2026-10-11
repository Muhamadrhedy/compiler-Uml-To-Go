package lexer

import "testing"

func TestScanClassDeclaration(t *testing.T) {
	source := "class User {\n    - id: int\n    + login(): bool\n}\n"
	got, err := New(source).Scan()
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	want := []Kind{
		TokenClass, TokenIdentifier, TokenLBrace, TokenNewline,
		TokenMinus, TokenIdentifier, TokenColon, TokenInt, TokenNewline,
		TokenPlus, TokenIdentifier, TokenLParen, TokenRParen, TokenColon, TokenBool, TokenNewline,
		TokenRBrace, TokenNewline, TokenEOF,
	}
	if len(got) != len(want) {
		t.Fatalf("got %d tokens, want %d", len(got), len(want))
	}
	for i, kind := range want {
		if got[i].Kind != kind {
			t.Errorf("token %d kind = %q, want %q", i, got[i].Kind, kind)
		}
	}
}

func TestScanComments(t *testing.T) {
	got, err := New("// a comment\nclass User {}\n").Scan()
	if err != nil {
		t.Fatal(err)
	}
	if got[0].Kind != TokenNewline {
		t.Fatalf("first token = %q, want NEWLINE", got[0].Kind)
	}
	if got[1].Kind != TokenClass {
		t.Fatalf("second token = %q, want CLASS", got[1].Kind)
	}
}

func TestUnexpectedCharacter(t *testing.T) {
	_, err := New("class User { @ }").Scan()
	if err == nil {
		t.Fatal("expected an error for unexpected character")
	}
}
