import { fileURLToPath, URL } from "node:url";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import vueJsx from "@vitejs/plugin-vue-jsx";
import unocss from "unocss/vite";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue(), vueJsx(), unocss()],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url))
    }
  },
  server: {
    proxy: {
      "/api/notifications/listen": {
        target: "ws://localhost:8444",
        ws: true
      },

      "/api": "http://localhost:8444",
      "/.well-known": "http://localhost:8444"
    }
  }
});
