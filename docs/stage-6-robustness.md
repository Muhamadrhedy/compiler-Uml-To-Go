# Stage 6 — Robustness and end-to-end tests

## Goal

Improve reliability around edge cases so UML2Code rejects problematic models clearly and does not write generated files when compilation fails.

## Changes

1. Semantic validation now rejects an empty program. This aligns `-check` behavior with the compiler's requirement that at least one class be present.
2. The Go generator detects distinct UML class names that normalize to the same Go type identifier, such as `user` and `User`.
3. The Go generator detects a field and a method that normalize to the same Go identifier.
4. Tests verify that empty or malformed input does not create an output file.
5. Tests verify that the input file is not overwritten if it is also specified as the output path.

## Verification commands

```bash
go test ./...
go vet ./...
go run ./cmd/uml2code examples/example.uml -o output.go
go build output.go
go run ./cmd/uml2code -check examples/example.uml
```

## Test coverage added

- Empty semantic program
- Class identifier normalization collision
- Field/method identifier collision
- Empty input CLI failure without output creation
- Malformed syntax CLI failure without output creation
- Input/output path collision protection

## Limitations

The language remains intentionally small. Inheritance, associations, constructors, generic types, and inferred method bodies are not implemented. Generated methods are stubs and should not be treated as finished business logic.
