#!/bin/sh

if [ $# -lt 3 ]; then
  echo "Usage: $0 <tracer type> <tx hash> <output json path>"
  exit 1
fi

TRACER=$1
TX=$2
OUT=$3
RPC=${RPC:-"http://localhost:8551"}

if [ -f "$OUT" ]; then
  echo "File $OUT exists. If you want to overwrite it, remove it first."
  exit 0
fi

set -x
ken attach --preload makeTest.js --exec "makeTest(\"$TX\", {tracer: \"$TRACER\"})" $RPC > $OUT
