import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
  server: {
    host: "0.0.0.0",
    cors: true,
  },
  plugins: [tailwindcss()],
  build: {
    manifest: true,
  },
});

