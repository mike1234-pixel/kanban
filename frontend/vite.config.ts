import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/boards': 'http://localhost:8080',
      '/columns': 'http://localhost:8080',
      '/cards': 'http://localhost:8080',
      '/login': 'http://localhost:8080',
      '/logout': 'http://localhost:8080',
      '/me': 'http://localhost:8080',
    },
  },
})
