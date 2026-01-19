#!/usr/bin/env bash
set -euo pipefail

if [ "$#" -lt 2 ]; then
  echo "usage: $0 <os> <component>..." >&2
  exit 1
fi

python3 - "$@" <<'PY'
import ntpath
import posixpath
import sys

args = sys.argv[1:]
if len(args) < 2:
    sys.stderr.write("need target os and at least one path component\n")
    sys.exit(1)

target = args[0].lower()
parts = args[1:]
lib = ntpath if target.startswith("win") else posixpath
joined = lib.join(*parts)
sys.stdout.write(lib.normpath(joined))
PY
