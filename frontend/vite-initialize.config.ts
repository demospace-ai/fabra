import react from "@vitejs/plugin-react-swc";
import { defineConfig, splitVendorChunkPlugin } from "vite";
import viteTsconfigPaths from "vite-tsconfig-paths";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react(), viteTsconfigPaths(), splitVendorChunkPlugin()],
  build: {
    lib: {
      entry: "src/initialize.ts",
    },
    rollupOptions: {
      output: [
        {
          dir: "build",
          entryFileNames: "initialize.js",
        },
      ],
    },
  },
  server: {
    port: 3000,
  },
});
