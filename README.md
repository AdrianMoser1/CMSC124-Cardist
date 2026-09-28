
# Cardist

## Creators

- Adrian Moser (AdrianMoser1)
- Lunele Moscos (LuneleIven)

## Overview

Cardist is a small, dynamically typed scripting language for describing turn-based card game rules, cards, effects, costs, and turn sequencing. It's built for someone designing a card game who wants to write "what a card does" as readable code with corresponding syntax rather than a giant switch statement in a general-purpose language. It simplifies the core idea of PVE card game building: it shouldn't be too complicated to make one, especially for aspiring game designers.

## Host language and build

- Host language: Go (version pinned in `go.mod`, currently 1.27.0)
- Version metadata: go.mod
- Build: `./build.sh`
A fresh clone needs the Go version in `go.mod` (or newer) on PATH; no other dependencies.
## Running it

| Command | What it does |
|---|---|
| `./run <file>` | Lab 0 placeholder: prints the file back unchanged. Becomes program execution in Lab 4. |
| `./run --tokenize <file>` | Scans the file and prints its token stream. |
| `./run` | Starts the REPL: scans one line at a time and prints that line's tokens. A bad line prints its error and returns the prompt. |

| Exit code | Meaning |
|---|---|
| 0 | clean run (token stream printed, or file echoed) |
| 64 | invalid command-line usage |
| 65 | scanner rejected the input (unterminated string, illegal character) |
| 66 | file could not be read/opened |
| 70 | reserved for a run that starts and then dies partway (runtime error); unused until Lab 3 |

## File extension

`.crd` — must match the `ext` field in every `tests/lab*/manifest.json`.

## Lexical structure
### Token list

Every token type the scanner can emit, listed once. Names are the `TOKEN_` constants in `scanner/scanner.go` minus the prefix.

| Category | Token types | Why it earns its place |
|---|---|---|
| Grouping | `LEFT_PAREN` `RIGHT_PAREN` | expression grouping and (later) call arguments |
| Blocks | `LEFT_BRACE` `RIGHT_BRACE` | block delimiters; braces do the scoping since whitespace is not significant |
| Arithmetic | `PLUS` `MINUS` `STAR` `SLASH` | the four arithmetic operators |
| Comparison / logic | `BANG` `BANG_EQUAL` `EQUAL_EQUAL` `GREATER` `GREATER_EQUAL` `LESSER` `LESSER_EQUAL` | comparisons and logical negation |
| Assignment | `EQUAL` | assignment, and `cost = 1` style declarations |
| Punctuation | `COMMA` `COLON` `SEMI_COLON` | `COMMA` for argument/parameter lists; `COLON` and `SEMI_COLON` are scanned and reserved (no grammar rule uses them yet) |
| Literals | `IDENTIFIER` `STRING` `NUMBER` | names, text, and numbers |
| General keywords | `VAR` `IF` `ELSE` `ELSEIF` `SWITCH` `WHILE` `FOR` `BREAK` `CONTINUE` `FUNC` `RETURN` `PRINT` `TRUE` `FALSE` `NIL` `AND` `OR` | general-purpose control flow and values |
| Domain keywords | `CARD` `ENERGY` `COST` `ENEMY` `INTENT` `ARTIFACT` `ELIXIR` `DEAL` `BLOCK` `APPLY` `DRAW` `BANISH` `TURN` `PLAYER` `EFFECT` | combat vocabulary (see design rationale) |
| End | `EOF` | tells the parser where the input stops |

### Keywords

| Keyword | Purpose |
|---|---|
| var | declares a variable |
| if | conditional branch |
| elseif | n or more alternate branch |
| else | last alternate branch |
| switch | multiple alternate branching |
| while | loop |
| true | boolean literal |
| false | boolean literal |
| nil | absence of a value |
| and | logical AND |
| or | logical OR |
| func | identify to be a function that can be called |
| print | prints for user view |
| return | return a value when from function |
| for | predetermined looping |
| break | disrupt the flow of iteration |
| continue | skip the current iteration |
| card | declares a card definition block |
| energy | refers to current energy pool |
| cost | declares a card's energy cost |
| enemy | declares an enemy definition block |
| intent | declares an enemy's telegraphed next action |
| artifact | declares an artifact definition block |
| elixir | declares an elixir definition block |
| deal | deals damage |
| block | grants block (damage mitigation) |
| apply | applies a status effect (e.g. weak, vulnerable, poison) |
| draw | draws a card |
| banish | removes a card from the fight entirely |
| turn | marks a turn-scoped block |
| player | refers to the player entity |
| effect | works alongside with apply to put a status in an entity |

