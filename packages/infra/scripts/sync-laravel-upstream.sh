#!/bin/sh
set -e

ROOT_DIR="$(cd "$(dirname "$0")/../../.." && pwd)"
OUTPUT_DIR="${ROOT_DIR}/packages/docs/parity"
FRAMEWORK_TAG="${1:-v13.3.0}"
SKELETON_TAG="${2:-v13.1.2}"

mkdir -p "${OUTPUT_DIR}"

TEMP_DIR=$(mktemp -d)
FRAMEWORK_DIR="${TEMP_DIR}/framework"
SKELETON_DIR="${TEMP_DIR}/laravel"

cleanup() { rm -rf "${TEMP_DIR}"; }
trap cleanup EXIT

git clone --depth=1 --branch "${FRAMEWORK_TAG}" https://github.com/laravel/framework.git "${FRAMEWORK_DIR}" 2>/dev/null
git clone --depth=1 --branch "${SKELETON_TAG}" https://github.com/laravel/laravel.git "${SKELETON_DIR}" 2>/dev/null

GENERATED_AT=$(date -u +"%Y-%m-%dT%H:%M:%S.000Z")
ILLUMINATE_DIR="${FRAMEWORK_DIR}/src/Illuminate"

# --- Framework manifest ---

framework_components_json() {
  local first=true
  echo "["

  for component_dir in "${ILLUMINATE_DIR}"/*/; do
    [ -d "${component_dir}" ] || continue
    component=$(basename "${component_dir}")
    composer_path="${component_dir}/composer.json"

    if [ -f "${composer_path}" ] && command -v jq >/dev/null 2>&1; then
      pkg=$(jq -r '.name // empty' "${composer_path}" 2>/dev/null)
      requires=$(jq -r '(.require // {}) | keys[] | select(startswith("illuminate/"))' "${composer_path}" 2>/dev/null | sort | jq -R . | jq -s .)
      suggests=$(jq -r '(.suggest // {}) | keys[]' "${composer_path}" 2>/dev/null | sort | jq -R . | jq -s .)
    else
      pkg=""
      requires="[]"
      suggests="[]"
    fi

    [ -z "${pkg}" ] && pkg="illuminate/$(echo "${component}" | tr '[:upper:]' '[:lower:]')"

    php_files=$(find "${component_dir}" -name '*.php' -type f | wc -l | tr -d ' ')
    php_lines=$(find "${component_dir}" -name '*.php' -type f -exec cat {} + 2>/dev/null | wc -l | tr -d ' ')
    [ -z "${php_lines}" ] && php_lines=0

    if [ "${first}" = true ]; then
      first=false
    else
      echo ","
    fi

    cat <<ENTRY
    {
      "component": "${component}",
      "package": "${pkg}",
      "requires": ${requires},
      "suggests": ${suggests},
      "phpFiles": ${php_files},
      "phpLines": ${php_lines}
    }
ENTRY
  done

  echo ""
  echo "  ]"
}

cat > "${OUTPUT_DIR}/framework-${FRAMEWORK_TAG}.json" <<EOF
{
  "source": "https://github.com/laravel/framework",
  "tag": "${FRAMEWORK_TAG}",
  "generatedAt": "${GENERATED_AT}",
  "components": $(framework_components_json)
}
EOF

# --- Skeleton manifest ---

skeleton_files_json() {
  local first=true
  echo "["

  find "${SKELETON_DIR}" -type f | while IFS= read -r file; do
    rel=$(echo "${file}" | sed "s|^${SKELETON_DIR}/||")
    case "${rel}" in
      .git/*) continue ;;
    esac

    if [ "${first}" = true ]; then
      first=false
    else
      printf ","
    fi

    printf "\n    \"%s\"" "${rel}"
  done

  echo ""
  echo "  ]"
}

cat > "${OUTPUT_DIR}/laravel-${SKELETON_TAG}.json" <<EOF
{
  "source": "https://github.com/laravel/laravel",
  "tag": "${SKELETON_TAG}",
  "generatedAt": "${GENERATED_AT}",
  "files": $(skeleton_files_json)
}
EOF

# --- Local components manifest ---

COMPONENTS_DIR="${ROOT_DIR}/packages/anvil"

local_components_json() {
  local first=true
  echo "["

  if [ -d "${COMPONENTS_DIR}" ]; then
    for component_dir in "${COMPONENTS_DIR}"/*/; do
      [ -d "${component_dir}" ] || continue
      component=$(basename "${component_dir}")

      file_count=$(find "${component_dir}" -type f | wc -l | tr -d ' ')
      has_doc=$([ -f "${component_dir}/doc.go" ] && echo "true" || echo "false")
      has_go_mod=$([ -f "${component_dir}/go.mod" ] && echo "true" || echo "false")
      has_package_json=$([ -f "${component_dir}/package.json" ] && echo "true" || echo "false")

      if [ "${first}" = true ]; then
        first=false
      else
        echo ","
      fi

      cat <<ENTRY
    {
      "component": "${component}",
      "files": ${file_count},
      "hasDoc": ${has_doc},
      "hasGoMod": ${has_go_mod},
      "hasPackageJson": ${has_package_json}
    }
ENTRY
    done
  fi

  echo ""
  echo "  ]"
}

cat > "${OUTPUT_DIR}/local-components.json" <<EOF
{
  "generatedAt": "${GENERATED_AT}",
  "components": $(local_components_json)
}
EOF

# --- Markdown report ---

# Read the framework JSON to build the table (requires jq)
if command -v jq >/dev/null 2>&1; then
  framework_table=$(jq -r '.components[] | "| \(.component) | `\(.package)` | \(if (.requires | length) > 0 then (.requires | join(", ")) else "-" end) | \(.phpFiles) | \(.phpLines) |"' "${OUTPUT_DIR}/framework-${FRAMEWORK_TAG}.json")
  local_table=$(jq -r '.components[] | "| \(.component) | \(.files) | \(if .hasGoMod then "yes" else "no" end) | \(if .hasPackageJson then "yes" else "no" end) | \(if .hasDoc then "yes" else "no" end) |"' "${OUTPUT_DIR}/local-components.json")
else
  framework_table="| (jq required for table generation) | | | | |"
  local_table="| (jq required for table generation) | | | | |"
fi

cat > "${OUTPUT_DIR}/README.md" <<EOF
# Laravel Parity Snapshot

- Framework tag: \`${FRAMEWORK_TAG}\`
- Skeleton tag: \`${SKELETON_TAG}\`
- Generated: \`${GENERATED_AT}\`

## Framework Components

| Component | Package | Illuminate Dependencies | PHP Files | PHP Lines |
| --- | --- | --- | ---: | ---: |
${framework_table}

## Local Components

| Component | Files | go.mod | package.json | doc.go |
| --- | ---: | --- | --- | --- |
${local_table}
EOF
