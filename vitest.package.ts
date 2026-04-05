import path from "node:path";
import { fileURLToPath } from "node:url";

import { defineConfig } from "vitest/config";

const repoRoot = path.dirname(fileURLToPath(import.meta.url));
const packageRoot = process.cwd();
const packageRelativePath = path.relative(repoRoot, packageRoot) || ".";
const normalizedPackagePath = packageRelativePath.split(path.sep).join("/");

export default defineConfig({
  resolve: {
    alias: {
      "@": path.resolve(packageRoot, "src"),
    },
  },
  root: packageRoot,
  test: {
    environment: "node",
    passWithNoTests: true,
    coverage: {
      provider: "v8",
      reporter: [
        "text",
        "json-summary",
        "lcov",
        ["lcovonly", { file: "coverage.out" }],
      ],
      reportsDirectory: path.join(
        repoRoot,
        "storage",
        "coverage",
        normalizedPackagePath,
      ),
    },
  },
});
