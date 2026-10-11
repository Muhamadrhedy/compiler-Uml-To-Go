# UML2Code

A small educational compiler written in Go. It reads a constrained text representation of UML class diagrams and will generate Go source code.

## Current milestone
Milestone 1 defines the U2C language and implements the lexer. Parser, AST, semantic analysis, and code generation will be added in later milestones.

## Requirements
- Go 1.23 or newer

## Run lexer tests
```bash
go test ./...
```

## Current limitation
The project currently tokenizes input only; it does not yet parse UML files or generate Go code.

## Tahap 3 — Semantic Analyzer

Semantic analyzer sudah tersedia di `internal/semantic`. Ia mendeteksi class duplikat, nama member duplikat dalam class, parameter duplikat, serta penggunaan `void` yang tidak sah. Detail aturan ada di `docs/stage-3-semantic.md`.

Jalankan seluruh pengujian dengan `go test ./...` dan pemeriksaan statis dengan `go vet ./...`.
