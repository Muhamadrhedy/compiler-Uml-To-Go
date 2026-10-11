# UML2Code Language (U2C) — MVP Specification

## Goal
UML2Code reads a small text representation of UML class diagrams and generates Go source code. The MVP supports classes, attributes, methods, visibility markers, primitive types, and line comments.

## Example
```text
class User {
    - id: int
    - name: string
    + login(): bool
    + updateName(newName: string): void
}
```

## Lexical rules
- Keywords: `class`, `int`, `float`, `string`, `bool`, `void`.
- Identifier: starts with a letter or `_`, then contains letters, digits, or `_`.
- Visibility: `+` public, `-` private, `#` protected.
- Punctuation: `{ } ( ) : ,`.
- Newlines separate member declarations. Spaces, tabs, and carriage returns are ignored.
- A `//` comment runs to the end of the line.

## EBNF grammar
```ebnf
Program       = { ClassDecl } ;
ClassDecl     = "class" , Identifier , "{" , { MemberDecl } , "}" ;
MemberDecl    = AttributeDecl | MethodDecl ;
AttributeDecl = Visibility , Identifier , ":" , Type , Newline ;
MethodDecl    = Visibility , Identifier , "(" , [ ParameterList ] , ")" , ":" , Type , Newline ;
ParameterList = Parameter , { "," , Parameter } ;
Parameter     = Identifier , ":" , Type ;
Visibility    = "+" | "-" | "#" ;
Type          = "int" | "float" | "string" | "bool" | "void" ;
Identifier    = LetterOrUnderscore , { Letter | Digit | "_" } ;
```

## Semantic constraints
- `void` is permitted only as a method return type, not as an attribute or parameter type.
- Class names must be unique in one program.
- Attribute and method names must not be duplicated within the same class in the MVP.
- The initial generator emits method stubs; it does not infer business logic from UML.

## Planned pipeline
Source → Lexer → Parser → AST → Semantic Analyzer → Go Code Generator.
```
