import fs from "node:fs";
import path from "node:path";

const root = path.resolve(new URL("../", import.meta.url).pathname);
const resources = path.join(root, "resources");
const output = path.join(root, "public", "build");

fs.mkdirSync(output, { recursive: true });

for (const [source, destination] of [
  [path.join(resources, "css", "app.css"), path.join(output, "app.css")],
  [path.join(resources, "js", "app.js"), path.join(output, "app.js")],
]) {
  fs.copyFileSync(source, destination);
}
