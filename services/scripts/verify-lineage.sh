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
#   (a) Pattern checks: ripgrep patterns sourced from
#       compliance/codes/scrub-terms.txt (one rule per line: ID|desc|pattern).
#       The pattern table is intentionally externalised so this script body
#       contains no upstream-brand tokens.
#   (b) Structural check: if the sibling bedrock-compliance checkout is
#       present (env BEDROCK_COMPLIANCE_PATH, or ../bedrock-compliance),
#       every `@bedrock/code-NNNN` token in bedrock must resolve to an
#       entry in compliance/codes/code-mapping.index.json.

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$ROOT"

if ! command -v rg >/dev/null 2>&1; then
  echo "verify-lineage: ripgrep ('rg') is required" >&2
  exit 2
fi

PATTERN_FILE="${BEDROCK_SCRUB_TERMS:-$ROOT/compliance/codes/scrub-terms.txt}"
if [ ! -f "$PATTERN_FILE" ]; then
  echo "verify-lineage: pattern file not found at $PATTERN_FILE" >&2
  exit 2
fi

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
  --glob '!compliance/codes/scrub-terms.txt'
  --glob '!compliance/**'
  --glob '!pnpm-lock.yaml'
  --glob '!**/pnpm-lock.yaml'
)

# Per-pattern allowlist: an ID followed by a comma-separated list of file
# globs whose hits we ignore. This lets us tolerate legitimate occurrences
# of an otherwise-banned token (e.g. an HTTP wire-protocol header, a real
# npm package name, or a fixed third-party brand inside a code comment that
# clearly does not assert lineage).
#
# Format: ID|glob1,glob2,...
allowlist=(
  # The Precognition HTTP wire-protocol header is part of the published
  # spec consumed by the laravel-precognition npm client. Keeping the
  # header name is required for interop.
  "V3|packages/inertia/**,packages/httpx/**"
)

hardfails=0
while IFS='|' read -r id desc pat; do
  case "$id" in
    ''|\#*) continue ;;
  esac
  [ -n "${pat:-}" ] || continue

  # Build per-pattern --glob exclusions from the allowlist.
  extra_globs=()
  for entry in "${allowlist[@]}"; do
    IFS='|' read -r aid aglobs <<<"$entry"
    if [ "$aid" = "$id" ]; then
      IFS=',' read -r -a parts <<<"$aglobs"
      for g in "${parts[@]}"; do
        extra_globs+=( --glob "!$g" )
      done
    fi
  done

  if hits=$(rg --no-heading --line-number --color never "${GLOBS[@]}" "${extra_globs[@]+"${extra_globs[@]}"}" -e "$pat" 2>/dev/null); then
    if [ -n "$hits" ]; then
      printf '[%s] %s — hits:\n' "$id" "$desc" >&2
      printf '%s\n' "$hits" >&2
      hardfails=$((hardfails + 1))
    fi
  fi
done <"$PATTERN_FILE"

# V9 — disallowed filenames. Prune large/cache dirs at any depth. The
# pattern list here is intentionally a regex sourced from the same file
# so the script body stays free of upstream tokens.
v9_iname_args=()
while IFS='|' read -r id desc pat; do
  case "$id" in
    ''|\#*) continue ;;
  esac
  if [ "$id" = "V1" ] || [ "$id" = "V3" ]; then
    # Convert the regex token (e.g. "Illuminate" or "\bLaravel\b") into
    # a basename glob by stripping leading/trailing anchors.
    tok="${pat//\\b/}"
    tok="${tok//\\/}"
    if [ -n "$tok" ] && [[ "$tok" != *"|"* && "$tok" != *"("* ]]; then
      v9_iname_args+=( -iname "*${tok}*" -o )
    fi
  fi
done <"$PATTERN_FILE"

if [ ${#v9_iname_args[@]} -gt 0 ]; then
  # Drop trailing -o
  unset 'v9_iname_args[${#v9_iname_args[@]}-1]'
  v9_hits=$(find . \
    \( -type d \( -name node_modules -o -name storage -o -name dist -o -name .git -o -name .cache -o -name .pnpm-store -o -name .vuepress -o -name compliance \) -prune \) -o \
    -type f \( "${v9_iname_args[@]}" \) -print 2>/dev/null || true)
  if [ -n "$v9_hits" ]; then
    printf '[V9] forbidden filenames:\n%s\n' "$v9_hits" >&2
    hardfails=$((hardfails + 1))
  fi
fi

# ---- structural check (optional) -----------------------------------------
COMPLIANCE="${BEDROCK_COMPLIANCE_PATH:-$ROOT/../bedrock-compliance}"
if [ -f "$COMPLIANCE/compliance/codes/code-mapping.index.json" ]; then
  index="$COMPLIANCE/compliance/codes/code-mapping.index.json"
  tokens=$(rg --no-heading --no-filename -o '@bedrock/code-[0-9]{4}' "${GLOBS[@]}" | sort -u)
  if [ -n "$tokens" ]; then
    missing=0
    while read -r token; do
      code="${token#@bedrock/}"
      if ! grep -q "\"$code\":" "$index" 2>/dev/null; then
        echo "[STRUCT] $code referenced in bedrock but missing from mapping" >&2
        missing=$((missing + 1))
      fi
    done <<<"$tokens"
    if [ "$missing" -gt 0 ]; then
      hardfails=$((hardfails + 1))
    fi
  fi
else
  echo "verify-lineage: bedrock-compliance not found at $COMPLIANCE; skipping structural check" >&2
fi

if [ "$hardfails" -gt 0 ]; then
  echo "verify-lineage: FAIL ($hardfails check(s) hit)" >&2
  exit 1
fi
echo "verify-lineage: OK"
exit 0
