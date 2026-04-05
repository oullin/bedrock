import fs from "node:fs";
import path from "node:path";

const root = path.resolve(new URL("../../", import.meta.url).pathname);
const packagesRoot = path.join(root, "packages", "illuminate");

const components = [
  "auth",
  "broadcasting",
  "bus",
  "cache",
  "collections",
  "concurrency",
  "conditionable",
  "config",
  "console",
  "container",
  "contracts",
  "cookie",
  "database",
  "encryption",
  "events",
  "filesystem",
  "foundation",
  "hashing",
  "http",
  "json-schema",
  "log",
  "macroable",
  "mail",
  "notifications",
  "pagination",
  "pipeline",
  "process",
  "queue",
  "redis",
  "reflection",
  "routing",
  "session",
  "support",
  "testing",
  "translation",
  "validation",
  "view",
];

const packageNames = {
  "json-schema": "jsonschema",
};

fs.mkdirSync(packagesRoot, { recursive: true });

for (const component of components) {
  const dir = path.join(packagesRoot, component);
  const pkg = packageNames[component] ?? component.replaceAll("-", "");
  const title = component
    .split("-")
    .map((segment) => segment.charAt(0).toUpperCase() + segment.slice(1))
    .join(" ");

  fs.mkdirSync(dir, { recursive: true });

  const goMod = path.join(dir, "go.mod");
  if (!fs.existsSync(goMod)) {
    fs.writeFileSync(
      goMod,
      `module github.com/gollin/packages/illuminate/${component}

go 1.26.0
`,
    );
  }

  const doc = path.join(dir, "doc.go");
  if (!fs.existsSync(doc)) {
    fs.writeFileSync(
      doc,
      `// Package ${pkg} provides Laravel-inspired ${title.toLowerCase()} primitives.
package ${pkg}
`,
    );
  }

  const packageJson = path.join(dir, "package.json");
  if (!fs.existsSync(packageJson)) {
    fs.writeFileSync(
      packageJson,
      JSON.stringify(
        {
          name: `@gollin/illuminate-${component}`,
          version: "0.0.0",
          private: true,
          scripts: {
            build: `sh -c 'mkdir -p ../../storage/dist/${component} && go build ./... && touch ../../storage/dist/${component}/.build-stamp'`,
            dev: "go test ./... -count=1",
            test: "go test ./...",
            "test:coverage": `sh -c 'mkdir -p ../../storage/coverage/go/${component} && go test ./... -coverprofile=../../storage/coverage/go/${component}/coverage.out'`,
            typecheck: "go test ./... -run '^$'",
            fmt: "gofmt -w .",
            "fmt:check": `sh -c 'test -z "$(gofmt -l .)"'`,
          },
        },
        null,
        2,
      ) + "\n",
    );
  }
}
