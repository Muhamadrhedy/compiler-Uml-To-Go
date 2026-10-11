# Stage 4 — Go Code Generator

The generator converts a validated UML2Code AST into formatted Go source code.

## Mapping rules

- UML `int` -> Go `int`
- UML `float` -> Go `float64`
- UML `string` -> Go `string`
- UML `bool` -> Go `bool`
- UML `void` method return -> no Go return type
- UML public (`+`) members -> exported Go identifiers (initial uppercase)
- UML private (`-`) and protected (`#`) members -> unexported Go identifiers (initial lowercase); Go has no `protected` visibility
- Method bodies are generated as stubs containing `panic("TODO: implement method")`

The output is formatted with `go/format`. The generator detects collisions introduced by visibility-based name conversion, such as a private `name` and public `Name` field mapping to the same Go identifier.

## Test

```bash
go test ./...
go vet ./...
```

This stage exposes `generator.Generate(*ast.Program) ([]byte, error)`. CLI integration is planned for the next stage.
