import { defineConfig } from "vite"
import react from "@vitejs/plugin-react"
import { viteSingleFile } from "vite-plugin-singlefile"

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react(), viteSingleFile()],
  server: {
    port: 8080,
  },
  // See https://github.com/richardtallent/vite-plugin-singlefile#how-do-i-use-it
  build: {
    assetsInlineLimit: 100000000,
    chunkSizeWarningLimit: 100000000,
    rollupOptions: {
			output: {
				manualChunks: () => "everything.js",
			},
		},
  }
})
