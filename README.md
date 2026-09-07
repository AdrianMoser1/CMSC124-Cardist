# Cardist

## Creators

- Adrian Moser (AdrianMoser1)
- Lunele Moscos (LuneleIven)

## Overview

Cardist is a small, dynamically typed scripting language for describing turn-based card game rules, cards, effects, costs, and turn sequencing. It's built for someone designing a card game who wants to write "what a card does" as readable code with corresponding syntax rather than a giant switch statement in a general-purpose language. It simplifies the core idea of PVE card game building: it shouldn’t be too complicated to make one, especially for aspiring game designers. 

## Host language and build

- Host language: Go 1.26.7
- Version metadata: go.mod
- Build: `./build.sh`
- Fresh clone needs Go 1.22+ on PATH; no other dependencies.

## Running it

| Command | What it does |
|---|---|
| `./run <file>` | Executes a program. Available from Lab 4. |
| `./run --tokenize <file>` | Prints the token stream. |
| `./run --parse <file>` | Not yet implemented — Lab 2. |
| `./run --eval <file>` | Not yet implemented — Lab 3. |
| `./run` | Starts the REPL; tokenizes one line at a time as of Lab 1. |

| Exit code | Meaning |
|---|---|
| 0 | file scans cleanly |
| 65 | scanner rejects the file (unterminated string, illegal character) |
| 70 | reserved for runtime errors starting Lab 3 |

## File extension

`.deck` — must match the `ext` field in every `tests/lab*/manifest.json`.

## Lexical structure

### Keywords

| Keyword | Purpose |
|---|---|
| var | declares a variable |
| if | conditional branch |
| else | alternate branch |
| while | loop |
| true | boolean literal |
| false | boolean literal |
| nil | absence of a value |
| and | logical AND |
| or | logical OR |
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
| `!` | logical negation | unary | right | 7 |
| `-` | arithmetic negation | unary | right | 7 |


### Literals

| Kind | Syntax | Produces |
|---|---|---|
| number | `6`, `1.5` | Go float64 |
| string | `"Strike"`, no escapes yet | Go string |
| boolean | `true`, `false` | Go bool |
| nil | `nil` | Go nil interface |

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
Fields, in order: token type, the raw source text, the literal value (`null` for non-literal tokens), and the 1-indexed source line. 

## Errors and diagnostics
Message format: 
```
[line 4] Error: Unterminated string.
[line 9] Error: Unexpected character '$'. 
```
| Failure | Exit code |
|---|---|
| lexical error (unterminated string, illegal character) | 65 |
| syntax error | 65 (Lab 2) |
| runtime error | 70 (Lab 3) |

## Design rationale 
Keywords split into two tiers: general-purpose control flow (`var`, `if`, `while`) and combat-domain vocabulary (`card`, `enemy`, `artifact`, `intent`, `deal`, `apply`). The domain tier is deliberately close to how PVE card games UI describes things — "deal 6 damage," "apply 2 weak"  so a rules author's script reads like a card's actual tooltip. `intent` was included to dictate what would the enemy do in the next turn given a unique cyclical actions per enemy, it's cheap to reserve as a keyword now and expensive to retrofit into existing tests once card scripts already use `intent` as a bare identifier. `elixir` and `artifact` share most of ‘card’'s shape (a name, a cost or trigger condition, and an effect block) but are kept as separate keywords rather than folded into one generic `item` block, since the design would treat trigger timing differently enough that conflating them now would make Lab 2's grammar harder to write correctly. 

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
Token(type=IDENTIFIER, lexeme=effect, literal=null, line=3)
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
- No multi-line strings — a newline inside a string is a scan error.
- No block comments.
- `enemy`/`intent`/`artifact`/`elixir` are reserved keywords but their
  runtime semantics will be implemented later on

## Changelog

| Activity | What changed in the language |
|---|---|
| Lab 1 | Scanner, initial keyword list, token format frozen. |