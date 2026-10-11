# UML2Code — Stage 6 (Robustness & End-to-End Tests)

UML2Code is a small compiler that translates a constrained text-based UML class language into Go source code.

## Pipeline

`source .uml → Lexer → Parser → AST → Semantic Analyzer → Go Code Generator → source .go`

## Requirements

- Go 1.23 or newer

## Run tests and static checks

```bash
go test ./...
go vet ./...
```

## Compile a UML file

```bash
go run ./cmd/uml2code examples/example.uml -o output.go
```

If `-o` is omitted, the output uses the input name with a `.go` extension:

```bash
go run ./cmd/uml2code examples/example.uml
```

Validate without generating output:

```bash
go run ./cmd/uml2code -check examples/example.uml
```

Show help:

```bash
go run ./cmd/uml2code -h
```

## Validate generated Go

```bash
go build output.go
```

The generator emits `package main` and an empty `main()` so the generated file can be built as a standalone Go program. UML describes class structure and method signatures, not business logic, so generated methods contain a TODO panic stub. Replace those stubs with real implementations before using the generated program.

## Stage 6 robustness improvements

- Reject empty UML programs during semantic validation.
- Detect class names that collide after conversion to Go exported identifiers (for example `user` and `User`).
- Detect field/method name collisions after Go identifier conversion.
- Add CLI tests confirming invalid input does not create output files.
- Add CLI test confirming the input file is not overwritten when selected as output.
- Preserve compiler errors and avoid writing output until all compilation stages succeed.

## Current language scope

- Classes, attributes, and methods
- Visibility markers `+`, `-`, `#`
- Types `int`, `float`, `string`, `bool`, `void` (void only as a method return type)
- Duplicate declaration and supported-type checks
- Go identifier collision checks during generation

Inheritance, relationships, constructors, generics, and method-body inference are out of scope for this MVP.
