// Package semantic performs semantic checks on a parsed UML2Code AST.
package semantic

import (
	"fmt"
	"strings"

	"uml2code/internal/ast"
)

// Diagnostic describes a semantic error in the source model.
type Diagnostic struct {
	Line    int
	Message string
}

func (d Diagnostic) Error() string {
	if d.Line > 0 {
		return fmt.Sprintf("semantic error at line %d: %s", d.Line, d.Message)
	}
	return "semantic error: " + d.Message
}

var allowedTypes = map[string]bool{
	"int": true, "float": true, "string": true, "bool": true, "void": true,
}

// Check validates program-level and member-level rules. It returns all
// diagnostics found, rather than stopping at the first semantic error.
func Check(program *ast.Program) []error {
	if program == nil {
		return []error{Diagnostic{Message: "program AST is nil"}}
	}

	var diagnostics []error
	if len(program.Classes) == 0 {
		return []error{Diagnostic{Line: 1, Message: "program must declare at least one class"}}
	}
	classes := make(map[string]int)
	classByName := make(map[string]ast.ClassDecl)
	for _, class := range program.Classes {
		if previousLine, exists := classes[class.Name]; exists {
			diagnostics = append(diagnostics, Diagnostic{Line: class.Line, Message: fmt.Sprintf("class %q is declared more than once (first declared at line %d)", class.Name, previousLine)})
		} else {
			classes[class.Name] = class.Line
			classByName[class.Name] = class
		}
		diagnostics = append(diagnostics, checkClass(class)...)
	}
	for _, class := range program.Classes {
		if class.Extends != "" {
			if _, exists := classByName[class.Extends]; !exists {
				diagnostics = append(diagnostics, Diagnostic{Line: class.Line, Message: fmt.Sprintf("parent class %q referenced by %q is not declared", class.Extends, class.Name)})
			}
			if class.Extends == class.Name {
				diagnostics = append(diagnostics, Diagnostic{Line: class.Line, Message: fmt.Sprintf("class %q cannot inherit from itself", class.Name)})
			}
		}
	}
	// Detect inheritance cycles using a three-state DFS.
	state := make(map[string]int)
	var visit func(string, []string)
	visit = func(name string, path []string) {
		if state[name] == 2 {
			return
		}
		if state[name] == 1 {
			diagnostics = append(diagnostics, Diagnostic{Message: fmt.Sprintf("inheritance cycle detected: %s", strings.Join(append(path, name), " -> "))})
			return
		}
		state[name] = 1
		if c, ok := classByName[name]; ok && c.Extends != "" {
			if _, exists := classByName[c.Extends]; exists {
				visit(c.Extends, append(path, name))
			}
		}
		state[name] = 2
	}
	for name := range classByName {
		visit(name, nil)
	}
	seenAssociations := make(map[string]bool)
	for _, association := range program.Associations {
		if _, exists := classByName[association.Source]; !exists {
			diagnostics = append(diagnostics, Diagnostic{Line: association.Line, Message: fmt.Sprintf("association source class %q is not declared", association.Source)})
		}
		if _, exists := classByName[association.Target]; !exists {
			diagnostics = append(diagnostics, Diagnostic{Line: association.Line, Message: fmt.Sprintf("association target class %q is not declared", association.Target)})
		}
		key := association.Source + "->" + association.Target
		if seenAssociations[key] {
			diagnostics = append(diagnostics, Diagnostic{Line: association.Line, Message: fmt.Sprintf("association %s -> %s is declared more than once", association.Source, association.Target)})
		}
		seenAssociations[key] = true
	}
	return diagnostics
}

func checkClass(class ast.ClassDecl) []error {
	var diagnostics []error
	// MVP rule: member names must be unique within a class. This avoids
	// ambiguous generated Go fields/methods and keeps the language simple.
	members := make(map[string]int)
	for _, attribute := range class.Attributes {
		if line, exists := members[attribute.Name]; exists {
			diagnostics = append(diagnostics, Diagnostic{Line: attribute.Line, Message: fmt.Sprintf("member %q is declared more than once in class %q (first declared at line %d)", attribute.Name, class.Name, line)})
		} else {
			members[attribute.Name] = attribute.Line
		}
		if !isAllowedType(attribute.Type, false) {
			diagnostics = append(diagnostics, Diagnostic{Line: attribute.Line, Message: fmt.Sprintf("type %q is not allowed for attribute %q", attribute.Type, attribute.Name)})
		}
	}
	for _, method := range class.Methods {
		if line, exists := members[method.Name]; exists {
			diagnostics = append(diagnostics, Diagnostic{Line: method.Line, Message: fmt.Sprintf("member %q is declared more than once in class %q (first declared at line %d)", method.Name, class.Name, line)})
		} else {
			members[method.Name] = method.Line
		}
		if !isAllowedType(method.ReturnType, true) {
			diagnostics = append(diagnostics, Diagnostic{Line: method.Line, Message: fmt.Sprintf("return type %q is not allowed for method %q", method.ReturnType, method.Name)})
		}
		parameters := make(map[string]int)
		for _, parameter := range method.Parameters {
			if line, exists := parameters[parameter.Name]; exists {
				diagnostics = append(diagnostics, Diagnostic{Line: parameter.Line, Message: fmt.Sprintf("parameter %q is duplicated in method %q (first declared at line %d)", parameter.Name, method.Name, line)})
			} else {
				parameters[parameter.Name] = parameter.Line
			}
			if !isAllowedType(parameter.Type, false) {
				diagnostics = append(diagnostics, Diagnostic{Line: parameter.Line, Message: fmt.Sprintf("type %q is not allowed for parameter %q", parameter.Type, parameter.Name)})
			}
		}
	}
	return diagnostics
}

func isAllowedType(typ string, allowVoid bool) bool {
	if !allowedTypes[typ] {
		return false
	}
	return allowVoid || !strings.EqualFold(typ, "void")
}
