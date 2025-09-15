import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import { fileURLToPath, URL } from 'url';
import fs from 'fs';

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    https: {
      key: fs.readFileSync('../backend/certs/private.pem'),
      cert: fs.readFileSync('../backend/certs/certificate.pem')
    },
    host: 'localhost',
    port: 5173,
    proxy: {
      '/api': {
        target: 'https://localhost:9090',
        changeOrigin: true,
        secure: false, // Allow self-signed certificates
        rewrite: (path) => path
      },
      '/news': {
        target: 'https://localhost:9090',
        changeOrigin: true,
        secure: false, // Allow self-signed certificates
        rewrite: (path) => path
      }
    }
  }
});