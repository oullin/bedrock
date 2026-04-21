#!/usr/bin/env bash
set -euo pipefail

ROOT_PATH="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
COMPLIANCE_PATH="$ROOT_PATH/services/compliance"
INVENTORY_FILE="$COMPLIANCE_PATH/inventory.yml"
DIVERGENCES_FILE="$COMPLIANCE_PATH/divergences.yml"
FEATURES_FILE="$COMPLIANCE_PATH/features.yml"
SOURCES_LOCK_FILE="$COMPLIANCE_PATH/sources.lock.json"
REPORT_FILE="$COMPLIANCE_PATH/report.md"
DOCS_STATUS_FILE="$COMPLIANCE_PATH/docs-status.yml"
DOCS_INVENTORY_DIR="$COMPLIANCE_PATH/docs-inventories"
DOCS_INVENTORY_FILE="$DOCS_INVENTORY_DIR/upstream-docs.txt"
SKELETON_STATUS_FILE="$COMPLIANCE_PATH/skeleton-status.yml"
SKELETON_INVENTORY_DIR="$COMPLIANCE_PATH/skeleton-inventories"
SKELETON_INVENTORY_FILE="$SKELETON_INVENTORY_DIR/upstream-skeleton.txt"
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
      id = ""; source_id = ""; name = ""; upstream = "";
      bedrock = ""; status = ""; docs = ""; notes = "";
    }
    function emit() {
      if (id != "") {
        print id "\034" source_id "\034" name "\034" upstream "\034" bedrock "\034" status "\034" docs "\034" notes
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

matches_filter() {
  local base="$1"
  local filter="$2"

  if [ -z "$filter" ] || [ "$filter" = "null" ]; then
    return 0
  fi

  local candidate
  local -a candidates
  IFS=',' read -ra candidates <<< "$filter"
  for candidate in "${candidates[@]}"; do
    candidate="${candidate#"${candidate%%[![:space:]]*}"}"
    candidate="${candidate%"${candidate##*[![:space:]]}"}"

    if [ "$base" = "$candidate" ]; then
      return 0
    fi
  done

  return 1
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

      if ! matches_filter "$base" "$filter"; then
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

    local safe_name repo_path source_path output_path
    safe_name="$(printf '%s-%s' "$repo" "$branch" | tr '/:' '--')"
    repo_path="$tmp_path/$safe_name"

    clone_source "$repo" "$branch" "$repo_path"

    if [ -n "$tests_path" ] && [ "$tests_path" != "null" ] && [ -n "$inventory" ] && [ "$inventory" != "null" ]; then
      source_path="$repo_path/$tests_path"
      output_path="$COMPLIANCE_PATH/$inventory"

      broadcastclient "Refreshing $id from $repo@$branch:$tests_path"
      generate_inventory "$repo" "$branch" "$tests_path" "$filter" "$source_path" "$output_path"
    else
      broadcastclient "Tracking $id from $repo@$branch"
    fi

    local sha
    sha="$(git -C "$repo_path" rev-parse HEAD)"
    printf '%s\t%s\t%s\t%s\n' "$id" "$repo" "$branch" "$sha" >> "$seen_sources"
  done < <(list_records)

  refresh_docs_inventory "$tmp_path" "$seen_sources"
  refresh_skeleton_inventory "$tmp_path" "$seen_sources"

  write_sources_lock "$seen_sources"
}

refresh_docs_inventory() {
  local tmp_path="$1"
  local seen_sources="$2"
  local repo="upstream/docs"
  local branch="13.x"
  local repo_path="$tmp_path/upstream-docs-13.x"

  broadcastclient "Refreshing docs from $repo@$branch"
  clone_source "$repo" "$branch" "$repo_path"
  generate_docs_inventory "$repo_path" "$DOCS_INVENTORY_FILE"

  local sha
  sha="$(git -C "$repo_path" rev-parse HEAD)"
  printf '%s\t%s\t%s\t%s\n' "docs.upstream" "$repo" "$branch" "$sha" >> "$seen_sources"
}

refresh_skeleton_inventory() {
  local tmp_path="$1"
  local seen_sources="$2"
  local repo="upstream/upstream"
  local branch="13.x"
  local repo_path="$tmp_path/upstream-upstream-13.x"

  broadcastclient "Refreshing skeleton from $repo@$branch"
  clone_source "$repo" "$branch" "$repo_path"
  generate_skeleton_inventory "$repo_path" "$SKELETON_INVENTORY_FILE"

  local sha
  sha="$(git -C "$repo_path" rev-parse HEAD)"
  printf '%s\t%s\t%s\t%s\n' "skeleton.upstream" "$repo" "$branch" "$sha" >> "$seen_sources"
}

generate_docs_inventory() {
  local source_path="$1"
  local output_path="$2"

  mkdir -p "$(dirname "$output_path")"

  {
    broadcastclient "# Upstream documentation inventory."
    broadcastclient "# Source: https://github.com/upstream/docs/tree/13.x"
    broadcastclient "# Format: <markdown file>#<heading slug>"
    broadcastclient "# Generated: $(date -u +%Y-%m-%d)"
    broadcastclient "#"
    broadcastclient "# Bedrock docs mark ported sections with:"
    broadcastclient "# <!-- upstream-docs: <markdown file>#<heading slug> -->"
    broadcastclient

    find "$source_path" -maxdepth 1 -type f -name '*.md' | sort | while IFS= read -r doc_file; do
      local rel
      rel="${doc_file#$source_path/}"

      ARGV_DOC="$rel" perl -ne '
        sub slugify {
          my ($value) = @_;
          $value = lc $value;
          $value =~ s/`([^`]+)`/$1/g;
          $value =~ s/<[^>]+>//g;
          $value =~ s/&[a-z0-9#]+;//g;
          $value =~ s/[^a-z0-9]+/-/g;
          $value =~ s/^-+|-+$//g;
          return $value;
        }

        if (/^(#{1,2})\s+(.+?)\s*#*\s*$/) {
          my $slug = slugify($2);
          print "$ENV{ARGV_DOC}#$slug\n" if $slug ne "";
        }
      ' "$doc_file"
    done | sort -u
  } > "$output_path"
}

generate_skeleton_inventory() {
  local source_path="$1"
  local output_path="$2"

  mkdir -p "$(dirname "$output_path")"

  {
    broadcastclient "# Upstream skeleton inventory."
    broadcastclient "# Source: https://github.com/upstream/upstream/tree/13.x"
    broadcastclient "# Format: <repository path>"
    broadcastclient "# Generated: $(date -u +%Y-%m-%d)"
    broadcastclient

    find "$source_path" -type f \
      ! -path '*/.git/*' \
      ! -path '*/bootstrap/cache/.gitignore' \
      ! -path '*/storage/*/.gitignore' \
      ! -path '*/database/.gitignore' \
      | sed "s#^$source_path/##" \
      | sort
  } > "$output_path"
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

plain_inventory_entries() {
  local file="$1"
  awk '
    /^[[:space:]]*$/ { next }
    /^#/ { next }
    { print }
  ' "$file"
}

status_entries() {
  local file="$1"
  local status="$2"

  [ -f "$file" ] || return 0

  awk -v wanted="$status" '
    /^[a-zA-Z_][a-zA-Z0-9_-]*:[[:space:]]*$/ {
      current = $0
      sub(/:.*/, "", current)
      next
    }
    /^[a-zA-Z_][a-zA-Z0-9_-]*:[[:space:]]*\[\][[:space:]]*$/ {
      current = ""
      next
    }
    current == wanted && /^[[:space:]]*-[[:space:]]+/ {
      value = $0
      sub(/^[[:space:]]*-[[:space:]]+/, "", value)
      gsub(/^"|"$/, "", value)
      print value
    }
  ' "$file"
}

docs_marker_entries() {
  rg -o --no-filename '<!--[[:space:]]*upstream-docs:[[:space:]]*[^[:space:]]+[[:space:]]*-->' "$ROOT_PATH/services/docs" -g '*.md' 2>/dev/null \
    | sed -E 's/^<!--[[:space:]]*upstream-docs:[[:space:]]*//; s/[[:space:]]*-->$//' \
    | sort -u || true
}

generic_inventory_stats() {
  local inventory_file="$1"
  local ported_file="$2"
  local adapted_file="$3"
  local excluded_file="$4"

  awk '
    FILENAME == ARGV[1] {
      ported[$0] = 1
      next
    }
    FILENAME == ARGV[2] {
      adapted[$0] = 1
      next
    }
    FILENAME == ARGV[3] {
      excluded[$0] = 1
      next
    }
    /^[[:space:]]*$/ || /^#/ {
      next
    }
    {
      total++

      if ($0 in ported) {
        ported_count++
      } else if ($0 in adapted) {
        adapted_count++
      } else if ($0 in excluded) {
        excluded_count++
      } else {
        missing_count++
      }
    }
    END {
      printf "%d\t%d\t%d\t%d\t%d\n", total + 0, ported_count + 0, adapted_count + 0, excluded_count + 0, missing_count + 0
    }
  ' "$ported_file" "$adapted_file" "$excluded_file" "$inventory_file"
}

pending_plain_inventory_entries() {
  local inventory_file="$1"
  local ported_file="$2"
  local adapted_file="$3"
  local excluded_file="$4"

  awk '
    FILENAME == ARGV[1] {
      ported[$0] = 1
      next
    }
    FILENAME == ARGV[2] {
      adapted[$0] = 1
      next
    }
    FILENAME == ARGV[3] {
      excluded[$0] = 1
      next
    }
    /^[[:space:]]*$/ || /^#/ {
      next
    }
    !($0 in ported) && !($0 in adapted) && !($0 in excluded) {
      print
    }
  ' "$ported_file" "$adapted_file" "$excluded_file" "$inventory_file"
}

docs_pending_by_file_rows() {
  local rows
  rows="$(pending_plain_inventory_entries "$DOCS_INVENTORY_FILE" "$1" "$2" "$3" \
    | awk -F '#' '{ count[$1]++ } END { for (file in count) print count[file] "\t" file }' \
    | sort -nr \
    | head -n 20)"

  if [ -z "$rows" ]; then
    broadcastclient "| n/a | 0 |"
    return 0
  fi

  printf '%s\n' "$rows" | while IFS=$'\t' read -r count file; do
    printf '| `%s` | %s |\n' "$file" "$count"
  done
}

docs_status_by_file_rows() {
  local status_file="$1"
  local label="$2"
  local rows
  rows="$(awk -F '#' '{ count[$1]++ } END { for (file in count) print file "\t" count[file] }' "$status_file" | sort)"

  if [ -z "$rows" ]; then
    return 0
  fi

  printf '%s\n' "$rows" | while IFS=$'\t' read -r file count; do
    printf '| `%s` | %s | %s |\n' "$file" "$label" "$count"
  done
}

skeleton_status_rows() {
  local rows
  rows="$(
    {
      pending_plain_inventory_entries "$SKELETON_INVENTORY_FILE" "$1" "$2" "$3" | awk '{ print $0 "\tPending" }'
      awk '{ print $0 "\tAdapted" }' "$2"
      awk '{ print $0 "\tExcluded" }' "$3"
    } | sort
  )"

  if [ -z "$rows" ]; then
    broadcastclient "| n/a | n/a |"
    return 0
  fi

  printf '%s\n' "$rows" | while IFS=$'\t' read -r file status; do
    printf '| `%s` | %s |\n' "$file" "$status"
  done
}

generic_summary_row() {
  local label="$1"
  local inventory_count="$2"
  local inventory_file="$3"
  local ported_file="$4"
  local adapted_file="$5"
  local excluded_file="$6"
  local stats total ported adapted excluded missing

  stats="$(generic_inventory_stats "$inventory_file" "$ported_file" "$adapted_file" "$excluded_file")"
  IFS=$'\t' read -r total ported adapted excluded missing <<< "$stats"

  printf '| %s | %s | %s | %s | %s | %s | %s |\n' \
    "$label" \
    "$inventory_count" \
    "$total" \
    "$(count_percent_display "$ported" "$total")" \
    "$(count_percent_display "$missing" "$total")" \
    "$(count_percent_display "$adapted" "$total")" \
    "$(count_percent_display "$excluded" "$total")"
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

percent_display() {
  local count="$1"
  local total="$2"

  if [ "$total" -eq 0 ]; then
    printf 'n/a'
    return 0
  fi

  awk -v count="$count" -v total="$total" 'BEGIN { printf "%.1f%%", (count / total) * 100 }'
}

count_percent_display() {
  local count="$1"
  local total="$2"

  printf '%s (%s)' "$count" "$(percent_display "$count" "$total")"
}

inventory_summary_stats() {
  local scope="$1"
  local ported_index="$2"
  local adapted_index="$3"
  local inventory_count=0 total_count=0 ported_count=0 adapted_count=0 missing_count=0
  local id status repo branch tests_path inventory filter upstream bedrock

  while IFS="$RECORD_SEPARATOR" read -r id status repo branch tests_path inventory filter upstream bedrock; do
    [ -n "$inventory" ] && [ "$inventory" != "null" ] || continue

    case "$scope" in
      all) ;;
      framework) [[ "$id" == framework.* ]] || continue ;;
      package) [[ "$id" == package.* ]] || continue ;;
      *)
        broadcastclient "upstream-compliance: unknown inventory summary scope: $scope" >&2
        return 1
        ;;
    esac

    local stats total ported adapted missing
    stats="$(inventory_stats "$COMPLIANCE_PATH/$inventory" "$ported_index" "$adapted_index")"
    IFS=$'\t' read -r total ported adapted missing <<< "$stats"

    inventory_count=$((inventory_count + 1))
    total_count=$((total_count + total))
    ported_count=$((ported_count + ported))
    adapted_count=$((adapted_count + adapted))
    missing_count=$((missing_count + missing))
  done < <(list_records)

  printf '%d\t%d\t%d\t%d\t%d\n' "$inventory_count" "$total_count" "$ported_count" "$adapted_count" "$missing_count"
}

inventory_summary_row() {
  local label="$1"
  local scope="$2"
  local ported_index="$3"
  local adapted_index="$4"
  local stats inventory_count total ported adapted missing

  stats="$(inventory_summary_stats "$scope" "$ported_index" "$adapted_index")"
  IFS=$'\t' read -r inventory_count total ported adapted missing <<< "$stats"

  printf '| %s | %s | %s | %s | %s | %s |\n' \
    "$label" \
    "$inventory_count" \
    "$total" \
    "$(count_percent_display "$ported" "$total")" \
    "$(count_percent_display "$missing" "$total")" \
    "$(count_percent_display "$adapted" "$total")"
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

  printf 'Ported tests: %s / %s (%s); Missing tests: %s (%s); Adapted tests: %s (%s)' \
    "$ported" \
    "$total" \
    "$(percent_display "$ported" "$total")" \
    "$missing" \
    "$(percent_display "$missing" "$total")" \
    "$adapted" \
    "$(percent_display "$adapted" "$total")"
}

report() {
  need rg

  local output_file tmp_file
  output_file="${1:-$REPORT_FILE}"
  tmp_file="$(mktemp "${TMPDIR:-/tmp}/bedrock-compliance-report.XXXXXX")"
  mv "$tmp_file" "$tmp_file.md"
  tmp_file="$tmp_file.md"
  local status_index_path ported_index adapted_index docs_ported_index docs_adapted_index docs_excluded_index skeleton_ported_index skeleton_adapted_index skeleton_excluded_index
  status_index_path="$(mktemp -d "${TMPDIR:-/tmp}/bedrock-compliance-status.XXXXXX")"
  ported_index="$status_index_path/ported.txt"
  adapted_index="$status_index_path/adapted.txt"
  docs_ported_index="$status_index_path/docs-ported.txt"
  docs_adapted_index="$status_index_path/docs-adapted.txt"
  docs_excluded_index="$status_index_path/docs-excluded.txt"
  skeleton_ported_index="$status_index_path/skeleton-ported.txt"
  skeleton_adapted_index="$status_index_path/skeleton-adapted.txt"
  skeleton_excluded_index="$status_index_path/skeleton-excluded.txt"
  build_status_index "$ported_index" "$adapted_index"
  docs_marker_entries > "$docs_ported_index"
  status_entries "$DOCS_STATUS_FILE" adapted | sort -u > "$docs_adapted_index"
  status_entries "$DOCS_STATUS_FILE" excluded | sort -u > "$docs_excluded_index"
  status_entries "$SKELETON_STATUS_FILE" ported | sort -u > "$skeleton_ported_index"
  status_entries "$SKELETON_STATUS_FILE" adapted | sort -u > "$skeleton_adapted_index"
  status_entries "$SKELETON_STATUS_FILE" excluded | sort -u > "$skeleton_excluded_index"

  {
    broadcastclient "# Upstream Compliance Report"
    broadcastclient
    broadcastclient "Generated: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
    broadcastclient
    broadcastclient "Source of truth: services/compliance"
    broadcastclient
    broadcastclient "## Compliance Workflow Requirements"
    broadcastclient
    broadcastclient "- Every mapped Bedrock package compliance pass must include an upstream feature audit: \`make sure we also have all the upstream features for <package>\`."
    broadcastclient "- The audit compares upstream source contracts, public APIs, middleware, events, and runtime behavior against the Bedrock surface."
    broadcastclient "- Bedrock-equivalent features need executable parity tests; PHP-only or intentionally different behavior belongs in \`services/compliance/divergences.yml\`."
    broadcastclient "- \`services/compliance/features.yml\` must contain feature audit coverage for every mapped Bedrock surface."
    broadcastclient
    broadcastclient "## Inventories"
    broadcastclient
    broadcastclient "| Inventory | Upstream Tests | Ported Tests | Missing Tests | Adapted Tests |"
    broadcastclient "| --- | ---: | ---: | ---: | ---: |"

    list_records | while IFS="$RECORD_SEPARATOR" read -r id status repo branch tests_path inventory filter upstream bedrock; do
      [ -n "$inventory" ] && [ "$inventory" != "null" ] || continue

      local file
      file="$COMPLIANCE_PATH/$inventory"

      local rel stats total ported adapted missing
      rel="${file#$COMPLIANCE_PATH/}"
      stats="$(inventory_stats "$file" "$ported_index" "$adapted_index")"
      IFS=$'\t' read -r total ported adapted missing <<< "$stats"
      printf '| %s | %s | %s | %s | %s |\n' \
        "$rel" \
        "$total" \
        "$(count_percent_display "$ported" "$total")" \
        "$(count_percent_display "$missing" "$total")" \
        "$(count_percent_display "$adapted" "$total")"
    done

    broadcastclient
    broadcastclient "## Framework Coverage"
    broadcastclient
    broadcastclient "| Upstream Source | Bedrock Surface | Port Status | Tests | Docs |"
    broadcastclient "| --- | --- | --- | --- | --- |"

    list_records | while IFS="$RECORD_SEPARATOR" read -r id status repo branch tests_path inventory filter upstream bedrock; do
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
        "$(markdown_cell "$upstream")" \
        "$display_bedrock" \
        "$(markdown_cell "$status")" \
        "$(markdown_cell "$display_tests")" \
        "$display_docs"
    done

    broadcastclient
    broadcastclient "## Package Coverage"
    broadcastclient
    broadcastclient "| Upstream Package | Bedrock Surface | Port Status | Tests | Docs |"
    broadcastclient "| --- | --- | --- | --- | --- |"

    list_records | while IFS="$RECORD_SEPARATOR" read -r id status repo branch tests_path inventory filter upstream bedrock; do
      case "$id" in
        package.*) ;;
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
        "$(markdown_cell "$upstream")" \
        "$display_bedrock" \
        "$(markdown_cell "$status")" \
        "$(markdown_cell "$display_tests")" \
        "$display_docs"
    done

    broadcastclient
    broadcastclient "## Test Porting Summary"
    broadcastclient
    broadcastclient "| Scope | Inventories | Upstream Tests | Ported Tests | Pending / Missing Tests | Adapted Tests |"
    broadcastclient "| --- | ---: | ---: | ---: | ---: | ---: |"
    inventory_summary_row "All inventories" "all" "$ported_index" "$adapted_index"
    inventory_summary_row "Framework inventories" "framework" "$ported_index" "$adapted_index"
    inventory_summary_row "Package inventories" "package" "$ported_index" "$adapted_index"
    broadcastclient
    broadcastclient "## Documentation Porting Summary"
    broadcastclient
    broadcastclient "| Scope | Doc Inventories | Upstream Sections | Ported Sections | Pending / Missing Sections | Adapted Sections | Excluded Sections |"
    broadcastclient "| --- | ---: | ---: | ---: | ---: | ---: | ---: |"
    generic_summary_row "Upstream docs" "1" "$DOCS_INVENTORY_FILE" "$docs_ported_index" "$docs_adapted_index" "$docs_excluded_index"
    broadcastclient
    broadcastclient "## Upstream Skeleton Demo Summary"
    broadcastclient
    broadcastclient "| Scope | Skeleton Inventories | Upstream Files | Ported Files | Pending / Missing Files | Adapted Files | Excluded Files |"
    broadcastclient "| --- | ---: | ---: | ---: | ---: | ---: | ---: |"
    generic_summary_row "upstream/upstream" "1" "$SKELETON_INVENTORY_FILE" "$skeleton_ported_index" "$skeleton_adapted_index" "$skeleton_excluded_index"
    broadcastclient
    broadcastclient "## Documentation Detail"
    broadcastclient
    broadcastclient "### Top Pending Documentation Files"
    broadcastclient
    broadcastclient "| Upstream Docs File | Pending Sections |"
    broadcastclient "| --- | ---: |"
    docs_pending_by_file_rows "$docs_ported_index" "$docs_adapted_index" "$docs_excluded_index"
    broadcastclient
    broadcastclient "### Classified Documentation Sections"
    broadcastclient
    broadcastclient "| Upstream Docs File | Status | Sections |"
    broadcastclient "| --- | --- | ---: |"
    docs_status_by_file_rows "$docs_adapted_index" "Adapted"
    docs_status_by_file_rows "$docs_excluded_index" "Excluded"
    broadcastclient
    broadcastclient "## Skeleton Detail"
    broadcastclient
    broadcastclient "| Upstream Skeleton File | Status |"
    broadcastclient "| --- | --- |"
    skeleton_status_rows "$skeleton_ported_index" "$skeleton_adapted_index" "$skeleton_excluded_index"
    broadcastclient
    broadcastclient "## Feature Coverage"
    broadcastclient
    broadcastclient "| Feature | Source | Bedrock Surface | Status | Docs | Notes |"
    broadcastclient "| --- | --- | --- | --- | --- | --- |"

    list_feature_records | while IFS="$RECORD_SEPARATOR" read -r id source_id name upstream bedrock status docs notes; do
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

    broadcastclient
    broadcastclient "## Sources Without Test Inventories"
    broadcastclient
    broadcastclient "| Source | Bedrock | Tracking | Reason |"
    broadcastclient "| --- | --- | --- | --- |"

    list_records | while IFS="$RECORD_SEPARATOR" read -r id status repo branch tests_path inventory filter upstream bedrock; do
      [ "$status" = "mapped" ] || continue
      [ -z "$inventory" ] || [ "$inventory" = "null" ] || continue

      local tracking reason
      case "$id" in
        framework.contracts)
          tracking="Concrete package inventories"
          reason="Interface-only component; compliance is verified through concrete package inventories."
          ;;
        package.broadcastclient)
          tracking="Feature inventory"
          reason="TypeScript client package; compliance is tracked by feature coverage until a stable upstream test inventory is available."
          ;;
        package.httppreview)
          tracking="Feature inventory"
          reason="JavaScript package; compliance is tracked by feature coverage until a stable upstream test inventory is available."
          ;;
        *)
          tracking="Mapped source"
          reason="No generated upstream test inventory configured."
          ;;
      esac

      printf '| `%s` | `%s` | %s | %s |\n' "$upstream" "$bedrock" "$tracking" "$reason"
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

  format_report_markdown "$tmp_file"

  rm -rf "$status_index_path"
  mv "$tmp_file" "$output_file"
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
      framework.collections) candidate="collection" ;;
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
  local inline_exclusions

  if inline_exclusions="$(rg -n '^[[:space:]]*// Excluded:' "$ROOT_PATH/packages" -g '*_test.go' 2>/dev/null)"; then
    broadcastclient "Inline test exclusions are not allowed; put exclusions in services/compliance:" >&2
    printf '%s\n' "$inline_exclusions" >&2
    return 1
  fi
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

check_plain_inventory_format() {
  local file="$1"
  local label="$2"

  if [ ! -f "$file" ]; then
    broadcastclient "Missing $label inventory file: ${file#$ROOT_PATH/}" >&2
    return 1
  fi

  awk '
    /^[[:space:]]*$/ || /^#/ { next }
    /^[^[:space:]]+$/ { next }
    { invalid = 1; print FILENAME ":" FNR ": invalid inventory entry: " $0 > "/dev/stderr" }
    END { exit invalid ? 1 : 0 }
  ' "$file"
}

check_status_entries_known() {
  local status_file="$1"
  local inventory_file="$2"
  local label="$3"
  local failed=0 entries known unknown

  entries="$(mktemp "${TMPDIR:-/tmp}/bedrock-compliance-status-entries.XXXXXX")"
  known="$(mktemp "${TMPDIR:-/tmp}/bedrock-compliance-known-entries.XXXXXX")"

  {
    status_entries "$status_file" ported
    status_entries "$status_file" adapted
    status_entries "$status_file" excluded
  } | sort -u > "$entries"

  plain_inventory_entries "$inventory_file" | sort -u > "$known"

  unknown="$(comm -23 "$entries" "$known")"

  if [ -n "$unknown" ]; then
    broadcastclient "$label status file references unknown inventory entries:" >&2
    printf '%s\n' "$unknown" >&2
    failed=1
  fi

  rm -f "$entries" "$known"
  return "$failed"
}

check_docs_markers_known() {
  local failed=0 markers known unknown

  markers="$(mktemp "${TMPDIR:-/tmp}/bedrock-compliance-doc-markers.XXXXXX")"
  known="$(mktemp "${TMPDIR:-/tmp}/bedrock-compliance-doc-known.XXXXXX")"

  docs_marker_entries > "$markers"
  plain_inventory_entries "$DOCS_INVENTORY_FILE" | sort -u > "$known"

  unknown="$(comm -23 "$markers" "$known")"

  if [ -n "$unknown" ]; then
    broadcastclient "Documentation markers reference unknown Upstream docs inventory entries:" >&2
    printf '%s\n' "$unknown" >&2
    failed=1
  fi

  rm -f "$markers" "$known"
  return "$failed"
}

check_docs_and_skeleton_tracking() {
  local failed=0

  [ -f "$DOCS_STATUS_FILE" ] || { broadcastclient "Missing $DOCS_STATUS_FILE" >&2; failed=1; }
  [ -f "$SKELETON_STATUS_FILE" ] || { broadcastclient "Missing $SKELETON_STATUS_FILE" >&2; failed=1; }

  check_plain_inventory_format "$DOCS_INVENTORY_FILE" "docs" || failed=1
  check_plain_inventory_format "$SKELETON_INVENTORY_FILE" "skeleton" || failed=1

  if [ "$failed" -eq 0 ]; then
    check_status_entries_known "$DOCS_STATUS_FILE" "$DOCS_INVENTORY_FILE" "Docs" || failed=1
    check_status_entries_known "$SKELETON_STATUS_FILE" "$SKELETON_INVENTORY_FILE" "Skeleton" || failed=1
    check_docs_markers_known || failed=1
  fi

  return "$failed"
}

check_features_format() {
  awk '
    /^[[:space:]]*$/ || /^#/ { next }
    /^version:[ \t]*[0-9]+[ \t]*$/ { next }
    /^features:[ \t]*$/ { next }
    /^  - id:[ \t]*[^[:space:]].*$/ { next }
    /^    (source_id|name|upstream|bedrock|status|docs|notes):[ \t]*.*$/ { next }
    { invalid = 1; print FILENAME ":" FNR ": invalid features entry: " $0 > "/dev/stderr" }
    END { exit invalid ? 1 : 0 }
  ' "$FEATURES_FILE"
}

check_features() {
  local failed=0 count=0 source_ids feature_ids required_source_ids audited_source_ids missing_audits
  source_ids="$(mktemp "${TMPDIR:-/tmp}/bedrock-compliance-source-ids.XXXXXX")"
  feature_ids="$(mktemp "${TMPDIR:-/tmp}/bedrock-compliance-feature-ids.XXXXXX")"
  required_source_ids="$(mktemp "${TMPDIR:-/tmp}/bedrock-compliance-required-feature-audits.XXXXXX")"
  audited_source_ids="$(mktemp "${TMPDIR:-/tmp}/bedrock-compliance-audited-sources.XXXXXX")"

  list_records | awk -F "$RECORD_SEPARATOR" '{ print $1 }' | sort -u > "$source_ids"
  list_records | awk -F "$RECORD_SEPARATOR" '$2 == "mapped" && $9 != "" && $9 != "null" { print $1 }' | sort -u > "$required_source_ids"
  : > "$feature_ids"
  : > "$audited_source_ids"

  check_features_format || failed=1

  while IFS="$RECORD_SEPARATOR" read -r id source_id name upstream bedrock status docs notes; do
    count=$((count + 1))

    if [ -z "$id" ] || [ -z "$source_id" ] || [ -z "$name" ] || [ -z "$upstream" ] || [ -z "$bedrock" ] || [ -z "$status" ]; then
      broadcastclient "Feature $id is missing a required field." >&2
      failed=1
    fi

    case "$status" in
      ported|partial|missing|excluded) ;;
      *)
        broadcastclient "Feature $id has invalid status: $status" >&2
        failed=1
        ;;
    esac

    if ! rg -Fxq -- "$source_id" "$source_ids"; then
      broadcastclient "Feature $id references unknown source_id: $source_id" >&2
      failed=1
    fi
    printf '%s\n' "$source_id" >> "$audited_source_ids"

    if rg -Fxq -- "$id" "$feature_ids"; then
      broadcastclient "Duplicate feature id: $id" >&2
      failed=1
    fi
    printf '%s\n' "$id" >> "$feature_ids"

    if [ -n "$docs" ] && [ "$docs" != "null" ]; then
      case "$docs" in
        services/docs/*|packages/*|basics/*|concepts/*) ;;
        *)
          broadcastclient "Feature $id docs path must be under services/docs: $docs" >&2
          failed=1
          ;;
      esac

      if [ "$status" = "ported" ] || [ "$status" = "partial" ]; then
        local docs_file
        docs_file="$(docs_file_for "$bedrock" "$docs")"

        if [ ! -f "$docs_file" ]; then
          broadcastclient "Feature $id docs path does not exist: ${docs_file#$ROOT_PATH/}" >&2
          failed=1
        fi
      fi
    fi
  done < <(list_feature_records)

  sort -u -o "$audited_source_ids" "$audited_source_ids"
  missing_audits="$(comm -23 "$required_source_ids" "$audited_source_ids")"

  if [ -n "$missing_audits" ]; then
    broadcastclient "Mapped Bedrock sources missing upstream feature audit coverage in $FEATURES_FILE:" >&2
    printf '%s\n' "$missing_audits" >&2
    failed=1
  fi

  if [ "$count" -eq 0 ]; then
    broadcastclient "$FEATURES_FILE must define at least one feature." >&2
    failed=1
  fi

  rm -f "$source_ids" "$feature_ids" "$required_source_ids" "$audited_source_ids"
  return "$failed"
}

normalize_report() {
  sed 's/^Generated: .*/Generated: <normalized>/'
}

format_report_markdown() {
  local file="$1"

  if command -v pnpm >/dev/null 2>&1 && [ -f "$ROOT_PATH/package.json" ]; then
    (cd "$ROOT_PATH" && pnpm exec oxfmt --ignore-path .gitignore "$file" >/dev/null)
  fi
}

check_report_current() {
  local expected actual failed=0
  expected="$(mktemp "${TMPDIR:-/tmp}/bedrock-compliance-report-expected.XXXXXX")"
  actual="$(mktemp "${TMPDIR:-/tmp}/bedrock-compliance-report-actual.XXXXXX")"

  report "$expected"

  normalize_report < "$expected" > "$expected.normalized"
  normalize_report < "$REPORT_FILE" > "$actual"

  if ! diff -u "$expected.normalized" "$actual" >&2; then
    broadcastclient "services/compliance/report.md is stale; run services/scripts/upstream-compliance.sh report." >&2
    failed=1
  fi

  rm -f "$expected" "$expected.normalized" "$actual"
  return "$failed"
}

check() {
  need rg

  [ -f "$INVENTORY_FILE" ] || { broadcastclient "Missing $INVENTORY_FILE" >&2; return 1; }
  [ -f "$DIVERGENCES_FILE" ] || { broadcastclient "Missing $DIVERGENCES_FILE" >&2; return 1; }
  [ -f "$FEATURES_FILE" ] || { broadcastclient "Missing $FEATURES_FILE" >&2; return 1; }
  [ -f "$SOURCES_LOCK_FILE" ] || { broadcastclient "Missing $SOURCES_LOCK_FILE" >&2; return 1; }
  [ -f "$DOCS_STATUS_FILE" ] || { broadcastclient "Missing $DOCS_STATUS_FILE" >&2; return 1; }
  [ -f "$SKELETON_STATUS_FILE" ] || { broadcastclient "Missing $SKELETON_STATUS_FILE" >&2; return 1; }

  check_no_package_local_parity
  check_excluded_not_implemented
  check_laravel_test_files_have_source_context
  check_no_inline_exclusions
  check_inventory_files_configured
  check_inventory_format
  check_docs_and_skeleton_tracking
  check_features
  check_report_current
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
