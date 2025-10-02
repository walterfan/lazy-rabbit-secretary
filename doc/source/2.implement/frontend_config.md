# Environment Configuration

This project uses environment variables to configure API endpoints and other settings. This eliminates hardcoded URLs and makes the application more flexible across different environments.

## Environment Files

- `.env.example` - Template file with all available environment variables
- `.env` - Local development configuration (not committed to git)
- `.env.production` - Production configuration (if needed)

## Available Environment Variables

### API Configuration

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `VITE_API_BASE_URL` | Complete API base URL | `https://localhost:9090` | `https://api.example.com` |
| `VITE_API_PROTOCOL` | API protocol (http/https) | `https` | `https` |
| `VITE_API_HOST` | API hostname | `localhost` | `api.example.com` |
| `VITE_API_PORT` | API port | `9090` | `443` |

### Usage Priority

1. **`VITE_API_BASE_URL`** - If set, this takes precedence over individual components
2. **Component-based** - If `VITE_API_BASE_URL` is not set, the system builds the URL from:
   - `VITE_API_PROTOCOL` + `VITE_API_HOST` + `VITE_API_PORT`

## Configuration Examples

### Development (Local)
```bash
# .env
VITE_API_BASE_URL=https://localhost:9090
VITE_API_PROTOCOL=https
VITE_API_HOST=localhost
VITE_API_PORT=9090
```

### Production
```bash
# .env.production
VITE_API_BASE_URL=https://api.yourdomain.com
VITE_API_PROTOCOL=https
VITE_API_HOST=api.yourdomain.com
VITE_API_PORT=443
```

### Docker/Container
```bash
# .env.docker
VITE_API_BASE_URL=https://backend:9090
VITE_API_PROTOCOL=https
VITE_API_HOST=backend
VITE_API_PORT=9090
```

## How It Works

### Frontend Configuration

The frontend uses a centralized API configuration utility (`src/utils/apiConfig.ts`) that:

1. Reads environment variables
2. Builds the appropriate API base URL
3. Provides helper functions for API calls
4. Logs configuration in development mode

### Vite Proxy Configuration

The Vite development server proxy configuration (`vite.config.ts`) automatically uses the same environment variables to:

1. Configure API proxy targets
2. Handle CORS and SSL certificates
3. Route requests to the correct backend

### Store Integration

All Pinia stores (`wikiStore.ts`, `postStore.ts`, etc.) now use the centralized `getApiUrl()` function instead of hardcoded URLs.

## Migration from Hardcoded URLs

### Before
```typescript
// Hardcoded URL
const response = await fetch('https://localhost:9090/api/v1/posts')
```

### After
```typescript
// Environment-based URL
import { getApiUrl } from '@/utils/apiConfig'
const response = await fetch(getApiUrl('/api/v1/posts'))
```

## Benefits

1. **Environment Flexibility** - Easy switching between dev/staging/production
2. **No Hardcoded URLs** - All URLs are configurable
3. **Centralized Configuration** - Single source of truth for API settings
4. **Docker Ready** - Easy containerization with environment variables
5. **Team Collaboration** - Each developer can have their own `.env` file

## Troubleshooting

### Common Issues

1. **API calls failing** - Check if `VITE_API_BASE_URL` is correctly set
2. **Proxy not working** - Verify `VITE_API_HOST` and `VITE_API_PORT` in vite.config.ts
3. **SSL errors** - Ensure `VITE_API_PROTOCOL=https` for HTTPS endpoints

### Debug Mode

In development, the API configuration is logged to the console:
```javascript
API Configuration: {
  baseUrl: "https://localhost:9090",
  environment: "development",
  isDevelopment: true,
  isProduction: false
}
```
