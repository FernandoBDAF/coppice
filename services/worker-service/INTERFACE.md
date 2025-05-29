# Worker Service Interface

## API Endpoints

### Job Management

1. **Create Job**

   - `POST /api/v1/jobs`
   - Creates a new job
   - Request body:
     ```json
     {
       "type": "string",
       "payload": "object",
       "priority": "integer",
       "schedule": "string (optional)"
     }
     ```
   - Response: Job object

2. **Get Job**

   - `GET /api/v1/jobs/{job_id}`
   - Retrieves job status and details
   - Response: Job object

3. **List Jobs**

   - `GET /api/v1/jobs`
   - Lists all jobs with optional filters
   - Query parameters:
     - `status`: Job status filter
     - `type`: Job type filter
     - `limit`: Maximum number of jobs
     - `offset`: Pagination offset
   - Response: Array of Job objects

4. **Cancel Job**
   - `DELETE /api/v1/jobs/{job_id}`
   - Cancels a running job
   - Response: Success message

### Task Operations

1. **Create Task**

   - `POST /api/v1/tasks`
   - Creates a new scheduled task
   - Request body:
     ```json
     {
       "type": "string",
       "schedule": "string",
       "payload": "object",
       "dependencies": ["string"]
     }
     ```
   - Response: Task object

2. **List Tasks**

   - `GET /api/v1/tasks`
   - Lists all tasks with optional filters
   - Query parameters:
     - `status`: Task status filter
     - `type`: Task type filter
     - `limit`: Maximum number of tasks
     - `offset`: Pagination offset
   - Response: Array of Task objects

3. **Update Task**

   - `PUT /api/v1/tasks/{task_id}`
   - Updates task schedule or dependencies
   - Request body:
     ```json
     {
       "schedule": "string",
       "dependencies": ["string"]
     }
     ```
   - Response: Updated Task object

4. **Delete Task**
   - `DELETE /api/v1/tasks/{task_id}`
   - Deletes a scheduled task
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

2. **Worker Metrics**

   - `GET /metrics/workers`
   - Returns worker performance metrics
   - Response: Worker metrics object

3. **Job Metrics**
   - `GET /metrics/jobs`
   - Returns job processing metrics
   - Response: Job metrics object

## Service Dependencies

### External Services

1. **RabbitMQ**

   - Purpose: Message queue for job processing
   - Operations:
     - Job queue management
     - Message publishing
     - Message consumption
     - Queue monitoring

2. **Redis**
   - Purpose: Job state and task scheduling
   - Operations:
     - Job state storage
     - Task scheduling
     - Rate limiting
     - Cache management

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

### Job Events

1. **Job Creation**

   - Topic: `jobs.created`
   - Events:
     - Job created
     - Job scheduled
     - Job queued

2. **Job Processing**

   - Topic: `jobs.processing`
   - Events:
     - Job started
     - Job completed
     - Job failed
     - Job retried

3. **Task Events**
   - Topic: `tasks.scheduled`
   - Events:
     - Task created
     - Task scheduled
     - Task triggered
     - Task completed

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

   - Job creation: 100 requests/minute
   - Task creation: 50 requests/minute
   - Job listing: 200 requests/minute
   - Task listing: 100 requests/minute

2. **Job Limits**
   - Maximum job size: 10MB
   - Maximum retry count: 3
   - Maximum job duration: 1 hour
   - Maximum concurrent jobs: 1000

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
