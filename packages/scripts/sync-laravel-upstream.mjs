import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { execFileSync } from "node:child_process";

const root = path.resolve(new URL("../../", import.meta.url).pathname);
const outputDir = path.join(root, "docs", "parity");
const frameworkTag = process.argv[2] || "v13.3.0";
const skeletonTag = process.argv[3] || "v13.1.2";

fs.mkdirSync(outputDir, { recursive: true });

function cloneRepo(url, tag, target) {
  execFileSync("git", ["clone", "--depth=1", "--branch", tag, url, target], {
    stdio: "ignore",
  });
}

function listFrameworkComponents(baseDir) {
  const illuminateDir = path.join(baseDir, "src", "Illuminate");
  const components = fs
    .readdirSync(illuminateDir, { withFileTypes: true })
    .filter((entry) => entry.isDirectory())
    .map((entry) => entry.name)
    .sort();

  return components.map((component) => {
    const componentDir = path.join(illuminateDir, component);
    const composerPath = path.join(componentDir, "composer.json");
    const composer = fs.existsSync(composerPath)
      ? JSON.parse(fs.readFileSync(composerPath, "utf8"))
      : {};
    const phpFiles = walk(componentDir).filter((file) => file.endsWith(".php"));

    return {
      component,
      package: composer.name ?? `illuminate/${component.toLowerCase()}`,
      requires: Object.keys(composer.require || {})
        .filter((dependency) => dependency.startsWith("illuminate/"))
        .sort(),
      suggests: Object.keys(composer.suggest || {}).sort(),
      phpFiles: phpFiles.length,
      phpLines: phpFiles.reduce(
        (total, file) => total + fs.readFileSync(file, "utf8").split("\n").length,
        0,
      ),
    };
  });
}

function listSkeletonFiles(baseDir) {
  return walk(baseDir)
    .map((file) => path.relative(baseDir, file))
    .filter((file) => !file.startsWith(".git/"))
    .sort();
}

function localComponentStatus() {
  const componentsDir = path.join(root, "packages", "illuminate");
  if (!fs.existsSync(componentsDir)) {
    return [];
  }

  return fs
    .readdirSync(componentsDir, { withFileTypes: true })
    .filter((entry) => entry.isDirectory())
    .map((entry) => {
      const dir = path.join(componentsDir, entry.name);
      const files = walk(dir);
      return {
        component: entry.name,
        files: files.length,
        hasDoc: fs.existsSync(path.join(dir, "doc.go")),
        hasGoMod: fs.existsSync(path.join(dir, "go.mod")),
        hasPackageJson: fs.existsSync(path.join(dir, "package.json")),
      };
    })
    .sort((left, right) => left.component.localeCompare(right.component));
}

function walk(dir) {
  const files = [];

  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const file = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      files.push(...walk(file));
      continue;
    }
    files.push(file);
  }

  return files;
}

const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), "laravel-sync-"));
const frameworkDir = path.join(tempDir, "framework");
const skeletonDir = path.join(tempDir, "laravel");

cloneRepo("https://github.com/laravel/framework.git", frameworkTag, frameworkDir);
cloneRepo("https://github.com/laravel/laravel.git", skeletonTag, skeletonDir);

const frameworkManifest = {
  source: "https://github.com/laravel/framework",
  tag: frameworkTag,
  generatedAt: new Date().toISOString(),
  components: listFrameworkComponents(frameworkDir),
};

const skeletonManifest = {
  source: "https://github.com/laravel/laravel",
  tag: skeletonTag,
  generatedAt: new Date().toISOString(),
  files: listSkeletonFiles(skeletonDir),
};

const localManifest = {
  generatedAt: new Date().toISOString(),
  components: localComponentStatus(),
};

fs.writeFileSync(
  path.join(outputDir, `framework-${frameworkTag}.json`),
  JSON.stringify(frameworkManifest, null, 2) + "\n",
);
fs.writeFileSync(
  path.join(outputDir, `laravel-${skeletonTag}.json`),
  JSON.stringify(skeletonManifest, null, 2) + "\n",
);
fs.writeFileSync(
  path.join(outputDir, "local-components.json"),
  JSON.stringify(localManifest, null, 2) + "\n",
);

const markdown = `# Laravel Parity Snapshot

- Framework tag: \`${frameworkTag}\`
- Skeleton tag: \`${skeletonTag}\`
- Generated: \`${new Date().toISOString()}\`

## Framework Components

| Component | Package | Illuminate Dependencies | PHP Files | PHP Lines |
| --- | --- | --- | ---: | ---: |
${frameworkManifest.components
  .map(
    (component) =>
      `| ${component.component} | \`${component.package}\` | ${
        component.requires.length ? component.requires.join(", ") : "-" 
      } | ${component.phpFiles} | ${component.phpLines} |`,
  )
  .join("\n")}

## Local Components

| Component | Files | go.mod | package.json | doc.go |
| --- | ---: | --- | --- | --- |
${localManifest.components
  .map(
    (component) =>
      `| ${component.component} | ${component.files} | ${component.hasGoMod ? "yes" : "no"} | ${component.hasPackageJson ? "yes" : "no"} | ${component.hasDoc ? "yes" : "no"} |`,
  )
  .join("\n")}
`;

fs.writeFileSync(path.join(outputDir, "README.md"), markdown);
