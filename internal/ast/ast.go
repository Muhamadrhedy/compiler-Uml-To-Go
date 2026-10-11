// Package ast defines the abstract syntax tree for UML2Code.
package ast

// Program is the root node of a parsed UML2Code source file.
type Program struct {
	Classes      []ClassDecl
	Associations []AssociationDecl
}

// ClassDecl represents one class declaration. Extends is empty when there is no parent.
type ClassDecl struct {
	Name       string
	Extends    string
	Attributes []AttributeDecl
	Methods    []MethodDecl
	Line       int
}

// AssociationDecl represents a directed, single-valued association: Source -> Target.
type AssociationDecl struct {
	Source, Target string
	Line           int
}

// AttributeDecl represents a class attribute.
type AttributeDecl struct {
	Name, Type, Visibility string
	Line                   int
}

// MethodDecl represents a method signature. Method bodies are generated later.
type MethodDecl struct {
	Name                   string
	Parameters             []Parameter
	ReturnType, Visibility string
	Line                   int
}

// Parameter represents one method parameter.
type Parameter struct {
	Name, Type string
	Line       int
}
