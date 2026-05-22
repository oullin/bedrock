#!/usr/bin/env bash
#
# Pretty, streaming front-end for `make format` and `make format-all`.
#
# Default mode runs `go-fmt` per Go module in parallel (see FORMAT_PARALLELISM)
# so each module produces a one-line progress entry as it finishes. The
# upstream `go-fmt` binary buffers all of its own output to the end of each
# invocation, so per-module invocation is the only way to get real streaming.
#
# Env:
#   FORMAT_BASE          git ref to diff against (default: origin/main)
#   FORMAT_PARALLELISM   concurrent go-fmt workers (default: 6)
#   FORMAT_BATCH=1       opt out of streaming; one batched go-fmt call (fastest
#                        wall time, no per-module progress — good for CI logs)
#
# Usage: format.sh <changed|all>

set -u -o pipefail

MODE="${1:-changed}"
case "$MODE" in
    changed|all) ;;
    *) echo "usage: format.sh <changed|all>" >&2; exit 2 ;;
esac

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$ROOT"

COMPOSE_FILE="go-fmt.compose.yaml"
GO_FMT_SERVICE="go-fmt"
GO_FMT_BIN="/usr/local/bin/go-fmt"
FORMAT_BASE="${FORMAT_BASE:-origin/main}"
FORMAT_PARALLELISM="${FORMAT_PARALLELISM:-6}"
FORMAT_BATCH="${FORMAT_BATCH:-0}"

# ── colors / pretty helpers ──────────────────────────────────────────────
if [[ -t 1 ]]; then
    BOLD=$'\033[1m'; DIM=$'\033[2m'; RESET=$'\033[0m'
    CYAN=$'\033[36m'; GREEN=$'\033[32m'; YELLOW=$'\033[33m'
    RED=$'\033[31m'; BLUE=$'\033[34m'; MAGENTA=$'\033[35m'
    export FORCE_COLOR=1
else
    BOLD=""; DIM=""; RESET=""
    CYAN=""; GREEN=""; YELLOW=""; RED=""; BLUE=""; MAGENTA=""
fi
export ROOT COMPOSE_FILE GO_FMT_SERVICE GO_FMT_BIN
export YELLOW GREEN RED DIM RESET

hr()    { printf "%s%s%s\n" "$DIM" "────────────────────────────────────────────────────────────" "$RESET"; }
stage() { printf "\n%s▶ %s%s\n" "$BOLD$CYAN" "$1" "$RESET"; }
ok()    { printf "%s✔ %s%s%s\n" "$GREEN" "$1" "${2:+ ${DIM}($2)${RESET}${GREEN}}" "$RESET"; }
warn()  { printf "%s! %s%s\n"   "$YELLOW" "$1" "$RESET"; }
fail()  { printf "%s✘ %s%s\n"   "$RED" "$1" "$RESET" >&2; }

human_secs() {
    local s=$1
    if   (( s < 60 ));   then printf "%ds" "$s"
    else                       printf "%dm%02ds" $((s/60)) $((s%60)); fi
}

# Prefix every line of stdin with a fixed-width colored tag, flushing per line.
prefix() {
    local tag="$1" color="$2"
    awk -v p="${color}$(printf '%-9s' "[$tag]")${RESET} " '{ print p $0; fflush(); }'
}

# ── change discovery ─────────────────────────────────────────────────────
list_changed() {
    local pattern="$1"
    {
        git diff --name-only --diff-filter=ACMRT "${FORMAT_BASE}"...HEAD -- "$pattern" 2>/dev/null || true
        git diff --name-only --diff-filter=ACMRT                       -- "$pattern" 2>/dev/null || true
        git ls-files --others --exclude-standard                       -- "$pattern" 2>/dev/null || true
    } | sort -u
}

changed_go_modules() {
    list_changed '*.go' | while IFS= read -r f; do
        d="$(dirname "$f")"
        while [[ "$d" != "." && ! -f "$d/go.mod" ]]; do d="$(dirname "$d")"; done
        [[ "$d" != "." ]] && echo "$d"
    done | sort -u
}

all_go_modules() {
    git ls-files 'packages/**/go.mod' 'services/**/go.mod' \
        | sed 's|/go.mod$||' \
        | grep -v '^packages/testing$' \
        | sort -u
}

