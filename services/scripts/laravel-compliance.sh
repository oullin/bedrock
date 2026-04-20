#!/usr/bin/env bash
set -euo pipefail

ROOT_PATH="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
COMPLIANCE_PATH="$ROOT_PATH/services/compliance"
INVENTORY_FILE="$COMPLIANCE_PATH/inventory.yml"
DIVERGENCES_FILE="$COMPLIANCE_PATH/divergences.yml"
SOURCES_LOCK_FILE="$COMPLIANCE_PATH/sources.lock.json"
REPORT_FILE="$COMPLIANCE_PATH/report.md"

usage() {
  cat <<'USAGE'
Usage: services/scripts/laravel-compliance.sh <command>

Commands:
  refresh   Clone configured upstream sources and regenerate inventory files.
  check     Validate compliance source-of-truth rules.
  report    Generate services/compliance/report.md from current inventories.
USAGE
}

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "laravel-compliance: missing required command: $1" >&2
    exit 1
  fi
}

list_records() {
  awk '
    function trim(value) {
      gsub(/^[ \t]+|[ \t]+$/, "", value)
      return value
    }
    function clean(value) {
      value = trim(value)
      gsub(/^"/, "", value)
      gsub(/"$/, "", value)
      return value
    }
    function reset_record() {
      id = ""; status = ""; repo = ""; branch = ""; tests_path = "";
      inventory = ""; filter = ""; laravel = ""; bedrock = "";
    }
    function emit() {
      if (id != "") {
        print id "\t" status "\t" repo "\t" branch "\t" tests_path "\t" inventory "\t" filter "\t" laravel "\t" bedrock
      }
    }
    /^  - id:/ {
      emit()
      reset_record()
      sub(/^  - id:[ \t]*/, "")
      id = clean($0)
      next
    }
    id != "" && /^    status:/ {
      sub(/^    status:[ \t]*/, "")
      status = clean($0)
      next
    }
    id != "" && /^    repo:/ {
      sub(/^    repo:[ \t]*/, "")
      repo = clean($0)
      next
    }
    id != "" && /^    branch:/ {
      sub(/^    branch:[ \t]*/, "")
      branch = clean($0)
      next
    }
    id != "" && /^    tests_path:/ {
      sub(/^    tests_path:[ \t]*/, "")
      tests_path = clean($0)
      next
    }
    id != "" && /^    inventory:/ {
      sub(/^    inventory:[ \t]*/, "")
      inventory = clean($0)
      next
    }
    id != "" && /^    filter:/ {
      sub(/^    filter:[ \t]*/, "")
      filter = clean($0)
      next
    }
    id != "" && /^    laravel:/ {
      sub(/^    laravel:[ \t]*/, "")
      laravel = clean($0)
      next
    }
    id != "" && /^    bedrock:/ {
      sub(/^    bedrock:[ \t]*/, "")
      bedrock = clean($0)
      next
    }
    END {
      emit()
    }
  ' "$INVENTORY_FILE"
}

php_test_methods() {
  local file="$1"

  perl -ne '
    if (/\#\[Test\]/) {
      $pending_test_attribute = 1;
    }

    if (/function\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(/) {
      $name = $1;

      if ($name =~ /^test/ || $pending_test_attribute) {
        print "$name\n";
      }

      $pending_test_attribute = 0;
    }
  ' "$file"
}

generate_inventory() {
  local repo="$1"
  local branch="$2"
  local tests_path="$3"
  local filter="$4"
  local source_path="$5"
  local output_path="$6"

  mkdir -p "$(dirname "$output_path")"

  {
    echo "# Laravel compliance test inventory."
    echo "# Source: https://github.com/$repo/tree/$branch/$tests_path"
    echo "# Format: <file>::<test method name>"
    echo "# Generated: $(date -u +%Y-%m-%d)"
    echo "#"
    echo "# Every entry below must have a Go counterpart, an adaptation entry, or an"
    echo "# explicit missing status in services/compliance/report.md."
    echo

    if [ ! -d "$source_path" ]; then
      return 0
    fi

    find "$source_path" -type f -name '*Test.php' | sort | while IFS= read -r php_file; do
      local base
      base="$(basename "$php_file")"

      if [ -n "$filter" ] && [ "$filter" != "null" ] && [ "$base" != "$filter" ]; then
        continue
      fi

      php_test_methods "$php_file" | while IFS= read -r method; do
        [ -n "$method" ] || continue
        printf '%s::%s\n' "$base" "$method"
      done
    done
  } > "$output_path"
}

clone_source() {
  local repo="$1"
  local branch="$2"
  local cache_path="$3"

  if [ -d "$cache_path/.git" ]; then
    return 0
  fi

  git clone --depth 1 --branch "$branch" "https://github.com/$repo.git" "$cache_path" >/dev/null 2>&1
}

refresh() {
  need git
  need perl

  local tmp_path
  tmp_path="$(mktemp -d "${TMPDIR:-/tmp}/bedrock-laravel-compliance.XXXXXX")"
  trap 'rm -rf "$tmp_path"' EXIT

  local seen_sources="$tmp_path/sources.tsv"
  : > "$seen_sources"

  list_records | while IFS=$'\t' read -r id status repo branch tests_path inventory filter laravel bedrock; do
    if [ -z "$repo" ] || [ "$repo" = "null" ] || [ -z "$branch" ] || [ "$branch" = "null" ]; then
      continue
    fi

    if [ -z "$tests_path" ] || [ "$tests_path" = "null" ] || [ -z "$inventory" ] || [ "$inventory" = "null" ]; then
      continue
    fi

    local safe_name repo_path source_path output_path
    safe_name="$(printf '%s-%s' "$repo" "$branch" | tr '/:' '--')"
    repo_path="$tmp_path/$safe_name"
    source_path="$repo_path/$tests_path"
    output_path="$COMPLIANCE_PATH/$inventory"

    echo "Refreshing $id from $repo@$branch:$tests_path"
    clone_source "$repo" "$branch" "$repo_path"
    generate_inventory "$repo" "$branch" "$tests_path" "$filter" "$source_path" "$output_path"

    local sha
    sha="$(git -C "$repo_path" rev-parse HEAD)"
    printf '%s\t%s\t%s\t%s\n' "$id" "$repo" "$branch" "$sha" >> "$seen_sources"
  done

  write_sources_lock "$seen_sources"
}

write_sources_lock() {
  local source_file="$1"
  local tmp_file
  tmp_file="$(mktemp "${TMPDIR:-/tmp}/bedrock-sources-lock.XXXXXX")"

  {
    echo "{"
    echo "  \"version\": 1,"
    printf '  "generated_at": "%s",\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    echo "  \"sources\": ["

    awk -F '\t' '!seen[$2 "|" $3]++ { print $2 "\t" $3 "\t" $4 }' "$source_file" | awk -F '\t' '
      BEGIN { first = 1 }
      {
        if (!first) {
          print ","
        }

        first = 0
        printf "    {\n"
        printf "      \"repo\": \"%s\",\n", $1
        printf "      \"branch\": \"%s\",\n", $2
        printf "      \"commit\": \"%s\",\n", $3
        printf "      \"url\": \"https://github.com/%s/tree/%s\"\n", $1, $2
        printf "    }"
      }
      END {
        if (!first) {
          print ""
        }
      }
    '

    echo "  ]"
    echo "}"
  } > "$tmp_file"

  mv "$tmp_file" "$SOURCES_LOCK_FILE"
}

inventory_entries() {
  local file="$1"
  awk '
    /^[[:space:]]*$/ { next }
    /^#/ { next }
    /::/ { print }
  ' "$file"
}

entry_status() {
  local entry="$1"
  local php_file="${entry%%::*}"
  local method="${entry##*::}"
  local class="${php_file%.php}"

  if rg -Fq "$class::$method" "$ROOT_PATH/packages" "$ROOT_PATH/services" 2>/dev/null; then
    echo "ported"
    return 0
  fi

  if rg -Fq "$entry" "$DIVERGENCES_FILE" 2>/dev/null; then
    echo "adapted"
    return 0
  fi

  echo "missing"
}

inventory_stats() {
  local file="$1"
  local total=0
  local ported=0
  local adapted=0
  local missing=0

  while IFS= read -r entry; do
    [ -n "$entry" ] || continue
    total=$((total + 1))

    case "$(entry_status "$entry")" in
      ported) ported=$((ported + 1)) ;;
      adapted) adapted=$((adapted + 1)) ;;
      *) missing=$((missing + 1)) ;;
    esac
  done < <(inventory_entries "$file")

  printf '%s\t%s\t%s\t%s\n' "$total" "$ported" "$adapted" "$missing"
}

