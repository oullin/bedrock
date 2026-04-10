#!/bin/sh
set -e

case "${npm_config_user_agent:-}" in
  *pnpm/*) ;;
  *)
    echo "This repository uses pnpm only." >&2
    echo "Run \`corepack enable\` if needed, then \`pnpm install\`." >&2
    exit 1
    ;;
esac
