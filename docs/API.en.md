# API Documentation

[![Go](https://img.shields.io/badge/Go-1.25-blue.svg)](https://go.dev/dl/)
[![Gin](https://img.shields.io/badge/Gin-1.10-blue.svg)](https://gin-gonic.com/)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-Ready-green.svg)](https://swagger.io/)

> This document covers the complete API specification, including authentication, data export, Webhook notifications, and other core interfaces.

---

## 📑 Table of Contents

- [Basic Information](#basic-information)
- [Authentication and Authorization](#authentication-and-authorization)
- [Account Management API](#account-management-api)
- [Task Management API](#task-management-api)
- [Redemption Center API](#redemption-center-api)
- [Data Export API](#data-export-api)
- [Webhook Notification API](#webhook-notification-api)
- [Announcement Management API](#announcement-management-api)
- [Common Error Codes](#common-error-codes)
- [Rate Limiting](#rate-limiting)
- [Webhook Signature Verification](#webhook-signature-verification)
- [SDKs and Tools](#sdks-and-tools)

---

## Basic Information

### Base URL

```
http://localhost:8080/api
```

Production example:
```
https://api.your-domain.com/api
```

### Protocols

- **HTTP**: `http://localhost:8080/api`
- **HTTPS**: `https://api.your-domain.com/api` (recommended for production)
- **WebSocket**: `ws://localhost:8080/ws` or `wss://api.your-domain.com/ws`

### Request Format

- **Content-Type**: `application/json`
- **Character Encoding**: UTF-8
- **Request Body**: JSON format

### Response Format

All API responses follow a unified format:

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

| Field | Type | Description |
|-------|------|-------------|
| `code` | integer | Business status code, 200 indicates success |
| `message` | string | Response message |
| `data` | object/array/null | Response data |

---

## Authentication and Authorization

### Authentication Method

The system uses **JWT Bearer Token** authentication:

1. User logs in to obtain `access_token` and `refresh_token`
2. Subsequent requests include `Authorization: Bearer {access_token}` in the Header
3. When `access_token` expires, use `refresh_token` to refresh

### Token Validity

| Token Type | Validity | Description |
|-----------|----------|-------------|
| Access Token | 15 minutes | Used for API request authentication |
| Refresh Token | 7 days | Used to refresh Access Token |

### Request Header Format

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

### Login

Obtain access tokens.

**Request Example**:

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "your_password"
}
```

**Success Response**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900,
    "user": {
      "id": 1,
      "username": "admin",
      "role": "admin",
      "created_at": "2026-01-01T00:00:00Z"
    }
  }
}
```

**Failure Response**:

```json
{
  "code": 401,
  "message": "Invalid username or password",
  "data": null
}
```

---

### Refresh Token

Use Refresh Token to obtain a new Access Token.

**Request Example**:

```http
POST /api/auth/refresh
Authorization: Bearer {refresh_token}
```

**Success Response**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900
  }
}
```

---

### Logout

Logout from the current session.

**Request Example**:

```http
POST /api/auth/logout
Authorization: Bearer {access_token}
```

**Success Response**:

```json
{
  "code": 200,
  "message": "logout success",
  "data": null
}
```

---

### Get Current User Info

**Request Example**:

```http
GET /api/auth/me
Authorization: Bearer {access_token}
```

**Success Response**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "role": "admin",
    "created_at": "2026-01-01T00:00:00Z",
    "last_login": "2026-07-06T10:00:00Z"
  }
}
```

---

## Account Management API

### Get Account List

**Request Example**:

```http
GET /api/accounts?page=1&page_size=20&keyword=
Authorization: Bearer {access_token}
```

**Query Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `page` | integer | No | Page number, default 1 |
| `page_size` | integer | No | Items per page, default 20, max 100 |
| `keyword` | string | No | Search keyword |

**Success Response**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 100,
    "page": 1,
    "page_size": 20,
    "list": [
      {
        "id": 1,
        "username": "account_001",
        "platform": "caiyun",
        "status": "active",
        "last_signin_at": "2026-07-06T08:00:00Z",
        "created_at": "2026-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

### Create Account

**Request Example**:

```http
POST /api/accounts
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "username": "new_account",
  "password": "secure_password",
  "platform": "caiyun",
  "remark": "Optional remark"
}
```

---

### Update Account

**Request Example**:

```http
PUT /api/accounts/{account_id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "remark": "Updated remark",
  "status": "active"
}
```

---

### Delete Account

**Request Example**:

```http
DELETE /api/accounts/{account_id}
Authorization: Bearer {access_token}
```

---

## Task Management API

### Get Task List

**Request Example**:

```http
GET /api/tasks?page=1&page_size=20&status=
Authorization: Bearer {access_token}
```

**Query Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `page` | integer | No | Page number, default 1 |
| `page_size` | integer | No | Items per page, default 20 |
| `status` | string | No | Task status: pending/running/success/failed |
| `type` | string | No | Task type |

---

### Create Task

**Request Example**:

```http
POST /api/tasks
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Daily Sign-in",
  "type": "signin",
  "cron": "0 8 * * *",
  "account_ids": [1, 2, 3]
}
```

---

### Get Task Detail

**Request Example**:

```http
GET /api/tasks/{task_id}
Authorization: Bearer {access_token}
```

---

### Delete Task

**Request Example**:

```http
DELETE /api/tasks/{task_id}
Authorization: Bearer {access_token}
```

---

## Redemption Center API

### Get Redemption Records

**Request Example**:

```http
GET /api/exchanges?page=1&page_size=20&status=
Authorization: Bearer {access_token}
```

---

### Create Redemption Task

**Request Example**:

```http
POST /api/exchanges
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "account_id": 1,
  "product_id": "product_001",
  "quantity": 1
}
```

---

### Get Product List

**Request Example**:

```http
GET /api/products?page=1&page_size=20
Authorization: Bearer {access_token}
```

---

## Data Export API

### Create Export Task

Create CSV/JSON format data export tasks.

**Request Example**:

```http
POST /api/export/tasks
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "format": "csv",
  "type": "exchange",
  "start_time": "2026-01-01T00:00:00Z",
  "end_time": "2026-01-31T23:59:59Z",
  "account_ids": [1, 2, 3]
}
```

**Request Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `format` | string | Yes | Export format: `csv` or `json` |
| `type` | string | Yes | Export type: `exchange`, `task`, `account` |
| `start_time` | string | No | Start time (ISO 8601) |
| `end_time` | string | No | End time (ISO 8601) |
| `account_ids` | array | No | Account ID list |

**Success Response**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "task_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "pending",
    "format": "csv",
    "type": "exchange",
    "created_at": "2026-07-06T10:00:00Z"
  }
}
```

---

### Get Export History

**Request Example**:

```http
GET /api/export/history?page=1&page_size=20
Authorization: Bearer {access_token}
```

**Query Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `page` | integer | No | Page number, default 1 |
| `page_size` | integer | No | Items per page, default 20 |
| `format` | string | No | Filter by format |
| `status` | string | No | Filter by status |

**Success Response**:

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
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "format": "csv",
        "type": "exchange",
        "status": "completed",
        "file_url": "/api/export/download/550e8400-e29b-41d4-a716-446655440000",
        "file_size": 1024000,
        "created_at": "2026-07-06T10:00:00Z",
        "completed_at": "2026-07-06T10:00:05Z"
      }
    ]
  }
}
```

---

### Download Export File

Download completed export files.

**Request Example**:

```http
GET /api/export/download/{task_id}
Authorization: Bearer {access_token}
```

**Response**:

- `Content-Type: text/csv` (CSV format)
- `Content-Type: application/json` (JSON format)
- `Content-Disposition: attachment; filename="export_20250101.csv"`

---

### Get Export Task Status

**Request Example**:

```http
GET /api/export/tasks/{task_id}
Authorization: Bearer {access_token}
```

**Success Response**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "processing",
    "progress": 50,
    "created_at": "2026-07-06T10:00:00Z"
  }
}
```

**Task Status Description**:

| Status | Description |
|--------|-------------|
| `pending` | Task queued |
| `processing` | Processing |
| `completed` | Completed, ready for download |
| `failed` | Processing failed |

---

## Webhook Notification API

### Create Webhook Endpoint

Create a Webhook receiver endpoint.

**Request Example**:

```http
POST /api/webhooks
Authorization: Bearer {access_token}
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

**Request Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | Yes | Webhook name |
| `url` | string | Yes | Receiver URL (HTTPS recommended) |
| `events` | array | Yes | List of subscribed event types |
| `is_active` | boolean | No | Whether enabled, default true |
| `headers` | object | No | Custom request headers |

**Success Response**:

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
    "secret": "whsec_xxxxxxxxxxxxxxxxxxxx",
    "created_at": "2026-07-06T10:00:00Z"
  }
}
```

> **Security Note**: After creating a Webhook, keep the `secret` safe for signature verification. This field is only returned during creation.

---

### Get Webhook List

**Request Example**:

```http
GET /api/webhooks
Authorization: Bearer {access_token}
```

---

### Update Webhook

**Request Example**:

```http
PUT /api/webhooks/{webhook_id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Updated Name",
  "is_active": false
}
```

---

### Delete Webhook

**Request Example**:

```http
DELETE /api/webhooks/{webhook_id}
Authorization: Bearer {access_token}
```

---

### Get Delivery Logs

View Webhook delivery history.

**Request Example**:

```http
GET /api/webhooks/{webhook_id}/deliveries?page=1&page_size=20
Authorization: Bearer {access_token}
```

**Query Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `page` | integer | No | Page number, default 1 |
| `page_size` | integer | No | Items per page, default 20 |
| `status` | string | No | Delivery status: success/failed/pending |

**Success Response**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 100,
    "page": 1,
    "page_size": 20,
    "list": [
      {
        "id": 1,
        "webhook_id": 1,
        "event": "exchange.completed",
        "status": "success",
        "http_status": 200,
        "response_body": "OK",
        "attempts": 1,
        "error_message": null,
        "created_at": "2026-07-06T10:00:00Z"
      }
    ]
  }
}
```

---

### Test Webhook

Manually trigger a Webhook test delivery.

**Request Example**:

```http
POST /api/webhooks/{webhook_id}/test
Authorization: Bearer {access_token}
```

---

## Announcement Management API

### Get Announcement List

**Request Example**:

```http
GET /api/announcements?page=1&page_size=20
Authorization: Bearer {access_token}
```

---

### Create Announcement

**Request Example**:

```http
POST /api/announcements
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "title": "System Maintenance Notice",
  "content": "The system will undergo maintenance tonight...",
  "is_pinned": true,
  "expire_at": "2026-12-31T23:59:59Z"
}
```

---

### Delete Announcement

**Request Example**:

```http
DELETE /api/announcements/{announcement_id}
Authorization: Bearer {access_token}
```

---

## Common Error Codes

| Code | Description | Recommended Action |
|------|-------------|-------------------|
| 400 | Bad Request - Invalid parameters | Check request parameter format and required fields |
| 401 | Unauthorized - Missing or expired token | Re-login to obtain new Token |
| 403 | Forbidden - Insufficient permissions | Check user role and permissions |
| 404 | Not Found - Resource does not exist | Check request path and resource ID |
| 409 | Conflict - Resource conflict | Check duplicate creation or status conflict |
| 422 | Business logic validation failed | Correct business parameters based on error message |
| 429 | Too Many Requests - Rate limit exceeded | Reduce request frequency, follow rate limits |
| 500 | Internal Server Error | Contact administrator to check logs |

**Common Error Response Format**：

```json
{
  "code": 400,
  "message": "Parameter error: missing required field 'username'",
  "data": null
}
```

---

## Rate Limiting

The system implements rate limiting on API requests to prevent abuse and brute force attacks.

| Interface Type | Limit | Description |
|---------------|-------|-------------|
| Login endpoint | 5 requests/minute/IP | Prevent brute force attacks |
| Other endpoints | 100 requests/minute/user | Regular API calls |
| Export endpoints | 10 requests/hour/user | Export task limits |
| Webhook delivery | 1000 requests/minute | System-level limit |

**Rate Limit Exceeded Response**：

```json
{
  "code": 429,
  "message": "Too many requests, please try again later",
  "data": {
    "retry_after": 60
  }
}
```

**Response Headers**：

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1699072800
Retry-After: 60
```

---

## Webhook Signature Verification

Each Webhook request carries a signature in the Header for request authenticity verification.

### Request Headers

```http
X-Caiyun-Signature: sha256=<signature>
X-Caiyun-Timestamp: 1699072800
```

### Verification Algorithm

Signature generation: `HMAC-SHA256(secret, timestamp + "\n" + body)`

**Python Example**：

```python
import hmac
import hashlib

def verify_signature(secret, timestamp, body, signature):
    """Verify Webhook signature"""
    expected = hmac.new(
        secret.encode(),
        f"{timestamp}\n{body}".encode(),
        hashlib.sha256
    ).hexdigest()
    return hmac.compare_digest(expected, signature)
```

**Node.js Example**：

```javascript
const crypto = require('crypto');

function verifySignature(secret, timestamp, body, signature) {
  const expected = crypto
    .createHmac('sha256', secret)
    .update(`${timestamp}\n${body}`)
    .digest('hex');
  return crypto.timingSafeEqual(Buffer.from(expected), Buffer.from(signature));
}
```

**Go Example**：

```go
import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "strings"
)

func verifySignature(secret, timestamp, body, signature string) bool {
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write([]byte(timestamp + "\n" + body))
    expected := hex.EncodeToString(mac.Sum(nil))
    return hmac.Equal([]byte(expected), []byte(signature))
}
```

> **Security Notes**：
> 1. Always verify `X-Caiyun-Timestamp`, reject requests older than 5 minutes
> 2. Use `timingSafeEqual` or equivalent functions to prevent timing attacks
> 3. Only use Webhooks in HTTPS environments

---

## Webhook Event Types

| Event Type | Trigger | Data Example |
|-----------|---------|-------------|
| `exchange.completed` | Exchange task completed | `{"exchange_id": 1, "account_id": 1, "status": "success"}` |
| `exchange.failed` | Exchange task failed | `{"exchange_id": 1, "account_id": 1, "error": "Network timeout"}` |
| `exchange.created` | Exchange task created | `{"exchange_id": 1, "account_id": 1, "product_id": "p1"}` |
| `task.completed` | Task execution completed | `{"task_id": 1, "type": "signin", "result": "success"}` |
| `task.failed` | Task execution failed | `{"task_id": 1, "type": "signin", "error": "Account expired"}` |
| `account.created` | Account created | `{"account_id": 1, "username": "new_account"}` |
| `account.updated` | Account updated | `{"account_id": 1, "field": "status", "old": "active", "new": "inactive"}` |
| `account.deleted` | Account deleted | `{"account_id": 1, "username": "deleted_account"}` |

---

## SDKs and Tools

### cURL Examples

```bash
# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your_password"}'

# Get account list
curl -X GET http://localhost:8080/api/accounts \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Create export task
curl -X POST http://localhost:8080/api/export/tasks \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"format":"csv","type":"exchange"}'
```

### Using Postman/Insomnia

Import the following environment variables:

```json
{
  "base_url": "http://localhost:8080/api",
  "access_token": "YOUR_ACCESS_TOKEN",
  "refresh_token": "YOUR_REFRESH_TOKEN"
}
```

---

## 📞 Technical Support

- **API Issues**: [GitHub Issues](https://github.com/xiaohai09/mobile-cloud-disk-management/issues)
- **Security Vulnerabilities**: [Security Advisories](https://github.com/xiaohai09/mobile-cloud-disk-management/security/advisories)

---

*Last updated: 2026-07-07*
