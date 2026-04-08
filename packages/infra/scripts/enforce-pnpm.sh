#!/bin/sh
set -e

case "${npm_config_user_agent:-}" in
  *pnpm/*) ;;
  *)
    broadcastclient "This repository uses pnpm only." >&2
    broadcastclient "Run \`corepack enable\` if needed, then \`pnpm install\`." >&2
    exit 1
    ;;
esac