report() {
  need rg

  local tmp_file
  tmp_file="$(mktemp "${TMPDIR:-/tmp}/bedrock-compliance-report.XXXXXX")"

  {
    echo "# Laravel Compliance Report"
    echo
    echo "Generated: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
    echo
    echo "Source of truth: services/compliance"
    echo
    echo "## Inventories"
    echo
    echo "| Inventory | Total | Ported | Adapted | Missing |"
    echo "| --- | ---: | ---: | ---: | ---: |"

    find "$COMPLIANCE_PATH/inventories" -type f -name '*.txt' | sort | while IFS= read -r file; do
      local rel stats total ported adapted missing
      rel="${file#$COMPLIANCE_PATH/}"
      stats="$(inventory_stats "$file")"
      IFS=$'\t' read -r total ported adapted missing <<< "$stats"
      printf '| %s | %s | %s | %s | %s |\n' "$rel" "$total" "$ported" "$adapted" "$missing"
    done

    echo
    echo "## Permanent Exclusions"
    echo
    awk '
      /^  - id:/ {
        id = $0
        sub(/^  - id:[ \t]*/, "", id)
        next
      }
      id != "" && /^    laravel:/ {
        laravel = $0
        sub(/^    laravel:[ \t]*/, "", laravel)
        next
      }
      id != "" && /^    status: excluded-permanent/ {
        printf "- `%s`", laravel
        pending = 1
        next
      }
      pending && /^    rationale:/ {
        rationale = $0
        sub(/^    rationale:[ \t]*/, "", rationale)
        printf " - %s\n", rationale
        id = ""; laravel = ""; pending = 0
      }
    ' "$INVENTORY_FILE"
  } > "$tmp_file"

  mv "$tmp_file" "$REPORT_FILE"
}

