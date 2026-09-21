package main

import (
	"bufio"
	"cardist/parser"
	"cardist/scanner"
	"fmt"
	"os"
)

func main() { //checks the command line arguments and runs the appropriate mode: REPL or file scanning
	args := os.Args[1:]

	switch { //chooses which mode to run based on the command line arguments
	case len(args) == 0: //runs the REPL
		runPrompt()
	case len(args) == 2 && args[0] == "--tokenize": //if the user wants to tokenize a file, it will call runTokenizeFile with the file path
		runTokenizeFile(args[1])
	case len(args) == 2 && args[0] == "--parse": //if the user wants to parse a file, it will call runParseFile with the file path
		runParseFile(args[1])
	case len(args) == 1: //Prints out an error message and exit with code 70 (unhandled error since its yet to be implemented)
		fmt.Fprintln(os.Stderr, "Error: execution is not implemented until Lab 4. Use --tokenize <file> or --parse <file>.")
		os.Exit(70)
	default: //if the user provides invalid arguments, it will print out the usage message and exit with code 64 (invalid command line usage)
		fmt.Fprintln(os.Stderr, "Usage: run [--tokenize | --parse] [file]")
		os.Exit(64)
	}
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

/*
runPrompt is the REPL: scans and parses one line at a time.
*/
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
