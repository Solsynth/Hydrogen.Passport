import { defineConfig } from 'vite'
import path from "path";
import react from '@vitejs/plugin-react-swc'
import unocss from "unocss/vite"

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react(), unocss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    proxy: {
      "/.well-known": "http://localhost:8444",
      "/api": "http://localhost:8444"
    }
  }
})
