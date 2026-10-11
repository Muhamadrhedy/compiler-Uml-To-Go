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

func TestRunRejectsEmptyInputWithoutWritingOutput(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "empty.uml")
	output := filepath.Join(dir, "output.go")
	if err := os.WriteFile(input, nil, 0644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"-o", output, input}); err == nil || !strings.Contains(err.Error(), "at least one class") {
		t.Fatalf("expected empty input error, got %v", err)
	}
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatalf("output should not be created on compile failure; stat error = %v", err)
	}
}

func TestRunRejectsOutputEqualToInput(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "sample.uml")
	original := []byte("class User {\n    - id: int\n}\n")
	if err := os.WriteFile(input, original, 0644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"-o", input, input}); err == nil || !strings.Contains(err.Error(), "different from input") {
		t.Fatalf("expected output/input path error, got %v", err)
	}
	got, err := os.ReadFile(input)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(original) {
		t.Fatalf("input file was modified unexpectedly:\n%s", got)
	}
}

func TestRunRejectsMalformedInputWithoutWritingOutput(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "bad.uml")
	output := filepath.Join(dir, "bad.go")
	if err := os.WriteFile(input, []byte("class User {\n - id int\n}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"-o", output, input}); err == nil {
		t.Fatal("expected malformed input error")
	}
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatalf("output should not be created on syntax failure; stat error = %v", err)
	}
}
