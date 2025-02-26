# API Key Authentication for ServiceRadar Web Interface

This document explains how the API key authentication works in the ServiceRadar web interface.

## Overview

The ServiceRadar web application has been updated to support API key authentication for all API requests to the backend. This follows the 12-factor application methodology by loading the API key from environment variables.

## Setup Instructions

### 1. Environment Configuration

Create a `.env` file in the root directory of the project with the following content:

```
VITE_API_KEY=your-api-key-here
```

Make sure this matches the API key configured on your server.

> **Note:** The `.env` file should never be committed to version control. A `.env.example` file is provided as a template.

### 2. For Production Deployment

In a production environment, you should set the environment variable according to your deployment platform:

- **Docker:** Add it to your Docker Compose file or Docker run command:
  ```
  docker run -e VITE_API_KEY=your-api-key-here ...
  ```

- **Kubernetes:** Add it to your deployment configuration:
  ```yaml
  env:
    - name: VITE_API_KEY
      value: your-api-key-here
  # Or better, use a secret:
  env:
    - name: VITE_API_KEY
      valueFrom:
        secretKeyRef:
          name: serviceradar-secrets
          key: api-key
  ```

- **Traditional hosting:** Set the environment variable in your hosting environment or through your CI/CD pipeline.

### 3. Development Environment

For local development, simply create a `.env` file as described above. The Vite development server will automatically load the environment variables.

## How It Works

1. The web interface reads the API key from the environment variable `VITE_API_KEY`.
2. All API requests are processed through a centralized API service that automatically appends the API key as an `X-API-Key` header.
3. The backend validates this header against its configured API key.
4. Static assets (JavaScript, CSS, images) are exempt from API key validation.

## API Service

The API service is implemented in `src/services/api.js` and provides methods for making authenticated API requests:

```javascript
// Example usage:
import { get, post } from '../services/api';

// GET request
const data = await get('/api/nodes');

// POST request
const response = await post('/api/some-endpoint', { key: 'value' });
```

## Security Considerations

- The API key is included in the compiled JavaScript bundle and can be viewed by users with access to your application. This provides a basic level of authentication but should not be considered fully secure.
- For more sensitive operations, consider implementing a more robust authentication system (OAuth, JWT, etc.).
- Always use HTTPS in production to prevent the API key from being intercepted.