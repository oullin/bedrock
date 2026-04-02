import path from "node:path";
import { fileURLToPath } from "node:url";

import { defineConfig } from "vitest/config";

const repoRoot = path.dirname(fileURLToPath(import.meta.url));
const workspaceRoot = process.cwd();
const workspaceRelativePath = path.relative(repoRoot, workspaceRoot) || ".";
const normalizedWorkspacePath = workspaceRelativePath.split(path.sep).join("/");

export default defineConfig({
  resolve: {
    alias: {
      "@": path.resolve(workspaceRoot, "src"),
    },
  },
  root: workspaceRoot,
  test: {
    environment: "node",
    passWithNoTests: true,
    coverage: {
      provider: "v8",
      reporter: ["text", "json-summary", "lcov"],
      reportsDirectory: path.join(
        repoRoot,
        "storage",
        "coverage",
        normalizedWorkspacePath,
      ),
    },
  },
});
