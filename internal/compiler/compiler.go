// Package compiler wires the UML2Code compilation pipeline together.
package compiler

import (
	"fmt"

	"uml2code/internal/ast"
	"uml2code/internal/generator"
	"uml2code/internal/lexer"
	"uml2code/internal/parser"
	"uml2code/internal/semantic"
)

// Compile runs lexical analysis, parsing, semantic checks, and Go code generation.
func Compile(source string) ([]byte, error) {
	tokens, err := lexer.New(source).Scan()
	if err != nil {
		return nil, fmt.Errorf("lexical analysis: %w", err)
	}

	program, err := parser.New(tokens).Parse()
	if err != nil {
		return nil, fmt.Errorf("syntax analysis: %w", err)
	}

	if diagnostics := semantic.Check(program); len(diagnostics) > 0 {
		return nil, fmt.Errorf("%s", joinErrors(diagnostics))
	}

	output, err := generator.Generate(program)
	if err != nil {
		return nil, fmt.Errorf("code generation: %w", err)
	}
	return output, nil
}

// ParseAndCheck validates a source file without generating code.
func ParseAndCheck(source string) (*ast.Program, error) {
	tokens, err := lexer.New(source).Scan()
	if err != nil {
		return nil, fmt.Errorf("lexical analysis: %w", err)
	}
	program, err := parser.New(tokens).Parse()
	if err != nil {
		return nil, fmt.Errorf("syntax analysis: %w", err)
	}
	if diagnostics := semantic.Check(program); len(diagnostics) > 0 {
		return nil, fmt.Errorf("%s", joinErrors(diagnostics))
	}
	return program, nil
}

func joinErrors(errors []error) string {
	message := "semantic analysis failed:"
	for _, err := range errors {
		message += "\n  - " + err.Error()
	}
	return message
}
