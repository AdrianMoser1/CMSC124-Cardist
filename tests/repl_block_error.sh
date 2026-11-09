#!/usr/bin/env bash
# The harness only runs files, so pipe into the REPL. If executeBlock didn't
# restore the environment, `print b;` would find the dead block's b and print 2.
out=$(printf 'var a = 1;\n{ var b = 2; print b / 0; }\nprint b;\nprint a;\n' | ./run 2>/tmp/repl_err | sed 's/> //g')
[[ "$out" == "1" ]] && grep -q "Undefined variable 'b'." /tmp/repl_err && echo "PASS repl_block_error" || { echo "FAIL repl_block_error: got '$out'"; exit 1; }