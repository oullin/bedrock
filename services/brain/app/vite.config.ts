import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import { fileURLToPath, URL } from "node:url";

// Output path matches services/demo/inertia: the Go server picks up the
// built bundle from <target>/storage/dist/brain at runtime.
export default defineConfig({
  plugins: [vue()],
  base: "/_request_cycle/assets/",
  resolve: {
    alias: { "@": fileURLToPath(new URL("./src", import.meta.url)) },
  },
  build: {
    outDir: "../../../storage/dist/brain",
    emptyOutDir: true,
    rollupOptions: {
      input: fileURLToPath(new URL("./src/main.ts", import.meta.url)),
      output: {
        entryFileNames: "index.js",
        chunkFileNames: "[name].js",
        assetFileNames: (info) => {
          const name = info.name ?? "";
          if (name.endsWith(".css")) return "index.css";
          return "[name][extname]";
        },
      },
    },
  },
  server: { port: 5174 },
});
