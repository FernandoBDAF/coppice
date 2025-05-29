# Queue Service Interface

## API Endpoints

### Queue Management

1. **Create Queue**

   - Method: `POST`
   - Path: `/api/v1/queues`
   - Request:
     ```json
     {
       "name": "string",
       "type": "direct|fanout|topic",
       "options": {
         "durable": "boolean",
         "auto_delete": "boolean",
         "arguments": {}
       }
     }
     ```
   - Response: `201 Created`
   - Auth: Required

2. **Get Queue**

   - Method: `GET`
   - Path: `/api/v1/queues/{queue_name}`
   - Response: `200 OK`
   - Auth: Required

3. **List Queues**

   - Method: `GET`
   - Path: `/api/v1/queues`
   - Query Parameters:
     - `type`: Queue type filter
     - `status`: Queue status filter
   - Response: `200 OK`
   - Auth: Required

4. **Delete Queue**
   - Method: `DELETE`
   - Path: `/api/v1/queues/{queue_name}`
   - Response: `204 No Content`
   - Auth: Required

### Message Operations

1. **Publish Message**

   - Method: `POST`
   - Path: `/api/v1/queues/{queue_name}/messages`
   - Request:
     ```json
     {
       "type": "string",
       "payload": "string",
       "headers": {},
       "options": {
         "persistent": "boolean",
         "priority": "integer"
       }
     }
     ```
   - Response: `201 Created`
   - Auth: Required

2. **Consume Messages**

   - Method: `GET`
   - Path: `/api/v1/queues/{queue_name}/messages`
   - Query Parameters:
     - `limit`: Maximum messages to fetch
     - `timeout`: Consumer timeout
   - Response: `200 OK`
   - Auth: Required

3. **Acknowledge Message**

   - Method: `POST`
   - Path: `/api/v1/messages/{message_id}/ack`
   - Response: `200 OK`
   - Auth: Required

4. **Reject Message**
   - Method: `POST`
   - Path: `/api/v1/messages/{message_id}/reject`
   - Request:
     ```json
     {
       "requeue": "boolean"
     }
     ```
   - Response: `200 OK`
   - Auth: Required

### Event Operations

1. **Publish Event**

   - Method: `POST`
   - Path: `/api/v1/events`
   - Request:
     ```json
     {
       "type": "string",
       "source": "string",
       "data": "string",
       "metadata": {}
     }
     ```
   - Response: `201 Created`
   - Auth: Required

2. **Subscribe to Events**
   - Method: `POST`
   - Path: `/api/v1/events/subscribe`
   - Request:
     ```json
     {
       "types": ["string"],
       "callback_url": "string"
     }
     ```
   - Response: `201 Created`
   - Auth: Required

### Health and Metrics

1. **Health Check**

   - Method: `GET`
   - Path: `/health`
   - Response: `200 OK`
   - Auth: None

2. **Queue Metrics**

   - Method: `GET`
   - Path: `/metrics/queues`
   - Response: `200 OK`
   - Auth: Required

3. **Message Metrics**
   - Method: `GET`
   - Path: `/metrics/messages`
   - Response: `200 OK`
   - Auth: Required

## Service Dependencies

### External Services

1. **RabbitMQ**

   - Purpose: Message broker
   - Operations:
     - Queue management
     - Message routing
     - Message persistence
     - Dead letter handling

2. **Redis**
   - Purpose: Cache and rate limiting
   - Operations:
     - Rate limiting
     - Message deduplication
     - Temporary storage
     - Lock management

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
     - Queue metrics
     - Message metrics
     - Health checks
     - Alerting

3. **Logging Service**
   - Purpose: Centralized logging
   - Operations:
     - Log collection
     - Log aggregation
     - Log analysis

## Message Queue Topics

1. **Queue Events**

   - `queue.created`
   - `queue.deleted`
   - `queue.updated`
   - `queue.error`

2. **Message Events**

   - `message.published`
   - `message.consumed`
   - `message.acknowledged`
   - `message.rejected`
   - `message.error`

3. **System Events**
   - `system.health`
   - `system.metrics`
   - `system.error`
   - `system.alert`

## Response Formats

1. **Success Response**

   ```json
   {
     "status": "success",
     "data": {},
     "message": "string"
   }
   ```

2. **Error Response**
   ```json
   {
     "status": "error",
     "error": {
       "type": "string",
       "message": "string",
       "details": []
     }
   }
   ```

## Rate Limiting

1. **API Limits**

   - Queue operations: 100 requests/minute
   - Message operations: 1000 requests/minute
   - Event operations: 500 requests/minute
   - Metrics operations: 60 requests/minute

2. **Message Limits**
   - Message size: 1MB
   - Batch size: 100 messages
   - Consumer prefetch: 50 messages
   - Retry attempts: 3

## Security Headers

1. **Required Headers**

   - `Authorization`: Bearer token
   - `X-Request-ID`: Request tracking
   - `Content-Type`: application/json

2. **Optional Headers**
   - `X-User-ID`: User context
   - `X-Service-ID`: Service identification
   - `X-Correlation-ID`: Request correlation

## CORS Configuration

```go
config := cors.Config{
    AllowedOrigins: []string{
        "https://api.example.com",
        "https://admin.example.com"
    },
    AllowedMethods: []string{
        "GET",
        "POST",
        "PUT",
        "DELETE",
        "OPTIONS"
    },
    AllowedHeaders: []string{
        "Authorization",
        "Content-Type",
        "X-Request-ID",
        "X-User-ID",
        "X-Service-ID",
        "X-Correlation-ID"
    },
    MaxAge: 12 * time.Hour
}
```
