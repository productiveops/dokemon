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
      },
      "/ws": {
        target: "http://localhost:9090",
        ws: true,
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
