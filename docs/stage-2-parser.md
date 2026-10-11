# Stage 2 — Parser and AST

## Goal
Convert lexer tokens into an Abstract Syntax Tree (AST).

## AST nodes
- `Program`: all class declarations in a source file.
- `ClassDecl`: class name, attributes, methods, and source line.
- `AttributeDecl`: name, type, visibility, and source line.
- `MethodDecl`: name, parameters, return type, visibility, and source line.
- `Parameter`: parameter name, type, and source line.

## Parser behavior
The parser uses recursive descent. It accepts multiple classes, attributes, methods, comma-separated parameters, blank lines, and comments already skipped by the lexer. Syntax errors include line and column information.

## Current limitations
- Members are one per line.
- Only the primitive types `int`, `float`, `string`, `bool`, and `void` are recognized.
- `void` is rejected for attributes and parameters.
- Inheritance, associations, method bodies, and generics are not supported yet.

## Run tests

```bash
go test ./...
```
