#!/usr/bin/env bash
set -euo pipefail

scratch=$(mktemp -d)
trap 'rm -f "$scratch/input.crd" "$scratch/expected" "$scratch/actual"; rmdir "$scratch"' EXIT

# The harness trims newlines, so compare the bytes from the public launcher.
printf '""\n' > "$scratch/input.crd"
printf '\n' > "$scratch/expected"
./run --parse "$scratch/input.crd" > "$scratch/actual"
cmp "$scratch/expected" "$scratch/actual"

printf '// No expressions.\n' > "$scratch/input.crd"
./run --parse "$scratch/input.crd" > "$scratch/actual"
test ! -s "$scratch/actual"