### Operators

| Operator | Category | Operands | Associativity | Precedence |
|---|---|---|---|---|
| `=` | assignment | binary | right | 1 |
| `or` | logical | binary | left | 2 |
| `and` | logical | binary | left | 3 |
| `==`, `!=` | comparison | binary | left | 4 |
| `<`, `<=`, `>`, `>=` | comparison | binary | left | 5 |
| `+`, `-` | arithmetic | binary | left | 6 |
| `*`, `/` | arithmetic | binary | left | 7 |
| `!` | logical negation | unary | right | 8 |
| `-` | arithmetic negation | unary | right | 8 |

### Punctuation

| Character | Token | Use |
|---|---|---|
| `(` `)` | LEFT_PAREN, RIGHT_PAREN | grouping |
| `{` `}` | LEFT_BRACE, RIGHT_BRACE | blocks |
| `,` | COMMA | separator in lists |
| `:` | COLON | reserved |
| `;` | SEMI_COLON | reserved |

### Literals

| Kind | Syntax | Produces |
|---|---|---|
| number | `6`, `1.5` | Go float64 |
| string | `"Strike"`, single line, no escapes yet | Go string |
| boolean | `true`, `false` | Go bool |
| nil | `nil` | Go nil interface |

- Numbers: digits, optionally `.` followed by at least one digit. `.5` is not a number (the `.` is an illegal character) and `5.` scans as `5` followed by an illegal `.`.
- Strings: must close on the same line. A newline before the closing quote is an `Unterminated string.` error. The lexeme keeps the quotes; the literal does not.

### Identifiers

- Start characters: letters, `_`
- Continue characters: letters, digits, `_`
- Case-sensitive: yes
- No reserved-word prefixing restriction — `energy_gained` and `intent_next`
  scan as identifiers, not as keyword-plus-suffix.

### Comments

- Line comments: `//`
- Block comments: not supported
- Nesting: not applicable
- Harness note: `comment_prefix` in `tests/lab1/manifest.json` is set
  to `//`.

## Whitespace and termination

- Whitespace significant: no, beyond separating tokens
- Statement terminator: none (newline-agnostic; block scoping does the work)
- Block delimiters: braces `{ }`
- Grouping delimiters: parentheses `( )`

## Token output format
```
Token(type=CARD, lexeme=card, literal=null, line=1)
```
Fields, in order: token type, the raw source text, the literal value (`null` for non-literal tokens), and the 1-indexed source line. String lexemes include their quotes and their literals do not; number literals print as Go prints a float64 (`1`, `1.5`).

## Errors and diagnostics
Message format (on stderr):
```
[line 4] Error: Unterminated string.
[line 9] Error: Unexpected character '$'.
```
Scanning continues after an error so several problems are reported in one pass. A multi-byte character (for example `©`) is reported once, as itself.

| Failure | Exit code |
|---|---|
| lexical error (unterminated string, illegal character) | 65 |
| syntax error | 65 (Lab 2) |
| runtime error | 70 (Lab 3) |

## Design rationale
Keywords split into two tiers: general-purpose control flow (`var`, `if`, `while`) and combat-domain vocabulary (`card`, `enemy`, `artifact`, `intent`, `deal`, `apply`). The domain tier is deliberately close to how PVE card games UI describes things — "deal 6 damage," "apply 2 weak"  so a rules author's script reads like a card's actual tooltip. `intent` was included to dictate what would the enemy do in the next turn given a unique cyclical actions per enemy, it's cheap to reserve as a keyword now and expensive to retrofit into existing tests once card scripts already use `intent` as a bare identifier. `elixir` and `artifact` share most of ‘card’'s shape (a name, a cost or trigger condition, and an effect block) but are kept as separate keywords rather than folded into one generic `item` block, since the design would treat trigger timing differently enough that conflating them now would make Lab 2's grammar harder to write correctly.

Strings are single-line so an unterminated string is caught at the end of its own line rather than swallowing the rest of the file, which keeps error line numbers honest.

## Testing conventions

| Folder | Activity | Mode | Flag |
|---|---|---|---|
| tests/lab1 | Scanner | sidecar | `--tokenize` |
| tests/lab2 | Parser | sidecar | `--parse` |
| tests/lab3 | Evaluator | inline | `--eval` |
| tests/lab4 | Context | inline | none |
| tests/lab5 | Functions | inline | none |

