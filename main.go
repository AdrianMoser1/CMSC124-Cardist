package main

import (
	"bufio"
	"cardist/parser"
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
		runTokenizeFile(args[1])
	case len(args) == 2 && args[0] == "--parse":
		runParseFile(args[1])
	case len(args) == 1:
		runEcho(args[0])
	default:
		fmt.Fprintln(os.Stderr, "Usage: run [--tokenize | --parse] [file]")
		os.Exit(64)
	}
}

func runEcho(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not read file %q: %v\n", path, err)
		os.Exit(66)
	}
	fmt.Print(string(data))
	os.Exit(0)
}

func runTokenizeFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not read file %q: %v\n", path, err)
		os.Exit(66)
	}

	scn := scanner.NewScanner(string(data))
	tokens := scn.ScanTokens()

	if scn.ErrorFound() {
		os.Exit(65)
	}

	for _, tok := range tokens {
		fmt.Println(tok)
	}
	os.Exit(0)
}

func runParseFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not read file %q: %v\n", path, err)
		os.Exit(66)
	}

	scn := scanner.NewScanner(string(data))
	tokens := scn.ScanTokens()

	if scn.ErrorFound() {
		os.Exit(65)
	}

	prs := parser.NewParser(tokens)
	expressions := prs.Parse()

	if prs.ErrorFound() {
		os.Exit(65)
	}

	printer := parser.NewAstPrinter()
	for _, expr := range expressions {
		if expr != nil {
			fmt.Println(printer.Print(expr))
		}
	}
	os.Exit(0)
}

func runPrompt() {
	reader := bufio.NewScanner(os.Stdin)
	printer := parser.NewAstPrinter()
	fmt.Print("> ")

	for reader.Scan() {
		line := reader.Text()
		scn := scanner.NewScanner(line)
		tokens := scn.ScanTokens()

		if !scn.ErrorFound() {
			prs := parser.NewParser(tokens)
			expressions := prs.Parse()

			if !prs.ErrorFound() {
				for _, expr := range expressions {
					if expr != nil {
						fmt.Println(printer.Print(expr))
					}
				}
			}
		}
		fmt.Print("> ")
	}

	if err := reader.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: reading input: %v\n", err)
		os.Exit(1)
	}
}
