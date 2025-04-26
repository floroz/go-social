import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite"; // Added import

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    tailwindcss(), // Added plugin
    react(),
  ],
});
