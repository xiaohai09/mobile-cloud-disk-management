# API Documentation

## Basic Information

- **Base URL**: `http://localhost:8080/api`
- **Authentication**: JWT Bearer Token
- **Request Format**: `application/json`
- **Response Format**: `application/json`

## Authentication

### Login

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password"
}
```

**Response Example**：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900,
    "user": {
      "id": 1,
      "username": "admin",
      "role": "admin"
    }
  }
}
```

### Refresh Token

```http
POST /api/auth/refresh
Authorization: Bearer <refresh_token>
```

### Logout

```http
POST /api/auth/logout
Authorization: Bearer <access_token>
```

## Data Export API

### Create Export Task

```http
POST /api/export/tasks
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "format": "csv",        // csv | json
  "type": "exchange",     // exchange | task | account
  "start_time": "2026-01-01T00:00:00Z",
  "end_time": "2026-01-31T23:59:59Z",
  "account_ids": [1, 2, 3]
}
```

**Response Example**：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "task_id": "uuid",
    "status": "pending",
    "format": "csv",
    "created_at": "2026-07-06T10:00:00Z"
  }
}
```

### Get Export History

```http
GET /api/export/history?page=1&page_size=20
Authorization: Bearer <access_token>
```

**Response Example**：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 50,
    "page": 1,
    "page_size": 20,
    "list": [
      {
        "id": 1,
        "format": "csv",
        "type": "exchange",
        "status": "completed",
        "file_url": "/api/export/download/1",
        "created_at": "2026-07-06T10:00:00Z"
      }
    ]
  }
}
```

### Download Export File

```http
GET /api/export/download/{task_id}
Authorization: Bearer <access_token>
```

## Webhook API

### Create Webhook Endpoint

```http
POST /api/webhooks
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "Transaction Completion Notification",
  "url": "https://example.com/webhook",
  "events": ["exchange.completed", "exchange.failed"],
  "is_active": true,
  "headers": {
    "X-Custom-Header": "value"
  }
}
```

**Response Example**：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "Transaction Completion Notification",
    "url": "https://example.com/webhook",
    "events": ["exchange.completed", "exchange.failed"],
    "is_active": true,
    "secret": "whsec_xxxxx",
    "created_at": "2026-07-06T10:00:00Z"
  }
}
```

### Get Webhook List

```http
GET /api/webhooks
Authorization: Bearer <access_token>
```

### Update Webhook

```http
PUT /api/webhooks/{webhook_id}
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "Updated Name",
  "is_active": false
}
```

### Delete Webhook

```http
DELETE /api/webhooks/{webhook_id}
Authorization: Bearer <access_token>
```

### Get Delivery Logs

```http
GET /api/webhooks/{webhook_id}/deliveries?page=1&page_size=20
Authorization: Bearer <access_token>
```

**Response Example**：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 100,
    "list": [
      {
        "id": 1,
        "event": "exchange.completed",
        "status": "success",
        "http_status": 200,
        "attempts": 1,
        "created_at": "2026-07-06T10:00:00Z"
      }
    ]
  }
}
```

## Common Error Codes

| Code | Description |
|------|-------------|
| 400 | Bad Request - Invalid parameters |
| 401 | Unauthorized - Missing or expired token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource does not exist |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error |

**General Response Format**：

```json
{
  "code": 400,
  "message": "Parameter error: missing required field 'username'",
  "data": null
}
```

## Webhook Event Types

| Event Type | Trigger |
|------------|---------|
| `exchange.completed` | Exchange task completed |
| `exchange.failed` | Exchange task failed |
| `exchange.created` | Exchange task created |
| `task.completed` | Task execution completed |
| `task.failed` | Task execution failed |
| `account.created` | Account created |
| `account.updated` | Account updated |
| `account.deleted` | Account deleted |

## Webhook Signature Verification

Each webhook request includes a signature in the Header:

```
X-Caiyun-Signature: sha256=<signature>
X-Caiyun-Timestamp: 1699072800
```

Verification method:

```python
import hmac
import hashlib

signature = hmac.new(
    webhook_secret.encode(),
    f"{timestamp}\n{body}".encode(),
    hashlib.sha256
).hexdigest()
```

## Rate Limiting

- Login endpoint: `5 requests/minute/IP`
- Other endpoints: `100 requests/minute/user`
- Export endpoints: `10 requests/hour/user`

Exceeding limits returns `429 Too Many Requests`.
