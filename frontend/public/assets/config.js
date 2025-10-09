// Runtime configuration for the application
// This file can be modified after build to change API endpoints

window.__API_CONFIG__ = {
  baseUrl: 'https://www.lazy-rabbit-studio.top', // Change this to your desired API URL
  // You can also set individual components:
  // protocol: 'https',
  // host: 'your-server.com',
  // port: '443'
};

// Optional: Log the configuration for debugging
console.log('Runtime API Configuration loaded:', window.__API_CONFIG__);
