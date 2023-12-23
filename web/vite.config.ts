import path from "path"
import react from "@vitejs/plugin-react"
import { defineConfig } from "vite"

// https://vitejs.dev/config/
export default defineConfig(({ command, mode }) => {
  let proxy = {}
  if (mode === "development") {
    proxy = {
      "/api": {
        target: "http://localhost:9090",
        secure: false, // Allow self-signed certificates
      },
      "/ws": {
        target: "http://localhost:9090",
        ws: true,
        secure: false, // Allow self-signed certificates
      },
    }
  }

  return {
    server: {
      proxy: proxy,
    },
    plugins: [react()],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },
  }
})
