# UML2Code — Stage 4

UML2Code is a small compiler project that transforms a text-based UML Class Diagram representation into Go source code.

## Implemented stages

1. Lexer — converts source text into tokens.
2. Parser — validates syntax and builds an AST.
3. Semantic analyzer — checks duplicate classes/members/parameters and supported types.
4. Go code generator — maps AST nodes into formatted Go declarations and method stubs.

## Run tests

```bash
go test ./...
go vet ./...
```

## Generator API

```go
source, err := generator.Generate(program)
if err != nil {
    // handle invalid AST or generation error
}
// source contains formatted Go code
```

## Example

Input `examples/example.uml` declares `User` and `Product` classes. The generator maps `float` to `float64`, converts public identifiers to exported Go names, and emits method stubs. Method bodies are not inferred from UML and are intentionally left as TODO stubs.

## Current limitation

The generator is available as a Go package. CLI integration (reading a `.uml` file and writing a `.go` file) is the next stage.
