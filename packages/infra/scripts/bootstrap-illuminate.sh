#!/bin/sh
set -e

ROOT_DIR="$(cd "$(dirname "$0")/../../.." && pwd)"
PACKAGES_DIR="${ROOT_DIR}/packages/anvil"

COMPONENTS="
auth
broadcasting
bus
cache
collections
concurrency
conditionable
config
console
container
contracts
cookie
database
encryption
events
filesystem
foundation
hashing
http
json-schema
log
macroable
mail
notifications
pagination
pipeline
process
queue
redis
reflection
routing
session
support
testing
translation
validation
view
"

package_name() {
  case "$1" in
    json-schema) echo "jsonschema" ;;
    *)           echo "$1" | tr -d '-' ;;
  esac
}

title_case() {
  echo "$1" | tr '-' ' ' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) substr($i,2)}1'
}

mkdir -p "${PACKAGES_DIR}"

for component in ${COMPONENTS}; do
  dir="${PACKAGES_DIR}/${component}"
  pkg=$(package_name "${component}")
  title=$(title_case "${component}")

  mkdir -p "${dir}"

  if [ ! -f "${dir}/go.mod" ]; then
    cat > "${dir}/go.mod" <<EOF
module github.com/bedrock/packages/anvil/${component}

go 1.26.0
EOF
  fi

  if [ ! -f "${dir}/doc.go" ]; then
    lower_title=$(echo "${title}" | tr '[:upper:]' '[:lower:]')
    cat > "${dir}/doc.go" <<EOF
// Package ${pkg} provides Laravel-inspired ${lower_title} primitives.
package ${pkg}
EOF
  fi

  if [ ! -f "${dir}/package.json" ]; then
    cat > "${dir}/package.json" <<EOF
{
  "name": "@bedrock/illuminate-${component}",
  "version": "0.0.0",
  "private": true,
  "scripts": {
    "build": "sh -c 'mkdir -p ../../../storage/dist/${component} && go build ./... && touch ../../../storage/dist/${component}/.build-stamp'",
    "dev": "go test ./... -count=1",
    "test": "go test ./...",
    "test:coverage": "sh -c 'mkdir -p ../../../storage/coverage/go/${component} && go test ./... -coverprofile=../../../storage/coverage/go/${component}/coverage.out'",
    "typecheck": "go test ./... -run '^\$'",
    "fmt": "gofmt -w .",
    "fmt:check": "sh -c 'test -z \"\$(gofmt -l .)\"'"
  }
}
EOF
  fi
done
