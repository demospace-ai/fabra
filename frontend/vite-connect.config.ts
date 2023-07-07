import react from "@vitejs/plugin-react-swc";
import { defineConfig, splitVendorChunkPlugin } from "vite";
import viteTsconfigPaths from "vite-tsconfig-paths";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react(), viteTsconfigPaths(), splitVendorChunkPlugin()],
  build: {
    rollupOptions: {
      input: "connect.html",
      output: [
        {
          dir: "build",
        },
      ],
    },
  },
  server: {
    port: 3000,
  },
});
