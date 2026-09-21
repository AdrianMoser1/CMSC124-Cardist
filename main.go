package main

import (
	"bufio"
	"cardist/scanner"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	switch {
	case len(args) == 0:
		runPrompt()
	case len(args) == 2 && args[0] == "--tokenize":
		runFile(args[1])
	case len(args) == 1:
		// Lab 0 checks the file pipeline before the language has execution
		// semantics, so preserve the source instead of interpreting it.
		runEcho(args[0])
	default:
		fmt.Fprintln(os.Stderr, "Usage: run [--tokenize] [file]")
		os.Exit(64)
	}
}

// Lab 0 tests the command pipeline before execution exists, so this mode must
// preserve the input exactly rather than attempting to interpret it.
func runEcho(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not read file %q: %v\n", path, err)
		os.Exit(66)
	}
	fmt.Print(string(data))
	os.Exit(0)
}

func runFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not read file %q: %v\n", path, err)
		os.Exit(66)
	}

	scn := scanner.NewScanner(string(data))
	tokens := scn.ScanTokens()

	// Keep stdout reserved for complete token streams; diagnostics belong on
	// stderr, and a lexical error must be reflected in the process status.
	if scn.ErrorFound() {
		os.Exit(65)
	}

	for _, tok := range tokens {
		fmt.Println(tok)
	}
	os.Exit(0)
}

// The REPL scans each line independently so one malformed input does not end
// the session or prevent the user from trying the next line.
func runPrompt() {
	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for reader.Scan() {
		line := reader.Text()
		scn := scanner.NewScanner(line)
		tokens := scn.ScanTokens()
		for _, tok := range tokens {
			fmt.Println(tok)
		}
		fmt.Print("> ")
	}

	if err := reader.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: reading input: %v\n", err)
		os.Exit(1)
	}
}