# ── run ──────────────────────────────────────────────────────────────────
start_ts=$(date +%s)
hr
printf "%s%s format%s %s(mode: %s%s%s)%s\n" \
    "$BOLD" "$CYAN" "$RESET" "$DIM" "$RESET$CYAN" "$MODE" "$RESET$DIM" "$RESET"
hr

# Safety net: bring container up if not already running.
if ! docker compose -f "$COMPOSE_FILE" ps --status running --services 2>/dev/null \
        | grep -qx "$GO_FMT_SERVICE"; then
    stage "Bringing up go-fmt container"
    docker compose -f "$COMPOSE_FILE" up -d "$GO_FMT_SERVICE" 2>&1 | prefix docker "$MAGENTA"
fi

# Collect the Go module list up front so we can show the count before launching.
modules=()
if [[ "$MODE" == "all" ]]; then
    while IFS= read -r m; do modules+=("$m"); done < <(all_go_modules)
else
    while IFS= read -r m; do modules+=("$m"); done < <(changed_go_modules)
fi

# ── parallel: pnpm fmt + go-fmt ──────────────────────────────────────────
stage "Formatting source (pnpm fmt + go-fmt in parallel)"

pnpm_status=0
(
    set -o pipefail
    pnpm fmt 2>&1 | prefix pnpm "$BLUE"
) &
pnpm_pid=$!

go_status=0
if (( ${#modules[@]} == 0 )); then
    printf "  %sgo-fmt:%s no Go modules to format\n" "$DIM" "$RESET"
elif [[ "$FORMAT_BATCH" == "1" ]]; then
    printf "  %sgo-fmt:%s batched (%d module(s), no per-module stream)\n" \
        "$DIM" "$RESET" "${#modules[@]}"
    paths=()
    for m in "${modules[@]}"; do paths+=("/work/$m"); done
    if ! docker compose -f "$COMPOSE_FILE" exec -T "$GO_FMT_SERVICE" "$GO_FMT_BIN" \
            format --cwd /work "${paths[@]}" 2>&1 | prefix go-fmt "$YELLOW"; then
        go_status=1
    fi
else
    printf "  %sgo-fmt:%s %d module(s), %d parallel worker(s) — streaming as each finishes\n" \
        "$DIM" "$RESET" "${#modules[@]}" "$FORMAT_PARALLELISM"
    go_t0=$(date +%s)
    # xargs returns 123 if any child exits non-zero (we treat that as failure).
    if ! printf '%s\n' "${modules[@]}" \
            | xargs -n 1 -P "$FORMAT_PARALLELISM" -I{} \
                "$ROOT/services/scripts/format-go-one.sh" {}; then
        go_status=1
    fi
    go_dt=$(( $(date +%s) - go_t0 ))
    if (( go_status == 0 )); then
        printf "  %s↳ go-fmt: %d module(s) in %s%s\n" \
            "$DIM" "${#modules[@]}" "$(human_secs $go_dt)" "$RESET"
    fi
fi

wait "$pnpm_pid" || pnpm_status=$?

# ── markdown ─────────────────────────────────────────────────────────────
md_status=0
md_files=()
if [[ "$MODE" == "all" ]]; then
    while IFS= read -r f; do md_files+=("$f"); done < <(git ls-files '*.md')
else
    while IFS= read -r f; do md_files+=("$f"); done < <(list_changed '*.md')
fi

if (( ${#md_files[@]} > 0 )); then
    stage "Formatting markdown (${#md_files[@]} file(s))"
    if ! pnpm exec oxfmt --ignore-path .gitignore "${md_files[@]}" 2>&1 | prefix oxfmt "$MAGENTA"; then
        md_status=1
    fi
else
    stage "Markdown — none changed, skipping"
fi

# ── summary ──────────────────────────────────────────────────────────────
elapsed=$(( $(date +%s) - start_ts ))
hr
failed=0
(( pnpm_status )) && { fail "pnpm fmt failed (exit $pnpm_status)"; failed=1; }
(( go_status   )) && { fail "go-fmt failed";                       failed=1; }
(( md_status   )) && { fail "oxfmt failed";                        failed=1; }

if (( failed == 0 )); then
    ok "Format complete" "$(human_secs "$elapsed")"
    exit 0
else
    fail "Format failed in $(human_secs "$elapsed")"
    exit 1
fi
