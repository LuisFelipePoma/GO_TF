import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'

// https://vite.dev/config/
export default defineConfig({
  base: '/',
  plugins: [react()],
  preview: {
    port: 4323,
    strictPort: true
  },
  server: {
    port: 4323,
    strictPort: true,
    host: true,
    origin: 'http://0.0.0.0:4323'
  }
})
