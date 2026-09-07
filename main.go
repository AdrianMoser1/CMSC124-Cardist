package main

import (
    "fmt"
    "cardist/scanner"
)

func main() {
    code := "({\n+ - = \n})"
    scn := scanner.NewScanner(code)
    tokens := scn.ScanTokens()

    for _, tok := range tokens {
        fmt.Println(tok)
    }
}