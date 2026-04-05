const userAgent = process.env.npm_config_user_agent ?? "";

if (!userAgent.includes("pnpm/")) {
  console.error("This repository uses pnpm only.");
  console.error("Run `corepack enable` if needed, then `pnpm install`.");
  process.exit(1);
}
