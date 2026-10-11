# Stage 5 — CLI integration

The CLI wires the existing compiler stages together and provides these commands:

- `go run ./cmd/uml2code examples/example.uml -o output.go`: compile input and write Go source.
- `go run ./cmd/uml2code examples/example.uml`: write `examples/example.go` by default.
- `go run ./cmd/uml2code -check examples/example.uml`: run lexer, parser, and semantic validation without writing output.

Errors are written to standard error and the process exits with a non-zero status. The `internal/compiler` package exposes `Compile` for full generation and `ParseAndCheck` for validation-only use.
