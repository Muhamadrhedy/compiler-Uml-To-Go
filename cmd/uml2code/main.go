// Command uml2code compiles UML2Code text files into Go source code.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"uml2code/internal/compiler"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "uml2code:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet("uml2code", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	outputPath := flags.String("o", "", "output Go file (default: input filename with .go extension)")
	checkOnly := flags.Bool("check", false, "validate the input without generating a Go file")
	showVersion := flags.Bool("version", false, "print version")
	flags.Usage = func() {
		fmt.Fprintln(flags.Output(), "UML2Code - compile text-based UML class declarations to Go")
		fmt.Fprintln(flags.Output(), "\nUsage:")
		fmt.Fprintln(flags.Output(), "  uml2code [options] <input.uml>")
		fmt.Fprintln(flags.Output(), "\nOptions:")
		flags.PrintDefaults()
	}
	args = normalizeArgs(args)
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}
	if *showVersion {
		fmt.Println("uml2code dev")
		return nil
	}
	if flags.NArg() != 1 {
		flags.Usage()
		return fmt.Errorf("expected exactly one input .uml file")
	}

	inputPath := flags.Arg(0)
	if filepath.Ext(inputPath) != ".uml" {
		return fmt.Errorf("input file must have a .uml extension: %s", inputPath)
	}
	source, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("cannot read input %q: %w", inputPath, err)
	}

	if *checkOnly {
		if _, err := compiler.ParseAndCheck(string(source)); err != nil {
			return err
		}
		fmt.Printf("OK: %s passed lexical, syntax, and semantic checks\n", inputPath)
		return nil
	}

	generated, err := compiler.Compile(string(source))
	if err != nil {
		return err
	}
	if *outputPath == "" {
		*outputPath = strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".go"
	}
	if filepath.Clean(*outputPath) == filepath.Clean(inputPath) {
		return fmt.Errorf("output path must be different from input path")
	}
	if err := os.WriteFile(*outputPath, generated, 0644); err != nil {
		return fmt.Errorf("cannot write output %q: %w", *outputPath, err)
	}
	fmt.Printf("Generated %s from %s\n", *outputPath, inputPath)
	return nil
}

// normalizeArgs lets users place options before or after the input filename.
func normalizeArgs(args []string) []string {
	var options, positional []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-o" || arg == "--output":
			options = append(options, arg)
			if i+1 < len(args) {
				i++
				options = append(options, args[i])
			}
		case strings.HasPrefix(arg, "-o=") || strings.HasPrefix(arg, "--output="):
			options = append(options, "-o="+strings.SplitN(arg, "=", 2)[1])
		case arg == "-check" || arg == "-version" || arg == "-h" || arg == "--help":
			options = append(options, arg)
		default:
			positional = append(positional, arg)
		}
	}
	return append(options, positional...)
}
