# Storage Service Interface

## API Endpoints

### Profile Operations

1. **Create Profile**

   - `POST /api/v1/profiles`
   - Creates a new profile
   - Request body:
     ```json
     {
       "email": "string",
       "name": "string",
       "metadata": {
         "preferences": {}
       }
     }
     ```
   - Response: Profile object

2. **Get Profile**

   - `GET /api/v1/profiles/{id}`
   - Retrieves profile details
   - Response: Profile object

3. **Update Profile**

   - `PUT /api/v1/profiles/{id}`
   - Updates profile information
   - Request body:
     ```json
     {
       "name": "string",
       "metadata": {
         "preferences": {}
       }
     }
     ```
   - Response: Updated Profile object

4. **Delete Profile**
   - `DELETE /api/v1/profiles/{id}`
   - Soft deletes a profile
   - Response: Success message

### Address Operations

1. **Create Address**

   - `POST /api/v1/profiles/{id}/addresses`
   - Adds an address to a profile
   - Request body:
     ```json
     {
       "type": "string",
       "street": "string",
       "city": "string",
       "state": "string",
       "country": "string",
       "postal_code": "string",
       "is_default": "boolean"
     }
     ```
   - Response: Address object

2. **List Addresses**

   - `GET /api/v1/profiles/{id}/addresses`
   - Lists all addresses for a profile
   - Response: Array of Address objects

3. **Update Address**

   - `PUT /api/v1/profiles/{id}/addresses/{address_id}`
   - Updates address information
   - Request body:
     ```json
     {
       "street": "string",
       "city": "string",
       "state": "string",
       "country": "string",
       "postal_code": "string",
       "is_default": "boolean"
     }
     ```
   - Response: Updated Address object

4. **Delete Address**
   - `DELETE /api/v1/profiles/{id}/addresses/{address_id}`
   - Removes an address
   - Response: Success message

### Contact Operations

1. **Create Contact**

   - `POST /api/v1/profiles/{id}/contacts`
   - Adds a contact to a profile
   - Request body:
     ```json
     {
       "type": "string",
       "value": "string",
       "is_verified": "boolean"
     }
     ```
   - Response: Contact object

2. **List Contacts**

   - `GET /api/v1/profiles/{id}/contacts`
   - Lists all contacts for a profile
   - Response: Array of Contact objects

3. **Update Contact**

   - `PUT /api/v1/profiles/{id}/contacts/{contact_id}`
   - Updates contact information
   - Request body:
     ```json
     {
       "value": "string",
       "is_verified": "boolean"
     }
     ```
   - Response: Updated Contact object

4. **Delete Contact**
   - `DELETE /api/v1/profiles/{id}/contacts/{contact_id}`
   - Removes a contact
   - Response: Success message

### Health and Metrics

1. **Health Check**

   - `GET /health`
   - Returns service health status
   - Response:
     ```json
     {
       "status": "string",
       "version": "string",
       "uptime": "string"
     }
     ```

2. **Database Metrics**

   - `GET /metrics/database`
   - Returns database performance metrics
   - Response: Database metrics object

3. **API Metrics**
   - `GET /metrics/api`
   - Returns API performance metrics
   - Response: API metrics object

## Service Dependencies

### External Services

1. **PostgreSQL**

   - Purpose: Primary data storage
   - Operations:
     - Data persistence
     - Transaction management
     - Query execution
     - Connection management

2. **Redis**
   - Purpose: Caching and rate limiting
   - Operations:
     - Query result caching
     - Rate limiting
     - Session storage
     - Temporary data storage

### Internal Services

1. **Auth Service**

   - Purpose: Authentication and authorization
   - Operations:
     - Token validation
     - Permission checking
     - User context

2. **Monitoring Service**

   - Purpose: Metrics and monitoring
   - Operations:
     - Metrics collection
     - Performance monitoring
     - Health checks
     - Alerting

3. **Logging Service**
   - Purpose: Centralized logging
   - Operations:
     - Log collection
     - Log aggregation
     - Log analysis

## Message Queue Topics

### Profile Events

1. **Profile Changes**

   - Topic: `profiles.changes`
   - Events:
     - Profile created
     - Profile updated
     - Profile deleted
     - Profile restored

2. **Address Changes**

   - Topic: `profiles.addresses`
   - Events:
     - Address added
     - Address updated
     - Address deleted
     - Default address changed

3. **Contact Changes**
   - Topic: `profiles.contacts`
   - Events:
     - Contact added
     - Contact updated
     - Contact deleted
     - Contact verified

## Response Formats

### Success Response

```json
{
  "status": "success",
  "data": {
    // Response data
  },
  "message": "string"
}
```

### Error Response

```json
{
  "status": "error",
  "error": {
    "code": "string",
    "message": "string",
    "details": ["string"]
  }
}
```

## Rate Limiting

1. **API Limits**

   - Profile creation: 100 requests/minute
   - Profile updates: 200 requests/minute
   - Profile queries: 500 requests/minute
   - Address operations: 200 requests/minute
   - Contact operations: 200 requests/minute

2. **Query Limits**
   - Maximum batch size: 100 records
   - Maximum query depth: 3 levels
   - Maximum result size: 1000 records
   - Maximum query timeout: 30 seconds

## Security Headers

### Required Headers

1. **Authorization**

   - `Authorization: Bearer <token>`
   - JWT token for authentication

2. **Request ID**
   - `X-Request-ID: <uuid>`
   - Unique request identifier

### Optional Headers

1. **Client Info**

   - `X-Client-ID: <string>`
   - Client identifier

2. **Trace ID**
   - `X-Trace-ID: <uuid>`
   - Distributed tracing ID

## CORS Configuration

```json
{
  "allowed_origins": ["https://api.example.com", "https://admin.example.com"],
  "allowed_methods": ["GET", "POST", "PUT", "DELETE"],
  "allowed_headers": ["Authorization", "Content-Type", "X-Request-ID"],
  "max_age": 3600
}
```
