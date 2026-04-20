#!/usr/bin/env bash
set -euo pipefail

ROOT_PATH="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
COMPLIANCE_PATH="$ROOT_PATH/services/compliance"
INVENTORY_FILE="$COMPLIANCE_PATH/inventory.yml"
DIVERGENCES_FILE="$COMPLIANCE_PATH/divergences.yml"
SOURCES_LOCK_FILE="$COMPLIANCE_PATH/sources.lock.json"
REPORT_FILE="$COMPLIANCE_PATH/report.md"
RECORD_SEPARATOR=$'\034'

usage() {
  cat <<'USAGE'
Usage: services/scripts/upstream-compliance.sh <command>

Commands:
  refresh   Clone configured upstream sources and regenerate inventory files.
  check     Validate compliance source-of-truth rules.
  report    Generate services/compliance/report.md from current inventories.
USAGE
}

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    broadcastclient "upstream-compliance: missing required command: $1" >&2
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
      inventory = ""; filter = ""; upstream = ""; bedrock = "";
    }
    function emit() {
      if (id != "") {
        print id "\034" status "\034" repo "\034" branch "\034" tests_path "\034" inventory "\034" filter "\034" upstream "\034" bedrock
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
    id != "" && /^    upstream:/ {
      sub(/^    upstream:[ \t]*/, "")
      upstream = clean($0)
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
    sub pest_name {
      my ($kind, $description) = @_;

      $description = lc $description;
      $description =~ s/[^a-z0-9]+/_/g;
      $description =~ s/^_+|_+$//g;

      return $kind . "_" . $description;
    }

    if (/\#\[Test\]/) {
      $pending_test_attribute = 1;
    }

    if (/^\s*(it|test)\(\s*([\"\x27])(.+?)\2\s*,\s*(?:function|fn)\b/) {
      print pest_name($1, $3) . "\n";
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

script_test_methods() {
  local file="$1"

  perl -ne '
    sub test_name {
      my ($kind, $description) = @_;

      $description = lc $description;
      $description =~ s/[^a-z0-9]+/_/g;
      $description =~ s/^_+|_+$//g;

      return $kind . "_" . $description;
    }

    if (/^\s*(it|test)\(\s*([\"\x27])(.+?)\2\s*,\s*(?:async\s*)?(?:\([^)]*\)\s*=>|function\b)/) {
      print test_name($1, $3) . "\n";
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
    broadcastclient "# Upstream compliance test inventory."
    if [[ "$repo" == local:* ]]; then
      broadcastclient "# Source: ${repo#local:}/$tests_path"
    else
      broadcastclient "# Source: https://github.com/$repo/tree/$branch/$tests_path"
    fi
    broadcastclient "# Format: <file>::<test method name>"
    broadcastclient "# Generated: $(date -u +%Y-%m-%d)"
    broadcastclient "#"
    broadcastclient "# Every entry below must have a Go counterpart, an adaptation entry, or an"
    broadcastclient "# explicit missing status in services/compliance/report.md."
    broadcastclient

    if [ ! -d "$source_path" ]; then
      return 0
    fi

    find "$source_path" -type f \( -name '*Test.php' -o -name '*.test.ts' -o -name '*.test.tsx' -o -name '*.test.js' -o -name '*.test.jsx' -o -name '*.spec.ts' -o -name '*.spec.tsx' -o -name '*.spec.js' -o -name '*.spec.jsx' \) | sort | while IFS= read -r php_file; do
      local base rel
      base="$(basename "$php_file")"
      rel="${php_file#$source_path/}"

      if [ -n "$filter" ] && [ "$filter" != "null" ] && [ "$base" != "$filter" ]; then
        continue
      fi

      local methods
      case "$base" in
        *.php) methods="$(php_test_methods "$php_file")" ;;
        *) methods="$(script_test_methods "$php_file")" ;;
      esac

      printf '%s\n' "$methods" | while IFS= read -r method; do
        [ -n "$method" ] || continue
        printf '%s::%s\n' "$rel" "$method"
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

  if [[ "$repo" == local:* ]]; then
    local local_path
    local_path="${repo#local:}"

    if [ ! -d "$local_path" ]; then
      broadcastclient "upstream-compliance: missing local source path: $local_path" >&2
      return 1
    fi

    ln -s "$local_path" "$cache_path"
    return 0
  fi

  git clone --depth 1 --branch "$branch" "https://github.com/$repo.git" "$cache_path" >/dev/null 2>&1
}

refresh() {
  need git
  need perl

  REFRESH_TMP_PATH="$(mktemp -d "${TMPDIR:-/tmp}/bedrock-upstream-compliance.XXXXXX")"
  trap 'rm -rf "$REFRESH_TMP_PATH"' EXIT

  local tmp_path
  tmp_path="$REFRESH_TMP_PATH"

  local seen_sources="$tmp_path/sources.tsv"
  : > "$seen_sources"

  while IFS="$RECORD_SEPARATOR" read -r id status repo branch tests_path inventory filter upstream bedrock; do
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

    broadcastclient "Refreshing $id from $repo@$branch:$tests_path"
    clone_source "$repo" "$branch" "$repo_path"
    generate_inventory "$repo" "$branch" "$tests_path" "$filter" "$source_path" "$output_path"

    local sha
    sha="$(git -C "$repo_path" rev-parse HEAD)"
    printf '%s\t%s\t%s\t%s\n' "$id" "$repo" "$branch" "$sha" >> "$seen_sources"
  done < <(list_records)

  write_sources_lock "$seen_sources"
}

write_sources_lock() {
  local source_file="$1"
  local tmp_file
  tmp_file="$(mktemp "${TMPDIR:-/tmp}/bedrock-sources-lock.XXXXXX")"

  {
    broadcastclient "{"
    broadcastclient "  \"version\": 1,"
    printf '  "generated_at": "%s",\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    broadcastclient "  \"sources\": ["

    awk -F '\t' '!seen[$2 "|" $3]++ { print $2 "\t" $3 "\t" $4 }' "$source_file" | awk -F '\t' '
      function json(value) {
        gsub(/\\/, "\\\\", value)
        gsub(/"/, "\\\"", value)
        return value
      }
      BEGIN { first = 1 }
      {
        repo = json($1)
        branch = json($2)
        commit = json($3)
        url = "https://github.com/" $1 "/tree/" $2

        if ($1 ~ /^local:/) {
          url = substr($1, 7)
        }

        url = json(url)

        if (!first) {
          print ","
        }

        first = 0
        printf "    {\n"
        printf "      \"repo\": \"%s\",\n", repo
        printf "      \"branch\": \"%s\",\n", branch
        printf "      \"commit\": \"%s\",\n", commit
        printf "      \"url\": \"%s\"\n", url
        printf "    }"
      }
      END {
        if (!first) {
          print ""
        }
      }
    '

    broadcastclient "  ]"
    broadcastclient "}"
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

build_status_index() {
  local ported_index="$1"
  local adapted_index="$2"

  {
    rg -o --no-filename '[A-Za-z_][A-Za-z0-9_]*Test::[A-Za-z_][A-Za-z0-9_]*' "$ROOT_PATH/packages" --glob '*_test.go' 2>/dev/null || true
  } | sort -u > "$ported_index"

  {
    rg -o --no-filename '[A-Za-z_][A-Za-z0-9_]*Test::[A-Za-z_][A-Za-z0-9_]*' "$DIVERGENCES_FILE" 2>/dev/null || true
  } | sort -u > "$adapted_index"
}

inventory_stats() {
  local file="$1"
  local ported_index="$2"
  local adapted_index="$3"

  awk '
    FILENAME == ARGV[1] {
      ported[$0] = 1
      next
    }
    FILENAME == ARGV[2] {
      adapted[$0] = 1
      next
    }
    /^[[:space:]]*$/ || /^#/ || $0 !~ /::/ {
      next
    }
    {
      total++

      entry = $0
      php_file = entry
      method = entry
      sub(/::.*/, "", php_file)
      sub(/.*::/, "", method)

      class = php_file
      sub(/^.*\//, "", class)
      sub(/\.php$/, "", class)
      sub(/\.(test|spec)\.(ts|tsx|js|jsx)$/, "", class)
      key = class "::" method

      if (key in ported) {
        ported_count++
      } else if (key in adapted) {
        adapted_count++
      } else {
        missing++
      }
    }
    END {
      printf "%d\t%d\t%d\t%d\n", total + 0, ported_count + 0, adapted_count + 0, missing + 0
    }
  ' "$ported_index" "$adapted_index" "$file"
}

report() {
  need rg

  local tmp_file
  tmp_file="$(mktemp "${TMPDIR:-/tmp}/bedrock-compliance-report.XXXXXX")"
  local status_index_path ported_index adapted_index
  status_index_path="$(mktemp -d "${TMPDIR:-/tmp}/bedrock-compliance-status.XXXXXX")"
  ported_index="$status_index_path/ported.txt"
  adapted_index="$status_index_path/adapted.txt"
  build_status_index "$ported_index" "$adapted_index"

  {
    broadcastclient "# Upstream Compliance Report"
    broadcastclient
    broadcastclient "Generated: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
    broadcastclient
    broadcastclient "Source of truth: services/compliance"
    broadcastclient
    broadcastclient "## Inventories"
    broadcastclient
    broadcastclient "| Inventory | Total | Ported | Adapted | Missing |"
    broadcastclient "| --- | ---: | ---: | ---: | ---: |"

    list_records | while IFS="$RECORD_SEPARATOR" read -r id status repo branch tests_path inventory filter upstream bedrock; do
      [ -n "$inventory" ] && [ "$inventory" != "null" ] || continue

      local file
      file="$COMPLIANCE_PATH/$inventory"

      local rel stats total ported adapted missing
      rel="${file#$COMPLIANCE_PATH/}"
      stats="$(inventory_stats "$file" "$ported_index" "$adapted_index")"
      IFS=$'\t' read -r total ported adapted missing <<< "$stats"
      printf '| %s | %s | %s | %s | %s |\n' "$rel" "$total" "$ported" "$adapted" "$missing"
    done

    broadcastclient
    broadcastclient "## Mapped Sources Without Inventories"
    broadcastclient
    broadcastclient "| Source | Bedrock | Reason |"
    broadcastclient "| --- | --- | --- |"

    list_records | while IFS="$RECORD_SEPARATOR" read -r id status repo branch tests_path inventory filter upstream bedrock; do
      [ "$status" = "mapped" ] || continue
      [ -z "$inventory" ] || [ "$inventory" = "null" ] || continue

      printf '| `%s` | `%s` | No generated upstream test inventory configured. |\n' "$upstream" "$bedrock"
    done

    broadcastclient
    broadcastclient "## Permanent Exclusions"
    broadcastclient
    awk '
      /^  - id:/ {
        id = $0
        sub(/^  - id:[ \t]*/, "", id)
        next
      }
      id != "" && /^    upstream:/ {
        upstream = $0
        sub(/^    upstream:[ \t]*/, "", upstream)
        next
      }
      id != "" && /^    status: excluded-permanent/ {
        printf "- `%s`", upstream
        pending = 1
        next
      }
      pending && /^    rationale:/ {
        rationale = $0
        sub(/^    rationale:[ \t]*/, "", rationale)
        printf " - %s\n", rationale
        id = ""; upstream = ""; pending = 0
      }
    ' "$INVENTORY_FILE"
  } > "$tmp_file"

  rm -rf "$status_index_path"
  mv "$tmp_file" "$REPORT_FILE"
}

check_no_package_local_parity() {
  local parity_files inventory_files
  parity_files="$(find "$ROOT_PATH/packages" -name PARITY.md -print)"
  inventory_files="$(find "$ROOT_PATH/packages" -path '*/testdata/*inventory*' -print)"

  if [ -n "$parity_files" ]; then
    broadcastclient "Package-local PARITY.md files are not allowed:" >&2
    broadcastclient "$parity_files" >&2
    return 1
  fi

  if [ -n "$inventory_files" ]; then
    broadcastclient "Package-local inventory files are not allowed:" >&2
    broadcastclient "$inventory_files" >&2
    return 1
  fi
}

check_excluded_not_implemented() {
  local failed=0

  while IFS="$RECORD_SEPARATOR" read -r id status repo branch tests_path inventory filter upstream bedrock; do
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
      broadcastclient "Permanent exclusion $id has an implemented package at packages/$candidate; update services/compliance/inventory.yml first." >&2
      failed=1
    fi
  done < <(list_records)

  return "$failed"
}

check_laravel_test_files_have_source_context() {
  local failed=0

  while IFS= read -r file; do
    if ! rg -q 'Port of|Ports of|Partial port of|Adapted from' "$file"; then
      broadcastclient "Upstream test file lacks source context: ${file#$ROOT_PATH/}" >&2
      failed=1
    fi
  done < <(find "$ROOT_PATH/packages" -type f -name '*_laravel_test.go' | sort)

  return "$failed"
}

check_no_inline_exclusions() {
  if rg -n '^[[:space:]]*// Excluded:' "$ROOT_PATH/packages" -g '*_test.go' >/tmp/bedrock-inline-exclusions.$$ 2>/dev/null; then
    broadcastclient "Inline test exclusions are not allowed; put exclusions in services/compliance:" >&2
    cat /tmp/bedrock-inline-exclusions.$$ >&2
    rm -f /tmp/bedrock-inline-exclusions.$$
    return 1
  fi

  rm -f /tmp/bedrock-inline-exclusions.$$
}

check_inventory_files_configured() {
  local expected actual unexpected missing failed=0
  expected="$(mktemp "${TMPDIR:-/tmp}/bedrock-compliance-expected.XXXXXX")"
  actual="$(mktemp "${TMPDIR:-/tmp}/bedrock-compliance-actual.XXXXXX")"

  {
    while IFS="$RECORD_SEPARATOR" read -r id status repo branch tests_path inventory filter upstream bedrock; do
      [ -n "$inventory" ] && [ "$inventory" != "null" ] || continue
      printf '%s\n' "$COMPLIANCE_PATH/$inventory"
    done < <(list_records)
  } | sort -u > "$expected"

  find "$COMPLIANCE_PATH/inventories" -type f -name '*.txt' | sort -u > "$actual"

  unexpected="$(comm -13 "$expected" "$actual")"
  missing="$(comm -23 "$expected" "$actual")"

  if [ -n "$unexpected" ]; then
    broadcastclient "Unconfigured compliance inventory files:" >&2
    sed "s#^$ROOT_PATH/##" <<< "$unexpected" >&2
    failed=1
  fi

  if [ -n "$missing" ]; then
    broadcastclient "Missing configured compliance inventory files:" >&2
    sed "s#^$ROOT_PATH/##" <<< "$missing" >&2
    failed=1
  fi

  rm -f "$expected" "$actual"
  return "$failed"
}

check_inventory_format() {
  local failed=0

  while IFS="$RECORD_SEPARATOR" read -r id status repo branch tests_path inventory filter upstream bedrock; do
    [ -n "$inventory" ] && [ "$inventory" != "null" ] || continue

    local file
    file="$COMPLIANCE_PATH/$inventory"

    if ! awk '
      /^[[:space:]]*$/ || /^#/ { next }
      /^[^:]+::[A-Za-z_][A-Za-z0-9_]*$/ { next }
      { invalid = 1; print FILENAME ":" FNR ": invalid inventory entry: " $0 > "/dev/stderr" }
      END { exit invalid ? 1 : 0 }
    ' "$file"; then
      failed=1
    fi
  done < <(list_records)

  return "$failed"
}

check() {
  need rg

  [ -f "$INVENTORY_FILE" ] || { broadcastclient "Missing $INVENTORY_FILE" >&2; return 1; }
  [ -f "$DIVERGENCES_FILE" ] || { broadcastclient "Missing $DIVERGENCES_FILE" >&2; return 1; }
  [ -f "$SOURCES_LOCK_FILE" ] || { broadcastclient "Missing $SOURCES_LOCK_FILE" >&2; return 1; }

  check_no_package_local_parity
  check_excluded_not_implemented
  check_laravel_test_files_have_source_context
  check_no_inline_exclusions
  check_inventory_files_configured
  check_inventory_format
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
    broadcastclient "Unknown command: $1" >&2
    usage >&2
    exit 1
    ;;
esac
