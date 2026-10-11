// Package parser implements recursive-descent parsing for UML2Code.
package parser

import (
	"fmt"

	"uml2code/internal/ast"
	"uml2code/internal/lexer"
)

// Parser consumes lexer tokens and builds an AST.
type Parser struct {
	tokens  []lexer.Token
	current int
}

// New creates a parser for the provided token stream.
func New(tokens []lexer.Token) *Parser { return &Parser{tokens: tokens} }

// Parse parses a complete program containing classes and optional associations.
func (p *Parser) Parse() (*ast.Program, error) {
	program := &ast.Program{}
	p.skipNewlines()
	for !p.check(lexer.TokenEOF) {
		if p.check(lexer.TokenClass) {
			classDecl, err := p.parseClass()
			if err != nil {
				return nil, err
			}
			program.Classes = append(program.Classes, classDecl)
		} else if p.check(lexer.TokenIdentifier) && p.peek().Lexeme == "association" {
			association, err := p.parseAssociation()
			if err != nil {
				return nil, err
			}
			program.Associations = append(program.Associations, association)
		} else {
			return nil, p.errorAt(p.peek(), "expected 'class' declaration or 'association' declaration")
		}
		if !p.check(lexer.TokenEOF) && !p.check(lexer.TokenNewline) {
			return nil, p.errorAt(p.peek(), "expected a newline after declaration")
		}
		p.skipNewlines()
	}
	return program, nil
}

func (p *Parser) parseAssociation() (ast.AssociationDecl, error) {
	keyword := p.advance()
	source, err := p.expect(lexer.TokenIdentifier, "expected source class name after 'association'")
	if err != nil {
		return ast.AssociationDecl{}, err
	}
	if _, err = p.expect(lexer.TokenArrow, "expected '->' between association classes"); err != nil {
		return ast.AssociationDecl{}, err
	}
	target, err := p.expect(lexer.TokenIdentifier, "expected target class name after '->'")
	if err != nil {
		return ast.AssociationDecl{}, err
	}
	return ast.AssociationDecl{Source: source.Lexeme, Target: target.Lexeme, Line: keyword.Line}, nil
}

func (p *Parser) parseClass() (ast.ClassDecl, error) {
	keyword, err := p.expect(lexer.TokenClass, "expected 'class' declaration")
	if err != nil {
		return ast.ClassDecl{}, err
	}
	name, err := p.expect(lexer.TokenIdentifier, "expected class name after 'class'")
	if err != nil {
		return ast.ClassDecl{}, err
	}
	classDecl := ast.ClassDecl{Name: name.Lexeme, Line: keyword.Line}
	// Optional inheritance clause: class Child extends Parent { ... }
	if p.check(lexer.TokenIdentifier) && p.peek().Lexeme == "extends" {
		p.advance()
		parent, err := p.expect(lexer.TokenIdentifier, "expected parent class name after 'extends'")
		if err != nil {
			return ast.ClassDecl{}, err
		}
		classDecl.Extends = parent.Lexeme
	}
	if _, err = p.expect(lexer.TokenLBrace, "expected '{' after class declaration"); err != nil {
		return ast.ClassDecl{}, err
	}
	p.skipNewlines()
	for !p.check(lexer.TokenRBrace) && !p.check(lexer.TokenEOF) {
		member, err := p.parseMember()
		if err != nil {
			return ast.ClassDecl{}, err
		}
		switch m := member.(type) {
		case ast.AttributeDecl:
			classDecl.Attributes = append(classDecl.Attributes, m)
		case ast.MethodDecl:
			classDecl.Methods = append(classDecl.Methods, m)
		}
		if !p.check(lexer.TokenRBrace) && !p.check(lexer.TokenNewline) && !p.check(lexer.TokenEOF) {
			return ast.ClassDecl{}, p.errorAt(p.peek(), "expected a newline or '}' after member declaration")
		}
		p.skipNewlines()
	}
	if _, err := p.expect(lexer.TokenRBrace, "expected '}' to close class"); err != nil {
		return ast.ClassDecl{}, err
	}
	return classDecl, nil
}

