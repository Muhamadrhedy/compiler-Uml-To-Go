// Package generator converts a validated UML2Code AST into Go source code.
package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"go/token"
	"strings"
	"unicode"
	"unicode/utf8"

	"uml2code/internal/ast"
)

// Generate converts a program AST into formatted Go source. Semantic validation
// should be run before calling Generate.
func Generate(program *ast.Program) ([]byte, error) {
	if program == nil {
		return nil, fmt.Errorf("cannot generate Go code from a nil AST")
	}
	if len(program.Classes) == 0 {
		return nil, fmt.Errorf("cannot generate Go code: program has no classes")
	}

	var out bytes.Buffer
	out.WriteString("package main\n\n")
	seenClasses := make(map[string]string)
	associationsBySource := make(map[string][]string)
	for _, association := range program.Associations {
		associationsBySource[association.Source] = append(associationsBySource[association.Source], association.Target)
	}
	for _, class := range program.Classes {
		className := goIdentifier(class.Name, true)
		if !token.IsIdentifier(className) || token.Lookup(className).IsKeyword() {
			return nil, fmt.Errorf("class %q cannot be converted to a valid Go identifier", class.Name)
		}
		if previous, exists := seenClasses[className]; exists {
			return nil, fmt.Errorf("classes %q and %q map to the same Go type %q", previous, class.Name, className)
		}
		seenClasses[className] = class.Name
		fmt.Fprintf(&out, "type %s struct {\n", className)
		seenFields := make(map[string]string)
		if class.Extends != "" {
			parentName := goIdentifier(class.Extends, true)
			if previous, exists := seenFields[parentName]; exists {
				return nil, fmt.Errorf("parent type %q conflicts with field %q in class %q", parentName, previous, class.Name)
			}
			seenFields[parentName] = "embedded parent " + class.Extends
			fmt.Fprintf(&out, "\t*%s\n", parentName)
		}
		for _, attr := range class.Attributes {
			name := goIdentifier(attr.Name, attr.Visibility == "+")
			if previous, exists := seenFields[name]; exists {
				return nil, fmt.Errorf("attributes %q and %q in class %q map to the same Go field %q", previous, attr.Name, class.Name, name)
			}
			seenFields[name] = attr.Name
			fmt.Fprintf(&out, "\t%s %s\n", name, goType(attr.Type))
		}
		for _, target := range associationsBySource[class.Name] {
			fieldName := goIdentifier(target, true)
			if previous, exists := seenFields[fieldName]; exists {
				return nil, fmt.Errorf("association field %q conflicts with %s in class %q", fieldName, previous, class.Name)
			}
			seenFields[fieldName] = "association to " + target
			fmt.Fprintf(&out, "\t%s *%s\n", fieldName, goIdentifier(target, true))
		}
		out.WriteString("}\n\n")

		seenMethods := make(map[string]string)
		for _, method := range class.Methods {
			name := goIdentifier(method.Name, method.Visibility == "+")
			if previous, exists := seenMethods[name]; exists {
				return nil, fmt.Errorf("methods %q and %q in class %q map to the same Go method %q", previous, method.Name, class.Name, name)
			}
			if previous, exists := seenFields[name]; exists {
				return nil, fmt.Errorf("attribute %q and method %q in class %q map to the same Go identifier %q", previous, method.Name, class.Name, name)
			}
			seenMethods[name] = method.Name
			fmt.Fprintf(&out, "func (receiver *%s) %s(", className, name)
			seenParams := make(map[string]string)
			for i, param := range method.Parameters {
				if i > 0 {
					out.WriteString(", ")
				}
				paramName := goIdentifier(param.Name, false)
				if previous, exists := seenParams[paramName]; exists {
					return nil, fmt.Errorf("parameters %q and %q in method %q map to the same Go identifier %q", previous, param.Name, method.Name, paramName)
				}
				seenParams[paramName] = param.Name
				fmt.Fprintf(&out, "%s %s", paramName, goType(param.Type))
			}
			out.WriteString(")")
			if method.ReturnType != "void" {
				fmt.Fprintf(&out, " %s", goType(method.ReturnType))
			}
			out.WriteString(" {\n\tpanic(\"TODO: implement method\")\n}\n\n")
		}
	}

	// The generated file is a standalone package main, so include an empty entry point.
	out.WriteString("func main() {}\n")

	formatted, err := format.Source(out.Bytes())
	if err != nil {
		return nil, fmt.Errorf("generated Go source is invalid: %w", err)
	}
	return formatted, nil
}

func goType(typ string) string {
	switch typ {
	case "int":
		return "int"
	case "float":
		return "float64"
	case "string":
		return "string"
	case "bool":
		return "bool"
	case "void":
		return ""
	default:
		return typ
	}
}

// goIdentifier maps UML visibility to Go's exported/unexported naming convention.
// Go has no protected visibility, so '#' maps to an unexported identifier too.
func goIdentifier(name string, exported bool) string {
	if name == "" {
		return "umlEmpty"
	}
	name = strings.TrimSpace(name)
	first, size := utf8.DecodeRuneInString(name)
	if exported {
		name = string(unicode.ToUpper(first)) + name[size:]
	} else {
		name = string(unicode.ToLower(first)) + name[size:]
	}
	if !token.IsIdentifier(name) || token.Lookup(name).IsKeyword() {
		name = "uml" + string(unicode.ToUpper([]rune(name)[0])) + string([]rune(name)[1:])
	}
	return name
}
