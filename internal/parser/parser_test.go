package parser

import (
	"strings"
	"testing"

	"uml2code/internal/lexer"
)

func parseSource(t *testing.T, source string) error {
	t.Helper()
	tokens, err := lexer.New(source).Scan()
	if err != nil {
		return err
	}
	_, err = New(tokens).Parse()
	return err
}

func TestParseClassWithAttributesAndMethods(t *testing.T) {
	source := `class User {
    - id: int
    - name: string
    + login(): bool
    + rename(newName: string): void
}`
	tokens, err := lexer.New(source).Scan()
	if err != nil {
		t.Fatal(err)
	}
	program, err := New(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Classes) != 1 {
		t.Fatalf("got %d classes, want 1", len(program.Classes))
	}
	class := program.Classes[0]
	if class.Name != "User" {
		t.Errorf("class name = %q, want User", class.Name)
	}
	if len(class.Attributes) != 2 {
		t.Errorf("got %d attributes, want 2", len(class.Attributes))
	}
	if len(class.Methods) != 2 {
		t.Errorf("got %d methods, want 2", len(class.Methods))
	}
	if got := class.Methods[1].Parameters[0].Name; got != "newName" {
		t.Errorf("parameter name = %q, want newName", got)
	}
}

func TestParseMultipleClasses(t *testing.T) {
	source := "class A {\n + run(): void\n}\n\nclass B {\n - id: int\n}"
	tokens, err := lexer.New(source).Scan()
	if err != nil {
		t.Fatal(err)
	}
	program, err := New(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Classes) != 2 {
		t.Fatalf("got %d classes, want 2", len(program.Classes))
	}
}

func TestParseRejectsInvalidSyntax(t *testing.T) {
	err := parseSource(t, "class User {\n - id int\n}")
	if err == nil {
		t.Fatal("expected syntax error")
	}
	if !strings.Contains(err.Error(), "line 2") {
		t.Errorf("expected line number in error, got %v", err)
	}
}

func TestParseRejectsVoidAttribute(t *testing.T) {
	err := parseSource(t, "class User {\n - id: void\n}")
	if err == nil {
		t.Fatal("expected void attribute to be rejected")
	}
}
