import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
  server: {
    host: "0.0.0.0",
    cors: true,
    // Prefix dev asset URLs (including url() in CSS) with the dev server
    // origin so they resolve to Vite, not the Go app serving the HTML.
    origin: "http://localhost:5173",
  },
  plugins: [tailwindcss()],
  build: {
    manifest: true,
  },
});

