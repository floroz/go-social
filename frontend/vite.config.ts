import { defineConfig } from "vite";
import path from "path"; // Added import
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite"; // Added import

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    tailwindcss(), // Added plugin
    react(),
  ],
  resolve: {
    // Added resolve configuration
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
});
