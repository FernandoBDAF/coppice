# Profile API Go Client Examples

This document provides Go client examples for interacting with the Profile API.

## Client Setup

### Basic Client

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/your-org/profile-service/client"
)

func main() {
    cfg := client.Config{
        BaseURL:    "https://api.profileservice.com/v1",
        APIKey:     "your-api-key",
        Timeout:    30 * time.Second,
        MaxRetries: 3,
    }

    c, err := client.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()
}
```

### Client with Custom Transport

```go
package main

import (
    "crypto/tls"
    "net/http"
    "time"

    "github.com/your-org/profile-service/client"
)

func main() {
    transport := &http.Transport{
        TLSClientConfig: &tls.Config{
            MinVersion: tls.VersionTLS12,
        },
        MaxIdleConns:        100,
        IdleConnTimeout:     90 * time.Second,
        DisableCompression:  false,
    }

    cfg := client.Config{
        BaseURL:    "https://api.profileservice.com/v1",
        APIKey:     "your-api-key",
        Timeout:    30 * time.Second,
        Transport:  transport,
        MaxRetries: 3,
    }

    c, err := client.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()
}
```

## Profile Operations

### List Profiles

```go
func listProfiles(ctx context.Context, c *client.Client) error {
    params := client.ListProfilesParams{
        Page:  1,
        Limit: 10,
        Sort:  "name",
        Order: "asc",
    }

    profiles, err := c.ListProfiles(ctx, params)
    if err != nil {
        return fmt.Errorf("failed to list profiles: %w", err)
    }

    for _, profile := range profiles {
        log.Printf("Profile: %+v", profile)
    }
    return nil
}
```

### Get Profile by ID

```go
func getProfile(ctx context.Context, c *client.Client, id string) error {
    profile, err := c.GetProfile(ctx, id)
    if err != nil {
        return fmt.Errorf("failed to get profile: %w", err)
    }

    log.Printf("Profile: %+v", profile)
    return nil
}
```

### Create Profile

```go
func createProfile(ctx context.Context, c *client.Client) error {
    profile := client.Profile{
        Name:      "John Doe",
        Email:     "john.doe@example.com",
        Bio:       "Software Engineer",
        ImageURLs: []string{"https://example.com/image1.jpg"},
    }

    created, err := c.CreateProfile(ctx, profile)
    if err != nil {
        return fmt.Errorf("failed to create profile: %w", err)
    }

    log.Printf("Created profile: %+v", created)
    return nil
}
```

### Update Profile

```go
func updateProfile(ctx context.Context, c *client.Client, id string) error {
    update := client.ProfileUpdate{
        Name: "John Doe Updated",
        Bio:  "Senior Software Engineer",
    }

    updated, err := c.UpdateProfile(ctx, id, update)
    if err != nil {
        return fmt.Errorf("failed to update profile: %w", err)
    }

    log.Printf("Updated profile: %+v", updated)
    return nil
}
```

### Delete Profile

```go
func deleteProfile(ctx context.Context, c *client.Client, id string) error {
    err := c.DeleteProfile(ctx, id)
    if err != nil {
        return fmt.Errorf("failed to delete profile: %w", err)
    }

    log.Printf("Successfully deleted profile %s", id)
    return nil
}
```

## Error Handling

### Custom Error Types

```go
type APIError struct {
    Code    string            `json:"code"`
    Message string            `json:"message"`
    Details map[string]string `json:"details"`
}

func handleError(err error) {
    var apiErr *APIError
    if errors.As(err, &apiErr) {
        switch apiErr.Code {
        case "validation_error":
            log.Printf("Validation error: %s", apiErr.Message)
            for field, detail := range apiErr.Details {
                log.Printf("  %s: %s", field, detail)
            }
        case "not_found":
            log.Printf("Resource not found: %s", apiErr.Message)
        case "rate_limited":
            log.Printf("Rate limited: %s", apiErr.Message)
        default:
            log.Printf("API error: %s", apiErr.Message)
        }
    } else {
        log.Printf("Unexpected error: %v", err)
    }
}
```

### Retry Logic

```go
func withRetry(ctx context.Context, c *client.Client, id string) error {
    var profile *client.Profile
    var err error

    for i := 0; i < 3; i++ {
        profile, err = c.GetProfile(ctx, id)
        if err == nil {
            break
        }

        if errors.Is(err, client.ErrRateLimited) {
            time.Sleep(time.Duration(i+1) * time.Second)
            continue
        }

        return err
    }

    if err != nil {
        return fmt.Errorf("failed after retries: %w", err)
    }

    log.Printf("Profile: %+v", profile)
    return nil
}
```

## Context Usage

### With Timeout

```go
func withTimeout(ctx context.Context, c *client.Client) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    profiles, err := c.ListProfiles(ctx, client.ListProfilesParams{})
    if err != nil {
        return fmt.Errorf("failed to list profiles: %w", err)
    }

    log.Printf("Found %d profiles", len(profiles))
    return nil
}
```

### With Request ID

```go
func withRequestID(ctx context.Context, c *client.Client) error {
    requestID := uuid.New().String()
    ctx = context.WithValue(ctx, client.RequestIDKey, requestID)

    profile, err := c.GetProfile(ctx, "123")
    if err != nil {
        return fmt.Errorf("failed to get profile: %w", err)
    }

    log.Printf("Profile: %+v", profile)
    return nil
}
```

## Best Practices

1. **Error Handling**

   - Use custom error types
   - Implement retry logic
   - Handle rate limiting
   - Log errors appropriately

2. **Context Usage**

   - Always pass context
   - Set appropriate timeouts
   - Include request IDs
   - Handle cancellation

3. **Resource Management**

   - Close client connections
   - Use connection pooling
   - Implement circuit breakers
   - Monitor resource usage

4. **Security**

   - Use TLS
   - Validate inputs
   - Handle sensitive data
   - Implement authentication

5. **Performance**
   - Use connection pooling
   - Implement caching
   - Handle rate limiting
   - Monitor metrics
