# Authentication Guide

This guide explains how to authenticate with the Profile Service APIs using various methods.

## Authentication Methods

### 1. JWT Authentication

The primary authentication method for the Profile Service APIs is JWT (JSON Web Token) based authentication.

#### Obtaining a Token

To obtain a JWT token, make a POST request to the `/auth/token` endpoint:

```bash
curl -X POST https://auth.profileservice.com/v1/auth/token \
  -H "Content-Type: application/json" \
  -d '{
    "grant_type": "password",
    "username": "user@example.com",
    "password": "your-password"
  }'
```

Response:

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Using the Token

Include the token in the Authorization header of your requests:

```bash
curl -X GET https://api.profileservice.com/v1/profiles \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

#### Refreshing the Token

When the access token expires, use the refresh token to obtain a new one:

```bash
curl -X POST https://auth.profileservice.com/v1/auth/token \
  -H "Content-Type: application/json" \
  -d '{
    "grant_type": "refresh_token",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

### 2. API Key Authentication

For service-to-service communication, use API key authentication.

#### Using API Keys

Include the API key in the X-API-Key header:

```bash
curl -X GET https://api.profileservice.com/v1/health \
  -H "X-API-Key: your-api-key"
```

## Go Client Examples

### 1. Authentication Client

```go
package auth

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type AuthClient struct {
    baseURL    string
    httpClient *http.Client
}

type TokenResponse struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"`
    ExpiresIn    int    `json:"expires_in"`
    RefreshToken string `json:"refresh_token"`
}

func NewAuthClient(baseURL string) *AuthClient {
    return &AuthClient{
        baseURL:    baseURL,
        httpClient: &http.Client{},
    }
}

func (c *AuthClient) GetToken(username, password string) (*TokenResponse, error) {
    payload := map[string]string{
        "grant_type": "password",
        "username":   username,
        "password":   password,
    }

    jsonData, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest("POST", c.baseURL+"/auth/token", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var tokenResp TokenResponse
    if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
        return nil, err
    }

    return &tokenResp, nil
}
```

### 2. API Client with Authentication

```go
package api

import (
    "net/http"
)

type APIClient struct {
    baseURL    string
    httpClient *http.Client
    token      string
}

func NewAPIClient(baseURL, token string) *APIClient {
    return &APIClient{
        baseURL:    baseURL,
        httpClient: &http.Client{},
        token:      token,
    }
}

func (c *APIClient) SetToken(token string) {
    c.token = token
}

func (c *APIClient) doRequest(method, path string, body interface{}) (*http.Response, error) {
    req, err := http.NewRequest(method, c.baseURL+path, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", "Bearer "+c.token)
    req.Header.Set("Content-Type", "application/json")

    return c.httpClient.Do(req)
}
```

## Postman Collection

### 1. Environment Setup

Create a new environment in Postman with the following variables:

- `base_url`: https://api.profileservice.com/v1
- `auth_url`: https://auth.profileservice.com/v1
- `access_token`: (leave empty)
- `refresh_token`: (leave empty)

### 2. Authentication Flow

1. **Get Token**

   - Method: POST
   - URL: {{auth_url}}/auth/token
   - Body:
     ```json
     {
       "grant_type": "password",
       "username": "user@example.com",
       "password": "your-password"
     }
     ```
   - Tests:
     ```javascript
     var jsonData = pm.response.json();
     pm.environment.set("access_token", jsonData.access_token);
     pm.environment.set("refresh_token", jsonData.refresh_token);
     ```

2. **Refresh Token**
   - Method: POST
   - URL: {{auth_url}}/auth/token
   - Body:
     ```json
     {
       "grant_type": "refresh_token",
       "refresh_token": "{{refresh_token}}"
     }
     ```

## Security Best Practices

1. **Token Management**

   - Store tokens securely
   - Implement token refresh logic
   - Handle token expiration gracefully
   - Revoke tokens when no longer needed

2. **API Key Security**

   - Rotate API keys regularly
   - Use different keys for different environments
   - Implement key revocation
   - Monitor key usage

3. **Error Handling**

   - Handle 401 Unauthorized responses
   - Implement retry logic with exponential backoff
   - Log authentication failures
   - Monitor authentication patterns

4. **Rate Limiting**
   - Respect rate limits
   - Implement request queuing
   - Handle 429 Too Many Requests responses
   - Monitor rate limit usage

## Next Steps

1. 🔄 Implement token refresh logic
2. 🔄 Add rate limiting handling
3. 🔄 Create error handling examples
4. 🔄 Add monitoring examples
5. 🔄 Create security testing guide
6. 🔄 Add compliance documentation
