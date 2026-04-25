#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "usage: scripts/publish-main.sh \"commit message\"" >&2
  exit 1
fi

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

msg="$1"

git add .
git commit -m "$msg"
git push origin HEAD:main

echo "Pushed to origin/main"
