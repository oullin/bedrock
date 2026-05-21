#!/usr/bin/env bash
#
# Format a single Go module via `go-fmt` and print one tagged summary line.
# Used by format.sh as the per-module worker.
#
# Args:
#   $1  Module path relative to repo root (e.g. packages/auth).
#
# Env (set by parent):
#   ROOT, COMPOSE_FILE, GO_FMT_SERVICE, GO_FMT_BIN
#   YELLOW, GREEN, RED, DIM, RESET   (may be empty when not on a TTY)

set -u -o pipefail

mod="$1"

: "${ROOT:=$PWD}"
: "${COMPOSE_FILE:=go-fmt.compose.yaml}"
: "${GO_FMT_SERVICE:=go-fmt}"
: "${GO_FMT_BIN:=/usr/local/bin/go-fmt}"
: "${YELLOW:=}"; : "${GREEN:=}"; : "${RED:=}"; : "${DIM:=}"; : "${RESET:=}"

t0=$(date +%s)
if out=$(docker compose -f "$ROOT/$COMPOSE_FILE" exec -T "$GO_FMT_SERVICE" \
            "$GO_FMT_BIN" format -format json --host-path "$ROOT/$mod" 2>&1); then
    rc=0
else
    rc=$?
fi
dt=$(( $(date +%s) - t0 ))

files=$(   printf '%s' "$out" | tr ',{}' '\n\n\n' | sed -n 's/.*"files":\([0-9][0-9]*\).*/\1/p'      | head -1)
changed=$( printf '%s' "$out" | tr ',{}' '\n\n\n' | sed -n 's/.*"changed":\([0-9][0-9]*\).*/\1/p'    | head -1)
violations=$(printf '%s' "$out" | tr ',{}' '\n\n\n' | sed -n 's/.*"violations":\([0-9][0-9]*\).*/\1/p' | head -1)

# Single-line summary keeps writes < PIPE_BUF (4KB) → atomic under xargs -P.
if (( rc == 0 )); then
    extra=""
    if [[ -n "${changed:-}" && "$changed" != "0" ]]; then extra="${extra}  ${changed} changed"; fi
    if [[ -n "${violations:-}" && "$violations" != "0" ]]; then extra="${extra}  ${violations} violation(s)"; fi
    printf "%s[go-fmt]%s %-46s %s%3ss%s  %4s files%s%s\n" \
        "$YELLOW" "$RESET" "$mod" "$DIM" "$dt" "$RESET" "${files:-?}" "$extra" ""
else
    {
        printf "%s[go-fmt]%s %-46s %sFAILED in %ss%s\n" \
            "$RED" "$RESET" "$mod" "$RED" "$dt" "$RESET"
        printf '%s\n' "$out" | sed "s/^/  /"
    } >&2
fi
exit "$rc"
