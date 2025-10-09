import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import { fileURLToPath, URL } from 'url';
import fs from 'fs';

export default defineConfig(({ mode }) => {
  // Get API configuration from environment variables
  const apiProtocol = process.env.VITE_API_PROTOCOL || 'https';
  const apiHost = process.env.VITE_API_HOST || 'localhost';
  const apiPort = process.env.VITE_API_PORT || '9090';
  const apiTarget = `${apiProtocol}://${apiHost}:${apiPort}`;

  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      }
    },
    build: {
      // Increase chunk size warning limit
      chunkSizeWarningLimit: 1000,
      // CSS optimization
      cssCodeSplit: true,
      // Minification options
      minify: 'terser',
      terserOptions: {
        compress: {
          drop_console: mode === 'production',
          drop_debugger: mode === 'production'
        }
      },
      // Copy public files including config.js
      rollupOptions: {
        input: {
          main: './index.html'
        }
      }
    },
    css: {
      // CSS processing options
      postcss: './postcss.config.cjs'
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
          target: apiTarget,
          changeOrigin: true,
          secure: false, // Allow self-signed certificates
          rewrite: (path) => path
        },
      }
    }
  };
});