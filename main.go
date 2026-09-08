package main

import (
	"bufio"
	"cardist/scanner"
	"fmt"
	"os"
)

func main() { //checks the command line arguments and runs the appropriate mode: REPL or file scanning
	args := os.Args[1:]

	switch { //chooses which mode to run based on the command line arguments
	case len(args) == 0: //runs the REPL
		runPrompt()
	case len(args) == 2 && args[0] == "--tokenize": //if the user wants to tokenize a file, it will call runFile with the file path
		runFile(args[1])
	case len(args) == 1: //Prints out an error message and exit with code 70 (unhandled error since its yet to be implemented)
		fmt.Fprintln(os.Stderr, "Error: execution is not implemented until Lab 4. Use --tokenize <file>.") //
		os.Exit(70)
	default: //if the user provides invalid arguments, it will print out the usage message and exit with code 64 (invalid command line usage)
		fmt.Fprintln(os.Stderr, "Usage: run [--tokenize] [file]")
		os.Exit(64)
	}
}

func runFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not read file %q: %v\n", path, err)
		os.Exit(66)
	}

	scn := scanner.NewScanner(string(data))
	tokens := scn.ScanTokens()

	// Diagnostics already went to stderr, inside the scanner
	// stdout only prints the token stream if the scan came back clean
	if scn.ErrorFound() {
		os.Exit(65)
	}

	for _, tok := range tokens {
		fmt.Println(tok)
	}
	os.Exit(0)
}

/*
runPrompt is the REPL: scans one line at a time.  Marking errors but still proceeds
it just prints its error and the prompt returns
*/
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
		os.Exit(1) // indicates an error while reading input from the user
	}
}