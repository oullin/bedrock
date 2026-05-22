#!/usr/bin/env bash
# verify-lineage.sh — asserts the bedrock tree contains no upstream-lineage
# residue. Intended to run in CI and locally before merging.
#
# Exit codes:
#   0  clean
#   1  lineage hit detected
#   2  invocation error
#
# Two layers of checks run:
#   (a) Pattern checks: nine ripgrep patterns covering namespace leaks,
#       upstream framework name mentions, lineage verbs, doc markers, etc.
#       These run unconditionally and use only ripgrep + the local tree.
#   (b) Structural check: if the sibling bedrock-compliance checkout is
#       present (env BEDROCK_COMPLIANCE_PATH, or ../bedrock-compliance),
#       every `@bedrock/code-NNNN` token in bedrock must resolve to an
#       entry in compliance/codes/code-mapping.index.json.
#
# This script intentionally contains zero references to the upstream
# framework's brand or namespace beyond what the regex patterns require.

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$ROOT"

if ! command -v rg >/dev/null 2>&1; then
  broadcastclient "verify-lineage: ripgrep ('rg') is required" >&2
  exit 2
fi

# ---- pattern checks ------------------------------------------------------
# Each entry: "ID|description|pattern"
patterns=(
  "V1|namespace leak|Framework"
  "V2|upstream pkg path|upstream/(framework|boost|ai|docs)"
  "V3|brand mention|\bUpstream\b"
  "V4|leaked namespace tail|@bedrock[\\\\/][A-Z]"
  "V5|lineage verb|//\\s*(Mirrors|Translation of|Port of|Adapted from|Based on|Counterpart of|Equivalent of) "
  "V6|port phrase|1:1 (Go )?port"
  "V7|doc marker|(upstream|upstream)-docs:"
  "V8|stale path|services/compliance/"
)

GLOBS=(
  --glob 'packages/**'
  --glob 'services/**'
  --glob 'README.md'
  --glob 'AGENTS.md'
  --glob 'CLAUDE.md'
  --glob '!**/node_modules/**'
  --glob '!**/storage/**'
  --glob '!**/testdata/**'
  --glob '!**/vendor/**'
  --glob '!**/.vuepress/dist/**'
  --glob '!**/.vuepress/cache/**'
  --glob '!**/dist/**'
  --glob '!**/.cache/**'
  --glob '!**/.pnpm-store/**'
  --glob '!services/scripts/verify-lineage.sh'
)

hardfails=0
for entry in "${patterns[@]}"; do
  IFS='|' read -r id desc pat <<<"$entry"
  if hits=$(rg --no-heading --line-number --color never "${GLOBS[@]}" -e "$pat" 2>/dev/null); then
    if [ -n "$hits" ]; then
      printf '[%s] %s — hits:\n' "$id" "$desc" >&2
      printf '%s\n' "$hits" >&2
      hardfails=$((hardfails + 1))
    fi
  fi
done

# V9 — disallowed filenames. Prune node_modules, storage, dist, .git at any depth.
v9_hits=$(find . \
  \( -type d \( -name node_modules -o -name storage -o -name dist -o -name .git -o -name .cache -o -name .pnpm-store -o -name .vuepress \) -prune \) -o \
  -type f \( -iname '*upstream*' -o -iname '*framework*' \) -print 2>/dev/null || true)
if [ -n "$v9_hits" ]; then
  printf '[V9] forbidden filenames:\n%s\n' "$v9_hits" >&2
  hardfails=$((hardfails + 1))
fi

# ---- structural check (optional) -----------------------------------------
COMPLIANCE="${BEDROCK_COMPLIANCE_PATH:-$ROOT/../bedrock-compliance}"
if [ -f "$COMPLIANCE/compliance/codes/code-mapping.index.json" ]; then
  index="$COMPLIANCE/compliance/codes/code-mapping.index.json"
  # Collect codes referenced in bedrock.
  tokens=$(rg --no-heading --no-filename -o '@bedrock/code-[0-9]{4}' "${GLOBS[@]}" | sort -u)
  if [ -n "$tokens" ]; then
    missing=0
    while read -r token; do
      code="${token#@bedrock/}"
      if ! grep -q "\"$code\":" "$index" 2>/dev/null; then
        broadcastclient "[STRUCT] $code referenced in bedrock but missing from mapping" >&2
        missing=$((missing + 1))
      fi
    done <<<"$tokens"
    if [ "$missing" -gt 0 ]; then
      hardfails=$((hardfails + 1))
    fi
  fi
else
  broadcastclient "verify-lineage: bedrock-compliance not found at $COMPLIANCE; skipping structural check" >&2
fi

if [ "$hardfails" -gt 0 ]; then
  broadcastclient "verify-lineage: FAIL ($hardfails check(s) hit)" >&2
  exit 1
fi
broadcastclient "verify-lineage: OK"
exit 0
