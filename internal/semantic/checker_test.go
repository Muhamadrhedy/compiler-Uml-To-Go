package semantic

import (
	"strings"
	"testing"

	"uml2code/internal/ast"
)

func TestCheckAcceptsValidProgram(t *testing.T) {
	program := &ast.Program{Classes: []ast.ClassDecl{{
		Name: "User", Line: 1,
		Attributes: []ast.AttributeDecl{{Name: "id", Type: "int", Line: 2}},
		Methods:    []ast.MethodDecl{{Name: "login", ReturnType: "bool", Line: 3, Parameters: []ast.Parameter{{Name: "attempts", Type: "int", Line: 3}}}},
	}}}
	if errors := Check(program); len(errors) != 0 {
		t.Fatalf("expected valid program, got %v", errors)
	}
}

func TestCheckFindsDuplicateClasses(t *testing.T) {
	program := &ast.Program{Classes: []ast.ClassDecl{{Name: "User", Line: 1}, {Name: "User", Line: 5}}}
	assertContains(t, Check(program), `class "User" is declared more than once`)
}

func TestCheckFindsDuplicateMembersAndParameters(t *testing.T) {
	program := &ast.Program{Classes: []ast.ClassDecl{{
		Name: "User", Line: 1,
		Attributes: []ast.AttributeDecl{{Name: "id", Type: "int", Line: 2}, {Name: "id", Type: "string", Line: 3}},
		Methods:    []ast.MethodDecl{{Name: "login", ReturnType: "bool", Line: 4, Parameters: []ast.Parameter{{Name: "x", Type: "int", Line: 4}, {Name: "x", Type: "string", Line: 4}}}},
	}}}
	errors := Check(program)
	assertContains(t, errors, `member "id" is declared more than once`)
	assertContains(t, errors, `parameter "x" is duplicated`)
}

func TestCheckRejectsVoidAttributeAndParameter(t *testing.T) {
	program := &ast.Program{Classes: []ast.ClassDecl{{
		Name: "User", Line: 1,
		Attributes: []ast.AttributeDecl{{Name: "value", Type: "void", Line: 2}},
		Methods:    []ast.MethodDecl{{Name: "save", ReturnType: "void", Line: 3, Parameters: []ast.Parameter{{Name: "value", Type: "void", Line: 3}}}},
	}}}
	errors := Check(program)
	assertContains(t, errors, `type "void" is not allowed for attribute`)
	assertContains(t, errors, `type "void" is not allowed for parameter`)
}

func TestCheckRejectsNilAST(t *testing.T) {
	if errors := Check(nil); len(errors) != 1 || !strings.Contains(errors[0].Error(), "nil") {
		t.Fatalf("expected nil AST diagnostic, got %v", errors)
	}
}

func assertContains(t *testing.T, errors []error, part string) {
	t.Helper()
	for _, err := range errors {
		if strings.Contains(err.Error(), part) {
			return
		}
	}
	t.Fatalf("expected an error containing %q, got %v", part, errors)
}

func TestCheckRejectsProgramWithoutClasses(t *testing.T) {
	assertContains(t, Check(&ast.Program{}), "at least one class")
}
