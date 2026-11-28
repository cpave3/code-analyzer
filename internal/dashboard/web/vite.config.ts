import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";

// `codeeval serve` embeds dist/ (go:embed). In development, run the Go
// server for the API on its default address (task web:dev runs both) and
// let Vite proxy /api to it.
export default defineConfig({
  plugins: [react()],
  base: "/",
  build: {
    outDir: "dist",
    emptyOutDir: true,
    chunkSizeWarningLimit: 900,
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes("node_modules/echarts") || id.includes("node_modules/zrender")) return "echarts";
          if (id.includes("node_modules/sigma") || id.includes("node_modules/graphology")) return "graph";
          return undefined;
        },
      },
    },
  },
  server: {
    port: 5173,
    proxy: { "/api": "http://127.0.0.1:8765" },
  },
  test: {
    environment: "jsdom",
    globals: true,
    include: ["src/**/*.test.ts", "src/**/*.test.tsx"],
  },
});