check_no_package_local_parity() {
  local parity_files inventory_files
  parity_files="$(find "$ROOT_PATH/packages" -name PARITY.md -print)"
  inventory_files="$(find "$ROOT_PATH/packages" -path '*/testdata/*inventory*' -print)"

  if [ -n "$parity_files" ]; then
    echo "Package-local PARITY.md files are not allowed:" >&2
    echo "$parity_files" >&2
    return 1
  fi

  if [ -n "$inventory_files" ]; then
    echo "Package-local inventory files are not allowed:" >&2
    echo "$inventory_files" >&2
    return 1
  fi
}

check_excluded_not_implemented() {
  local failed=0

  while IFS=$'\t' read -r id status repo branch tests_path inventory filter laravel bedrock; do
    [ "$status" = "excluded-permanent" ] || continue

    case "$id" in
      framework.foundation) candidate="foundation" ;;
      framework.view) candidate="view" ;;
      framework.collections) candidate="collections" ;;
      framework.macroable) candidate="macroable" ;;
      framework.reflection) candidate="reflection" ;;
      package.browser_kit_testing) candidate="browser-kit-testing" ;;
      package.vite_plugin) candidate="vite-plugin" ;;
      package.starter_kits) candidate="starter-kits" ;;
      package.*) candidate="${id#package.}" ;;
      *) candidate="" ;;
    esac

    if [ -n "$candidate" ] && [ -d "$ROOT_PATH/packages/$candidate" ]; then
      echo "Permanent exclusion $id has an implemented package at packages/$candidate; update services/compliance/inventory.yml first." >&2
      failed=1
    fi
  done < <(list_records)

  return "$failed"
}

check_laravel_test_files_have_source_context() {
  local failed=0

  while IFS= read -r file; do
    if ! rg -q 'Port of|Ports of|Partial port of|Adapted from' "$file"; then
      echo "Laravel test file lacks source context: ${file#$ROOT_PATH/}" >&2
      failed=1
    fi
  done < <(find "$ROOT_PATH/packages" -type f -name '*_laravel_test.go' | sort)

  return "$failed"
}

check_no_inline_exclusions() {
  if rg -n '^[[:space:]]*// Excluded:' "$ROOT_PATH/packages" -g '*_test.go' >/tmp/bedrock-inline-exclusions.$$ 2>/dev/null; then
    echo "Inline test exclusions are not allowed; put exclusions in services/compliance:" >&2
    cat /tmp/bedrock-inline-exclusions.$$ >&2
    rm -f /tmp/bedrock-inline-exclusions.$$
    return 1
  fi

  rm -f /tmp/bedrock-inline-exclusions.$$
}

check() {
  need rg

  [ -f "$INVENTORY_FILE" ] || { echo "Missing $INVENTORY_FILE" >&2; return 1; }
  [ -f "$DIVERGENCES_FILE" ] || { echo "Missing $DIVERGENCES_FILE" >&2; return 1; }
  [ -f "$SOURCES_LOCK_FILE" ] || { echo "Missing $SOURCES_LOCK_FILE" >&2; return 1; }

  check_no_package_local_parity
  check_excluded_not_implemented
  check_laravel_test_files_have_source_context
  check_no_inline_exclusions

  find "$COMPLIANCE_PATH/inventories" -type f -name '*.txt' | sort | while IFS= read -r file; do
    if ! inventory_entries "$file" >/dev/null; then
      echo "Invalid inventory file: ${file#$ROOT_PATH/}" >&2
      return 1
    fi
  done
}

case "${1:-}" in
  refresh)
    refresh
    report
    ;;
  check)
    check
    ;;
  report)
    report
    ;;
  -h|--help|help|"")
    usage
    ;;
  *)
    echo "Unknown command: $1" >&2
    usage >&2
    exit 1
    ;;
esac
