import fs from "node:fs/promises";
import path from "node:path";
import process from "node:process";
import { spawn } from "node:child_process";

const repoRoot = process.cwd();
const [, , task, ...extraArgs] = process.argv;

if (!task) {
  console.error("Usage: node ./scripts/run-turbo.mjs <task> [extra args]");
  process.exit(1);
}

const turboArgs = [
  "turbo",
  "run",
  task,
  "--cache-dir=./storage/.turbo",
  ...extraArgs,
];

const child = spawn("pnpm", turboArgs, {
  cwd: repoRoot,
  stdio: "inherit",
});

const cleanupWorkspaceTurboDirs = async () => {
  for (const workspaceGroup of ["apps", "packages"]) {
    const workspaceGroupPath = path.join(repoRoot, workspaceGroup);

    let entries = [];

    try {
      entries = await fs.readdir(workspaceGroupPath, { withFileTypes: true });
    } catch {
      continue;
    }

    await Promise.all(
      entries
        .filter((entry) => entry.isDirectory())
        .map((entry) =>
          fs.rm(path.join(workspaceGroupPath, entry.name, ".turbo"), {
            force: true,
            recursive: true,
          }),
        ),
    );
  }
};

child.on("close", async (code, signal) => {
  await cleanupWorkspaceTurboDirs();

  if (signal) {
    process.kill(process.pid, signal);
    return;
  }

  process.exit(code ?? 1);
});