func (p *Parser) parseMember() (any, error) {
	visibilityToken := p.advance()
	visibility := visibilityToken.Lexeme
	if visibilityToken.Kind != lexer.TokenPlus && visibilityToken.Kind != lexer.TokenMinus && visibilityToken.Kind != lexer.TokenHash {
		return nil, p.errorAt(visibilityToken, "expected visibility '+' (public), '-' (private), or '#' (protected)")
	}
	name, err := p.expect(lexer.TokenIdentifier, "expected member name after visibility")
	if err != nil {
		return nil, err
	}
	if p.match(lexer.TokenColon) {
		typ, err := p.parseType()
		if err != nil {
			return nil, err
		}
		if typ.Lexeme == "void" {
			return nil, p.errorAt(typ, "attribute type cannot be void")
		}
		return ast.AttributeDecl{Name: name.Lexeme, Type: typ.Lexeme, Visibility: visibility, Line: name.Line}, nil
	}
	if p.match(lexer.TokenLParen) {
		method := ast.MethodDecl{Name: name.Lexeme, Visibility: visibility, Line: name.Line}
		if !p.check(lexer.TokenRParen) {
			for {
				paramName, err := p.expect(lexer.TokenIdentifier, "expected parameter name")
				if err != nil {
					return nil, err
				}
				if _, err := p.expect(lexer.TokenColon, "expected ':' after parameter name"); err != nil {
					return nil, err
				}
				paramType, err := p.parseType()
				if err != nil {
					return nil, err
				}
				if paramType.Lexeme == "void" {
					return nil, p.errorAt(paramType, "parameter type cannot be void")
				}
				method.Parameters = append(method.Parameters, ast.Parameter{Name: paramName.Lexeme, Type: paramType.Lexeme, Line: paramName.Line})
				if !p.match(lexer.TokenComma) {
					break
				}
			}
		}
		if _, err := p.expect(lexer.TokenRParen, "expected ')' after parameters"); err != nil {
			return nil, err
		}
		if _, err := p.expect(lexer.TokenColon, "expected ':' before method return type"); err != nil {
			return nil, err
		}
		returnType, err := p.parseType()
		if err != nil {
			return nil, err
		}
		method.ReturnType = returnType.Lexeme
		return method, nil
	}
	return nil, p.errorAt(p.peek(), "expected ':' for an attribute or '(' for a method")
}

func (p *Parser) parseType() (lexer.Token, error) {
	t := p.peek()
	switch t.Kind {
	case lexer.TokenInt, lexer.TokenFloat, lexer.TokenString, lexer.TokenBool, lexer.TokenVoid:
		p.advance()
		return t, nil
	default:
		return lexer.Token{}, p.errorAt(t, "expected type: int, float, string, bool, or void")
	}
}
func (p *Parser) expect(kind lexer.Kind, message string) (lexer.Token, error) {
	if p.check(kind) {
		return p.advance(), nil
	}
	return lexer.Token{}, p.errorAt(p.peek(), message)
}
func (p *Parser) match(kind lexer.Kind) bool {
	if !p.check(kind) {
		return false
	}
	p.advance()
	return true
}
func (p *Parser) check(kind lexer.Kind) bool { return p.peek().Kind == kind }
func (p *Parser) advance() lexer.Token {
	if p.current < len(p.tokens)-1 {
		p.current++
	}
	return p.tokens[p.current-1]
}
func (p *Parser) peek() lexer.Token {
	if len(p.tokens) == 0 {
		return lexer.Token{Kind: lexer.TokenEOF, Line: 1, Column: 1}
	}
	if p.current >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[p.current]
}
func (p *Parser) skipNewlines() {
	for p.match(lexer.TokenNewline) {
	}
}
func (p *Parser) errorAt(t lexer.Token, message string) error {
	return fmt.Errorf("parse error at line %d, column %d: %s (got %q)", t.Line, t.Column, message, t.Lexeme)
}
