import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      "/gameInfo": "http://localhost:8080",
      "/changeDir": "http://localhost:8080",
      "/addSnake": "http://localhost:8080",
      "/createGame": "http://localhost:8080",
      "/quitAll": "http://localhost:8080",
    },
  },
});
