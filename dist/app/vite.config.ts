import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import tailwindcss from "@tailwindcss/vite";
import { resolve } from "path";

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: {
      "@": resolve(__dirname, "resources/js"),
    },
  },
  build: {
    outDir: "public/build",
    manifest: true,
    rollupOptions: {
      input: "resources/js/app.ts",
    },
  },
  server: {
    origin: "http://localhost:5173",
  },
});
