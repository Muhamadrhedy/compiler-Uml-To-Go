package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunGeneratesGoFile(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "sample.uml")
	output := filepath.Join(dir, "sample.go")
	if err := os.WriteFile(input, []byte("class User {\n    - id: int\n}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"-o", output, input}); err != nil {
		t.Fatalf("run() error = %v", err)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "type User struct") {
		t.Fatalf("output missing generated type:\n%s", data)
	}
}

func TestRunCheckOnly(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "sample.uml")
	if err := os.WriteFile(input, []byte("class User {\n    - id: int\n}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"-check", input}); err != nil {
		t.Fatalf("run() check error = %v", err)
	}
}
