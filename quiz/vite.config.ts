import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5174,
    // The quizzes are the QUIZ*.md files in the repo root, outside this folder.
    fs: { allow: ['..'] },
  },
})
