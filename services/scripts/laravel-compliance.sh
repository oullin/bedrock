#!/usr/bin/env bash
set -euo pipefail

ROOT_PATH="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
COMPLIANCE_PATH="$ROOT_PATH/services/compliance"
INVENTORY_FILE="$COMPLIANCE_PATH/inventory.yml"
DIVERGENCES_FILE="$COMPLIANCE_PATH/divergences.yml"
FEATURES_FILE="$COMPLIANCE_PATH/features.yml"
SOURCES_LOCK_FILE="$COMPLIANCE_PATH/sources.lock.json"
REPORT_FILE="$COMPLIANCE_PATH/report.md"
RECORD_SEPARATOR=$'\034'

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
        print id "\034" status "\034" repo "\034" branch "\034" tests_path "\034" inventory "\034" filter "\034" laravel "\034" bedrock
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

list_feature_records() {
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
      id = ""; source_id = ""; name = ""; laravel = "";
      bedrock = ""; status = ""; docs = ""; notes = "";
    }
    function emit() {
      if (id != "") {
        print id "\034" source_id "\034" name "\034" laravel "\034" bedrock "\034" status "\034" docs "\034" notes
      }
    }
    /^  - id:/ {
      emit()
      reset_record()
      sub(/^  - id:[ \t]*/, "")
      id = clean($0)
      next
    }
    id != "" && /^    source_id:/ {
      sub(/^    source_id:[ \t]*/, "")
      source_id = clean($0)
      next
    }
    id != "" && /^    name:/ {
      sub(/^    name:[ \t]*/, "")
      name = clean($0)
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
    id != "" && /^    status:/ {
      sub(/^    status:[ \t]*/, "")
      status = clean($0)
      next
    }
    id != "" && /^    docs:/ {
      sub(/^    docs:[ \t]*/, "")
      docs = clean($0)
      next
    }
    id != "" && /^    notes:/ {
      sub(/^    notes:[ \t]*/, "")
      notes = clean($0)
      next
    }
    END {
      emit()
    }
  ' "$FEATURES_FILE"
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
    echo "# Laravel compliance test inventory."
    if [[ "$repo" == local:* ]]; then
      echo "# Source: ${repo#local:}/$tests_path"
    else
      echo "# Source: https://github.com/$repo/tree/$branch/$tests_path"
    fi
    echo "# Format: <file>::<test method name>"
    echo "# Generated: $(date -u +%Y-%m-%d)"
    echo "#"
    echo "# Every entry below must have a Go counterpart, an adaptation entry, or an"
    echo "# explicit missing status in services/compliance/report.md."
    echo

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
      echo "laravel-compliance: missing local source path: $local_path" >&2
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

  REFRESH_TMP_PATH="$(mktemp -d "${TMPDIR:-/tmp}/bedrock-laravel-compliance.XXXXXX")"
  trap 'rm -rf "$REFRESH_TMP_PATH"' EXIT

  local tmp_path
  tmp_path="$REFRESH_TMP_PATH"

  local seen_sources="$tmp_path/sources.tsv"
  : > "$seen_sources"

  while IFS="$RECORD_SEPARATOR" read -r id status repo branch tests_path inventory filter laravel bedrock; do
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
  done < <(list_records)

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

build_status_index() {
  local ported_index="$1"
  local adapted_index="$2"

  {
    rg -o --no-filename '[A-Za-z_][A-Za-z0-9_]*Test::[A-Za-z_][A-Za-z0-9_]*' "$ROOT_PATH/packages" -g '*_test.go' 2>/dev/null || true
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

markdown_cell() {
  local value="$1"
  value="${value//|/\\|}"
  printf '%s' "$value"
}

docs_file_for() {
  local bedrock="$1"
  local docs="$2"

  if [ -n "$docs" ] && [ "$docs" != "null" ]; then
    case "$docs" in
      services/docs/*) printf '%s\n' "$ROOT_PATH/$docs" ;;
      packages/*|basics/*|concepts/*) printf '%s\n' "$ROOT_PATH/services/docs/$docs" ;;
      *) printf '%s\n' "$ROOT_PATH/services/docs/packages/$docs" ;;
    esac
    return 0
  fi

  [ -n "$bedrock" ] && [ "$bedrock" != "null" ] || return 0
  printf '%s\n' "$ROOT_PATH/services/docs/packages/$bedrock.md"
}

docs_display() {
  local bedrock="$1"
  local docs="${2:-}"

  if { [ -z "$bedrock" ] || [ "$bedrock" = "null" ]; } && { [ -z "$docs" ] || [ "$docs" = "null" ]; }; then
    printf 'n/a'
    return 0
  fi

  local docs_file
  docs_file="$(docs_file_for "$bedrock" "$docs")"

  if [ -n "$docs_file" ] && [ -f "$docs_file" ]; then
    printf '`%s`' "${docs_file#$ROOT_PATH/services/docs/}"
  else
    printf 'missing'
  fi
}

tests_display() {
  local inventory="$1"
  local ported_index="$2"
  local adapted_index="$3"

  if [ -z "$inventory" ] || [ "$inventory" = "null" ]; then
    printf 'no inventory'
    return 0
  fi

  local file stats total ported adapted missing
  file="$COMPLIANCE_PATH/$inventory"
  stats="$(inventory_stats "$file" "$ported_index" "$adapted_index")"
  IFS=$'\t' read -r total ported adapted missing <<< "$stats"

  printf '%s/%s ported, %s adapted, %s missing' "$ported" "$total" "$adapted" "$missing"
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

    list_records | while IFS="$RECORD_SEPARATOR" read -r id status repo branch tests_path inventory filter laravel bedrock; do
      [ -n "$inventory" ] && [ "$inventory" != "null" ] || continue

      local file
      file="$COMPLIANCE_PATH/$inventory"

      local rel stats total ported adapted missing
      rel="${file#$COMPLIANCE_PATH/}"
      stats="$(inventory_stats "$file" "$ported_index" "$adapted_index")"
      IFS=$'\t' read -r total ported adapted missing <<< "$stats"
      printf '| %s | %s | %s | %s | %s |\n' "$rel" "$total" "$ported" "$adapted" "$missing"
    done

    echo
    echo "## Framework Coverage"
    echo
    echo "| Laravel Source | Bedrock Surface | Port Status | Tests | Docs |"
    echo "| --- | --- | --- | --- | --- |"

    list_records | while IFS="$RECORD_SEPARATOR" read -r id status repo branch tests_path inventory filter laravel bedrock; do
      case "$id" in
        framework.*) ;;
        *) continue ;;
      esac

      local display_bedrock display_tests display_docs
      if [ -n "$bedrock" ] && [ "$bedrock" != "null" ]; then
        display_bedrock="\`$bedrock\`"
      else
        display_bedrock="n/a"
      fi

      display_tests="$(tests_display "$inventory" "$ported_index" "$adapted_index")"
      display_docs="$(docs_display "$bedrock" "")"

      printf '| `%s` | %s | `%s` | %s | %s |\n' \
        "$(markdown_cell "$laravel")" \
        "$display_bedrock" \
        "$(markdown_cell "$status")" \
        "$(markdown_cell "$display_tests")" \
        "$display_docs"
    done

    echo
    echo "## Feature Coverage"
    echo
    echo "| Feature | Source | Bedrock Surface | Status | Docs | Notes |"
    echo "| --- | --- | --- | --- | --- | --- |"

    list_feature_records | while IFS="$RECORD_SEPARATOR" read -r id source_id name laravel bedrock status docs notes; do
      local display_bedrock display_docs display_notes
      if [ -n "$bedrock" ] && [ "$bedrock" != "null" ]; then
        display_bedrock="\`$bedrock\`"
      else
        display_bedrock="n/a"
      fi

      display_docs="$(docs_display "$bedrock" "$docs")"
      display_notes="${notes:-}"
      [ -n "$display_notes" ] && [ "$display_notes" != "null" ] || display_notes="n/a"

      printf '| %s | `%s` | %s | `%s` | %s | %s |\n' \
        "$(markdown_cell "$name")" \
        "$(markdown_cell "$source_id")" \
        "$display_bedrock" \
        "$(markdown_cell "$status")" \
        "$display_docs" \
        "$(markdown_cell "$display_notes")"
    done

    echo
    echo "## Mapped Sources Without Inventories"
    echo
    echo "| Source | Bedrock | Reason |"
    echo "| --- | --- | --- |"

    list_records | while IFS="$RECORD_SEPARATOR" read -r id status repo branch tests_path inventory filter laravel bedrock; do
      [ "$status" = "mapped" ] || continue
      [ -z "$inventory" ] || [ "$inventory" = "null" ] || continue

      printf '| `%s` | `%s` | No generated upstream test inventory configured. |\n' "$laravel" "$bedrock"
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

  rm -rf "$status_index_path"
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

  while IFS="$RECORD_SEPARATOR" read -r id status repo branch tests_path inventory filter laravel bedrock; do
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
  local inline_exclusions

  if inline_exclusions="$(rg -n '^[[:space:]]*// Excluded:' "$ROOT_PATH/packages" -g '*_test.go' 2>/dev/null)"; then
    echo "Inline test exclusions are not allowed; put exclusions in services/compliance:" >&2
    printf '%s\n' "$inline_exclusions" >&2
    return 1
  fi
}

check_inventory_files_configured() {
  local expected actual unexpected missing failed=0
  expected="$(mktemp "${TMPDIR:-/tmp}/bedrock-compliance-expected.XXXXXX")"
  actual="$(mktemp "${TMPDIR:-/tmp}/bedrock-compliance-actual.XXXXXX")"

  {
    while IFS="$RECORD_SEPARATOR" read -r id status repo branch tests_path inventory filter laravel bedrock; do
      [ -n "$inventory" ] && [ "$inventory" != "null" ] || continue
      printf '%s\n' "$COMPLIANCE_PATH/$inventory"
    done < <(list_records)
  } | sort -u > "$expected"

  find "$COMPLIANCE_PATH/inventories" -type f -name '*.txt' | sort -u > "$actual"

  unexpected="$(comm -13 "$expected" "$actual")"
  missing="$(comm -23 "$expected" "$actual")"

  if [ -n "$unexpected" ]; then
    echo "Unconfigured compliance inventory files:" >&2
    sed "s#^$ROOT_PATH/##" <<< "$unexpected" >&2
    failed=1
  fi

  if [ -n "$missing" ]; then
    echo "Missing configured compliance inventory files:" >&2
    sed "s#^$ROOT_PATH/##" <<< "$missing" >&2
    failed=1
  fi

  rm -f "$expected" "$actual"
  return "$failed"
}

check_inventory_format() {
  local failed=0

  while IFS="$RECORD_SEPARATOR" read -r id status repo branch tests_path inventory filter laravel bedrock; do
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

check_features_format() {
  awk '
    /^[[:space:]]*$/ || /^#/ { next }
    /^version:[ \t]*[0-9]+[ \t]*$/ { next }
    /^features:[ \t]*$/ { next }
    /^  - id:[ \t]*[^[:space:]].*$/ { next }
    /^    (source_id|name|laravel|bedrock|status|docs|notes):[ \t]*.*$/ { next }
    { invalid = 1; print FILENAME ":" FNR ": invalid features entry: " $0 > "/dev/stderr" }
    END { exit invalid ? 1 : 0 }
  ' "$FEATURES_FILE"
}

check_features() {
  local failed=0 count=0 source_ids feature_ids
  source_ids="$(mktemp "${TMPDIR:-/tmp}/bedrock-compliance-source-ids.XXXXXX")"
  feature_ids="$(mktemp "${TMPDIR:-/tmp}/bedrock-compliance-feature-ids.XXXXXX")"

  list_records | awk -F "$RECORD_SEPARATOR" '{ print $1 }' | sort -u > "$source_ids"
  : > "$feature_ids"

  check_features_format || failed=1

  while IFS="$RECORD_SEPARATOR" read -r id source_id name laravel bedrock status docs notes; do
    count=$((count + 1))

    if [ -z "$id" ] || [ -z "$source_id" ] || [ -z "$name" ] || [ -z "$laravel" ] || [ -z "$bedrock" ] || [ -z "$status" ]; then
      echo "Feature $id is missing a required field." >&2
      failed=1
    fi

    case "$status" in
      ported|partial|missing|excluded) ;;
      *)
        echo "Feature $id has invalid status: $status" >&2
        failed=1
        ;;
    esac

    if ! rg -Fxq -- "$source_id" "$source_ids"; then
      echo "Feature $id references unknown source_id: $source_id" >&2
      failed=1
    fi

    if rg -Fxq -- "$id" "$feature_ids"; then
      echo "Duplicate feature id: $id" >&2
      failed=1
    fi
    printf '%s\n' "$id" >> "$feature_ids"

    if [ -n "$docs" ] && [ "$docs" != "null" ]; then
      case "$docs" in
        services/docs/*|packages/*|basics/*|concepts/*) ;;
        *)
          echo "Feature $id docs path must be under services/docs: $docs" >&2
          failed=1
          ;;
      esac

      if [ "$status" = "ported" ] || [ "$status" = "partial" ]; then
        local docs_file
        docs_file="$(docs_file_for "$bedrock" "$docs")"

        if [ ! -f "$docs_file" ]; then
          echo "Feature $id docs path does not exist: ${docs_file#$ROOT_PATH/}" >&2
          failed=1
        fi
      fi
    fi
  done < <(list_feature_records)

  if [ "$count" -eq 0 ]; then
    echo "$FEATURES_FILE must define at least one feature." >&2
    failed=1
  fi

  rm -f "$source_ids" "$feature_ids"
  return "$failed"
}

check() {
  need rg

  [ -f "$INVENTORY_FILE" ] || { echo "Missing $INVENTORY_FILE" >&2; return 1; }
  [ -f "$DIVERGENCES_FILE" ] || { echo "Missing $DIVERGENCES_FILE" >&2; return 1; }
  [ -f "$FEATURES_FILE" ] || { echo "Missing $FEATURES_FILE" >&2; return 1; }
  [ -f "$SOURCES_LOCK_FILE" ] || { echo "Missing $SOURCES_LOCK_FILE" >&2; return 1; }

  check_no_package_local_parity
  check_excluded_not_implemented
  check_laravel_test_files_have_source_context
  check_no_inline_exclusions
  check_inventory_files_configured
  check_inventory_format
  check_features
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
