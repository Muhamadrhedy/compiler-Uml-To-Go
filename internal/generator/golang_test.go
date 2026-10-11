package generator

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"uml2code/internal/ast"
)

func TestGenerateCreatesValidGo(t *testing.T) {
	program := &ast.Program{Classes: []ast.ClassDecl{{
		Name: "User",
		Attributes: []ast.AttributeDecl{
			{Name: "id", Type: "int", Visibility: "-"},
			{Name: "name", Type: "string", Visibility: "+"},
			{Name: "active", Type: "bool", Visibility: "#"},
			{Name: "price", Type: "float", Visibility: "-"},
		},
		Methods: []ast.MethodDecl{
			{Name: "login", ReturnType: "bool", Visibility: "+"},
			{Name: "updateName", Parameters: []ast.Parameter{{Name: "newName", Type: "string"}}, ReturnType: "void", Visibility: "-"},
			{Name: "getPrice", ReturnType: "float", Visibility: "+"},
		},
	}}}

	got, err := Generate(program)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "generated.go", got, parser.AllErrors); err != nil {
		t.Fatalf("generated source is not valid Go: %v\n%s", err, got)
	}
	for _, want := range []string{
		"type User struct",
		"id     int",
		"Name string",
		"active bool",
		"price  float64",
		"func (receiver *User) Login() bool",
		"func (receiver *User) updateName(newName string)",
		"func (receiver *User) GetPrice() float64",
		"panic(\"TODO: implement method\")",
	} {
		if !strings.Contains(string(got), want) {
			t.Errorf("generated source missing %q:\n%s", want, got)
		}
	}
}

func TestGenerateRejectsNilAndEmptyProgram(t *testing.T) {
	if _, err := Generate(nil); err == nil {
		t.Fatal("Generate(nil) expected an error")
	}
	if _, err := Generate(&ast.Program{}); err == nil {
		t.Fatal("Generate(empty program) expected an error")
	}
}

func TestGenerateRejectsFieldNameCollision(t *testing.T) {
	program := &ast.Program{Classes: []ast.ClassDecl{{
		Name: "User",
		Attributes: []ast.AttributeDecl{
			{Name: "foo", Type: "string", Visibility: "-"},
			{Name: "Foo", Type: "string", Visibility: "#"},
		},
	}}}
	if _, err := Generate(program); err == nil {
		t.Fatal("expected generated field collision error")
	}
}

func TestGenerateRejectsClassNameCollision(t *testing.T) {
	program := &ast.Program{Classes: []ast.ClassDecl{
		{Name: "user"},
		{Name: "User"},
	}}
	if _, err := Generate(program); err == nil || !strings.Contains(err.Error(), "same Go type") {
		t.Fatalf("expected class name collision error, got %v", err)
	}
}

func TestGenerateRejectsFieldAndMethodNameCollision(t *testing.T) {
	program := &ast.Program{Classes: []ast.ClassDecl{{
		Name:       "User",
		Attributes: []ast.AttributeDecl{{Name: "login", Type: "bool", Visibility: "+"}},
		Methods:    []ast.MethodDecl{{Name: "login", ReturnType: "bool", Visibility: "+"}},
	}}}
	if _, err := Generate(program); err == nil || !strings.Contains(err.Error(), "same Go identifier") {
		t.Fatalf("expected field/method collision error, got %v", err)
	}
}
