package main

import (
	"bufio"
	"cardist/evaluator"
	"cardist/parser"
	"cardist/scanner"
	"fmt"
	"os"
	"strings"
)

const usage = "Usage: run [--tokenize | --parse | --eval] [file]"

func main() {
	args := os.Args[1:]

	switch {
	case len(args) == 0:
		runPrompt()
	case len(args) == 2 && args[0] == "--tokenize":
		runTokenizeFile(args[1])
	case len(args) == 2 && args[0] == "--parse":
		runParseFile(args[1])
	case len(args) == 2 && args[0] == "--eval":
		runEvalFile(args[1])
	case len(args) == 1 && strings.HasPrefix(args[0], "--"):
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(64)
	case len(args) == 1:
		runScriptFile(args[0])
	default:
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(64)
	}
}

func runScriptFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Error: could not read file %q: %v\n",
			path,
			err,
		)
		os.Exit(66)
	}

	scn := scanner.NewScanner(string(data))
	tokens := scn.ScanTokens()

	if scn.ErrorFound() {
		os.Exit(65)
	}

	prs := parser.NewParser(tokens)
	script := prs.ParseScript()

	// Parse the whole script before executing any statements.
	if prs.ErrorFound() {
		os.Exit(65)
	}

	ev := evaluator.NewEvaluator()
	ev.Emit = printValue

	if err := ev.ExecuteScript(script); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(70)
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
	program := prs.Parse()

	if prs.ErrorFound() {
		os.Exit(65)
	}

	printer := parser.NewAstPrinter()
	if out := printer.Print(program); out != "" {
		fmt.Println(out)
	}
	os.Exit(0)
}

// runEvalFile is the Lab 3 tested path: ./run --eval <file>.
//
// Order matters for the exit codes:
//  1. scan + parse the WHOLE file first; any failure there is a static error
//     (exit 65) and nothing is evaluated or printed.
//  2. only then evaluate, printing each value as soon as it is computed, so
//     output produced before a runtime error still reaches stdout.
//  3. the first runtime error goes to stderr and the process exits 70.
func runEvalFile(path string) {
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
	program := prs.Parse()
	if prs.ErrorFound() {
		os.Exit(65)
	}

	ev := evaluator.NewEvaluator()
	ev.Emit = printValue
	// One call walks the whole tree from the Program root.
	if _, err := ev.Evaluate(program); err != nil {
		// err.Error() is already "[line N] Runtime error: <message>".
		fmt.Fprintln(os.Stderr, err)
		os.Exit(70)
	}
	os.Exit(0)
}

// printValue is the single place a computed value reaches stdout.
func printValue(value any) {
	fmt.Println(evaluator.FormatValue(value))
}

// runPrompt is the REPL. One evaluator for the whole session, so variables persist.
// shortcut: one line = one entry, so a block must fit on one line; buffer lines
// until braces balance if multi-line blocks are ever needed.
func runPrompt() {
	reader := bufio.NewScanner(os.Stdin)
	ev := evaluator.NewEvaluator()
	ev.Emit = printValue
	fmt.Print("> ")

	for reader.Scan() {
		scn := scanner.NewScanner(reader.Text())
		tokens := scn.ScanTokens()

		if !scn.ErrorFound() {
			prs := parser.NewParser(tokens)
			var err error
			// tokens always ends with EOF, so tokens[n-2] is the last real token.
			if n := len(tokens); n > 1 && (tokens[n-2].Type == scanner.TOKEN_SEMI_COLON || tokens[n-2].Type == scanner.TOKEN_RIGHT_BRACE) {
				if script := prs.ParseScript(); !prs.ErrorFound() {
					err = ev.ExecuteREPL(script)
				}
			} else if program := prs.Parse(); !prs.ErrorFound() {
				// No trailing ; or }: Lab 3 expression path, so `2 + 3` still prints 5.
				_, err = ev.Evaluate(program)
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
		fmt.Print("> ")
	}

	if err := reader.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: reading input: %v\n", err)
		os.Exit(1)
	}
}