Run locally with:

```bash
curl -sSL https://raw.githubusercontent.com/WhiteLicorice/cmsc-124-harness/v1.1/run_tests.py -o run_tests.py
./build.sh
python3 run_tests.py tests/lab1
```

## Sample code

```
card "Strike" {
  cost = 1
  effect {
    deal 6 to enemy
  }
}
```

(A more elaborate `enemy "Cultist" { intent { ... } }` block scans the
same way — `intent`, `if`, `turn`, `player`, and comparison operators
all resolve to their own token types.)

Output (`--tokenize` on just the "Strike" card above):

```
Token(type=CARD, lexeme=card, literal=null, line=1)
Token(type=STRING, lexeme="Strike", literal=Strike, line=1)
Token(type=LEFT_BRACE, lexeme={, literal=null, line=1)
Token(type=COST, lexeme=cost, literal=null, line=2)
Token(type=EQUAL, lexeme==, literal=null, line=2)
Token(type=NUMBER, lexeme=1, literal=1, line=2)
Token(type=EFFECT, lexeme=effect, literal=null, line=3)
Token(type=LEFT_BRACE, lexeme={, literal=null, line=3)
Token(type=DEAL, lexeme=deal, literal=null, line=4)
Token(type=NUMBER, lexeme=6, literal=6, line=4)
Token(type=IDENTIFIER, lexeme=to, literal=null, line=4)
Token(type=ENEMY, lexeme=enemy, literal=null, line=4)
Token(type=RIGHT_BRACE, lexeme=}, literal=null, line=5)
Token(type=RIGHT_BRACE, lexeme=}, literal=null, line=6)
Token(type=EOF, lexeme=, literal=null, line=7)
```

## Known limitations

- No string escape sequences.
- Strings cannot span lines.
- Decimals must have a present whole number followed by '.' then one digit after it(Defense mechanism; poka-yoke)
- If a file has any lexical error, no tokens are printed for it: the whole file is rejected ("nothing about a rejected file belongs on stdout").
- No block comments.
- `enemy`/`intent`/`artifact`/`elixir` are reserved keywords but their
  runtime semantics will be implemented later on.

## Changelog

| Activity | What changed in the language |
|---|---|
| Lab 1 | Scanner, initial keyword list, token format frozen. |

## Bug log

| Bug | What broke | Fix |
|---|---|---|
| Rejected files leaked partial output | Tokens printed to stdout as scanning happened, so a rejected file still showed partial tokens before the error — violated "nothing about a rejected file belongs on stdout" | Buffer all tokens, only flush to stdout if the scan succeeded |
| `HadError`/`ErrorFound` unexported | Method named `errorFound()` (lowercase) — the `main` package couldn't call it; compile error | Capitalized to `ErrorFound()` |
| Missing `reader.Err()` check | REPL's `bufio.Scanner` loop never distinguished normal EOF from an actual stdin read failure | Added the check; exit code `1` on a genuine read error |
| `TokenType.String()` missing 8 cases | `print`, `return`, `for`, `function`/`func`, `break`, `continue`, `elseif`, `switch` all printed `type=UNKNOWN` despite being correctly recognized as keywords internally | Added the missing `case` statements |
| `operators.expected` CRLF mismatch | `.expected` file saved with `\r\n`; the program only ever emits `\n` — byte-for-byte mismatch in the test harness | Re-saved LF-only; added `.gitattributes` (`* text=auto eol=lf`) so it can't recur |
| `tests/lab0` regression | Rewriting `main.go` for `--tokenize` accidentally hijacked plain `./run <file>` to always reject with "not implemented until Lab 4," breaking Lab 0's original echo-the-file contract | Split into `runEcho()` (Lab 0's placeholder) vs `runFile()` (Lab 1's real scanning) — both now coexist |
| Keywords table broken Markdown | Table split into two blocks mid-document with no header/separator on the second half, plus a duplicated `or` row | Merged into one continuous table |
| README/scanner drift on strings | README said strings can't span lines; scanner allowed it and incremented the line counter in two places | Scanner now rejects a newline inside a string; only `scanToken` counts lines |
| Error message drift | README showed `Unexpected character '$'.`; scanner printed `Unexpected character: $` | Scanner message now matches the README; multi-byte characters reported once |
| REPL printed the wrong thing for Lab 1 | REPL printed parse trees instead of tokens | `replParses` switch in `main.go`; Lab 1 setting prints tokens |