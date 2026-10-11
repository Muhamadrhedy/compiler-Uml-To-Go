package compiler

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestCompileEndToEnd(t *testing.T) {
	source := `class User {
    - id: int
    + login(): bool
}`
	got, err := Compile(source)
	if err != nil {
		t.Fatalf("Compile() error = %v", err)
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "generated.go", got, parser.AllErrors); err != nil {
		t.Fatalf("generated source is invalid Go: %v\n%s", err, got)
	}
	for _, want := range []string{"type User struct", "func (receiver *User) Login() bool", "func main() {}"} {
		if !strings.Contains(string(got), want) {
			t.Errorf("generated source missing %q:\n%s", want, got)
		}
	}
}

func TestCompileRejectsSemanticError(t *testing.T) {
	source := `class User {
    - id: int
    - id: string
}`
	if _, err := Compile(source); err == nil || !strings.Contains(err.Error(), "declared more than once") {
		t.Fatalf("expected duplicate member diagnostic, got %v", err)
	}
}
