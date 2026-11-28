#!/usr/bin/env bash
set -e
mkdir -p build
binary=build/interpreter
case "$OSTYPE" in
  msys*|cygwin*|win32*) binary="$binary.exe" ;;
esac
go build -o "$binary" .
