# API 文档

[![Go](https://img.shields.io/badge/Go-1.25-blue.svg)](https://go.dev/dl/)
[![Gin](https://img.shields.io/badge/Gin-1.10-blue.svg)](https://gin-gonic.com/)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-Ready-green.svg)](https://swagger.io/)

> 本文档涵盖系统的完整 API 规范，包括认证、数据导出、Webhook 通知等核心接口。

---

## 📑 目录

- [基础信息](#基础信息)
- [认证与授权](#认证与授权)
- [账号管理 API](#账号管理-api)
- [任务管理 API](#任务管理-api)
- [兑换中心 API](#兑换中心-api)
- [数据导出 API](#数据导出-api)
- [Webhook 通知 API](#webhook-通知-api)
- [公告管理 API](#公告管理-api)
- [通用错误码](#通用错误码)
- [速率限制](#速率限制)
- [Webhook 签名验证](#webhook-签名验证)
- [SDK 与工具](#sdk-与工具)

---

## 基础信息

### Base URL

```
http://localhost:8080/api
```

生产环境示例：
```
https://api.your-domain.com/api
```

### 协议

- **HTTP**: `http://localhost:8080/api`
- **HTTPS**: `https://api.your-domain.com/api`（生产环境推荐）
- **WebSocket**: `ws://localhost:8080/ws` 或 `wss://api.your-domain.com/ws`

### 请求格式

- **Content-Type**: `application/json`
- **字符编码**: UTF-8
- **请求体**: JSON 格式

### 响应格式

所有 API 响应遵循统一格式：

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | integer | 业务状态码，200 表示成功 |
| `message` | string | 响应消息 |
| `data` | object/array/null | 响应数据 |

---

## 认证与授权

### 认证方式

系统采用 **JWT Bearer Token** 认证机制：

1. 用户登录获取 `access_token` 和 `refresh_token`
2. 后续请求在 Header 中携带 `Authorization: Bearer {access_token}`
3. `access_token` 过期后使用 `refresh_token` 刷新

### Token 有效期

| Token 类型 | 有效期 | 说明 |
|-----------|--------|------|
| Access Token | 15 分钟 | 用于 API 请求认证 |
| Refresh Token | 7 天 | 用于刷新 Access Token |

### 请求头格式

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

### 登录

获取访问令牌。

**请求示例**：

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "your_password"
}
```

**成功响应**：

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

**失败响应**：

```json
{
  "code": 401,
  "message": "用户名或密码错误",
  "data": null
}
```

---

### 刷新 Token

使用 Refresh Token 获取新的 Access Token。

**请求示例**：

```http
POST /api/auth/refresh
Authorization: Bearer {refresh_token}
```

**成功响应**：

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

### 登出

注销当前会话。

**请求示例**：

```http
POST /api/auth/logout
Authorization: Bearer {access_token}
```

**成功响应**：

```json
{
  "code": 200,
  "message": "logout success",
  "data": null
}
```

---

### 获取当前用户信息

**请求示例**：

```http
GET /api/auth/me
Authorization: Bearer {access_token}
```

**成功响应**：

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

## 账号管理 API

### 获取账号列表

**请求示例**：

```http
GET /api/accounts?page=1&page_size=20&keyword=
Authorization: Bearer {access_token}
```

**查询参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `page` | integer | 否 | 页码，默认 1 |
| `page_size` | integer | 否 | 每页数量，默认 20，最大 100 |
| `keyword` | string | 否 | 搜索关键词 |

**成功响应**：

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

### 创建账号

**请求示例**：

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

### 更新账号

**请求示例**：

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

### 删除账号

**请求示例**：

```http
DELETE /api/accounts/{account_id}
Authorization: Bearer {access_token}
```

---

## 任务管理 API

### 获取任务列表

**请求示例**：

```http
GET /api/tasks?page=1&page_size=20&status=
Authorization: Bearer {access_token}
```

**查询参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `page` | integer | 否 | 页码，默认 1 |
| `page_size` | integer | 否 | 每页数量，默认 20 |
| `status` | string | 否 | 任务状态：pending/running/success/failed |
| `type` | string | 否 | 任务类型 |

---

### 创建任务

**请求示例**：

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

### 获取任务详情

**请求示例**：

```http
GET /api/tasks/{task_id}
Authorization: Bearer {access_token}
```

---

### 删除任务

**请求示例**：

```http
DELETE /api/tasks/{task_id}
Authorization: Bearer {access_token}
```

---

## 兑换中心 API

### 获取兑换记录

**请求示例**：

```http
GET /api/exchanges?page=1&page_size=20&status=
Authorization: Bearer {access_token}
```

---

### 创建兑换任务

**请求示例**：

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

### 获取产品列表

**请求示例**：

```http
GET /api/products?page=1&page_size=20
Authorization: Bearer {access_token}
```

---

## 数据导出 API

### 创建导出任务

创建 CSV/JSON 格式的数据导出任务。

**请求示例**：

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

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `format` | string | 是 | 导出格式：`csv` 或 `json` |
| `type` | string | 是 | 导出类型：`exchange`、`task`、`account` |
| `start_time` | string | 否 | 开始时间（ISO 8601） |
| `end_time` | string | 否 | 结束时间（ISO 8601） |
| `account_ids` | array | 否 | 账号 ID 列表 |

**成功响应**：

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

### 获取导出历史

**请求示例**：

```http
GET /api/export/history?page=1&page_size=20
Authorization: Bearer {access_token}
```

**查询参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `page` | integer | 否 | 页码，默认 1 |
| `page_size` | integer | 否 | 每页数量，默认 20 |
| `format` | string | 否 | 筛选格式 |
| `status` | string | 否 | 筛选状态 |

**成功响应**：

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

### 下载导出文件

下载已完成的导出文件。

**请求示例**：

```http
GET /api/export/download/{task_id}
Authorization: Bearer {access_token}
```

**响应**：

- `Content-Type: text/csv`（CSV 格式）
- `Content-Type: application/json`（JSON 格式）
- `Content-Disposition: attachment; filename="export_20250101.csv"`

---

### 获取导出任务状态

**请求示例**：

```http
GET /api/export/tasks/{task_id}
Authorization: Bearer {access_token}
```

**成功响应**：

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

**任务状态说明**：

| 状态 | 说明 |
|------|------|
| `pending` | 任务排队中 |
| `processing` | 正在处理 |
| `completed` | 处理完成，可下载 |
| `failed` | 处理失败 |

---

## Webhook 通知 API

### 创建 Webhook 端点

创建 Webhook 接收端点。

**请求示例**：

```http
POST /api/webhooks
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "交易完成通知",
  "url": "https://example.com/webhook",
  "events": ["exchange.completed", "exchange.failed"],
  "is_active": true,
  "headers": {
    "X-Custom-Header": "value"
  }
}
```

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | Webhook 名称 |
| `url` | string | 是 | 接收 URL（HTTPS 推荐） |
| `events` | array | 是 | 订阅的事件类型列表 |
| `is_active` | boolean | 否 | 是否启用，默认 true |
| `headers` | object | 否 | 自定义请求头 |

**成功响应**：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "交易完成通知",
    "url": "https://example.com/webhook",
    "events": ["exchange.completed", "exchange.failed"],
    "is_active": true,
    "secret": "whsec_xxxxxxxxxxxxxxxxxxxx",
    "created_at": "2026-07-06T10:00:00Z"
  }
}
```

> **安全提示**：创建 Webhook 后请妥善保存 `secret`，用于签名验证。该字段仅在创建时返回。

---

### 获取 Webhook 列表

**请求示例**：

```http
GET /api/webhooks
Authorization: Bearer {access_token}
```

---

### 更新 Webhook

**请求示例**：

```http
PUT /api/webhooks/{webhook_id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "更新后的名称",
  "is_active": false
}
```

---

### 删除 Webhook

**请求示例**：

```http
DELETE /api/webhooks/{webhook_id}
Authorization: Bearer {access_token}
```

---

### 获取投递日志

查看 Webhook 投递历史记录。

**请求示例**：

```http
GET /api/webhooks/{webhook_id}/deliveries?page=1&page_size=20
Authorization: Bearer {access_token}
```

**查询参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `page` | integer | 否 | 页码，默认 1 |
| `page_size` | integer | 否 | 每页数量，默认 20 |
| `status` | string | 否 | 投递状态：success/failed/pending |

**成功响应**：

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

### 测试 Webhook

手动触发 Webhook 测试投递。

**请求示例**：

```http
POST /api/webhooks/{webhook_id}/test
Authorization: Bearer {access_token}
```

---

## 公告管理 API

### 获取公告列表

**请求示例**：

```http
GET /api/announcements?page=1&page_size=20
Authorization: Bearer {access_token}
```

---

### 创建公告

**请求示例**：

```http
POST /api/announcements
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "title": "系统维护通知",
  "content": "系统将于今晚进行维护升级...",
  "is_pinned": true,
  "expire_at": "2026-12-31T23:59:59Z"
}
```

---

### 删除公告

**请求示例**：

```http
DELETE /api/announcements/{announcement_id}
Authorization: Bearer {access_token}
```

---

## 通用错误码

| 错误码 | 说明 | 处理建议 |
|--------|------|---------|
| 400 | 请求参数错误 | 检查请求参数格式和必填项 |
| 401 | 未认证或 Token 过期 | 重新登录获取新 Token |
| 403 | 无权限访问 | 检查用户角色和权限 |
| 404 | 资源不存在 | 检查请求路径和资源 ID |
| 409 | 资源冲突 | 检查重复创建或状态冲突 |
| 422 | 业务逻辑校验失败 | 根据错误信息修正业务参数 |
| 429 | 请求过于频繁 | 降低请求频率，遵守速率限制 |
| 500 | 服务器内部错误 | 联系管理员查看日志 |

**通用错误响应格式**：

```json
{
  "code": 400,
  "message": "参数错误：缺少必填字段 'username'",
  "data": null
}
```

---

## 速率限制

系统对 API 请求实施速率限制，防止滥用和暴力攻击。

| 接口类型 | 限制 | 说明 |
|---------|------|------|
| 登录接口 | 5 次/分钟/IP | 防止暴力破解 |
| 其他接口 | 100 次/分钟/用户 | 常规 API 调用 |
| 导出接口 | 10 次/小时/用户 | 导出任务限制 |
| Webhook 投递 | 1000 次/分钟 | 系统级限制 |

**超出限制响应**：

```json
{
  "code": 429,
  "message": "请求过于频繁，请稍后再试",
  "data": {
    "retry_after": 60
  }
}
```

**响应头**：

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1699072800
Retry-After: 60
```

---

## Webhook 签名验证

每个 Webhook 请求都会在 Header 中携带签名，用于验证请求 authenticity。

### 请求头

```http
X-Caiyun-Signature: sha256=<signature>
X-Caiyun-Timestamp: 1699072800
```

### 验证算法

签名生成方式：`HMAC-SHA256(secret, timestamp + "\n" + body)`

**Python 示例**：

```python
import hmac
import hashlib

def verify_signature(secret, timestamp, body, signature):
    """验证 Webhook 签名"""
    expected = hmac.new(
        secret.encode(),
        f"{timestamp}\n{body}".encode(),
        hashlib.sha256
    ).hexdigest()
    return hmac.compare_digest(expected, signature)
```

**Node.js 示例**：

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

**Go 示例**：

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

> **安全提示**：
> 1. 始终验证 `X-Caiyun-Timestamp`，拒绝超过 5 分钟的请求
> 2. 使用 `timingSafeEqual` 或等效函数防止时序攻击
> 3. 仅在 HTTPS 环境下使用 Webhook

---

## Webhook 事件类型

| 事件类型 | 触发时机 | 数据示例 |
|---------|---------|---------|
| `exchange.completed` | 兑换任务完成 | `{"exchange_id": 1, "account_id": 1, "status": "success"}` |
| `exchange.failed` | 兑换任务失败 | `{"exchange_id": 1, "account_id": 1, "error": "Network timeout"}` |
| `exchange.created` | 兑换任务创建 | `{"exchange_id": 1, "account_id": 1, "product_id": "p1"}` |
| `task.completed` | 任务执行完成 | `{"task_id": 1, "type": "signin", "result": "success"}` |
| `task.failed` | 任务执行失败 | `{"task_id": 1, "type": "signin", "error": "Account expired"}` |
| `account.created` | 账号创建 | `{"account_id": 1, "username": "new_account"}` |
| `account.updated` | 账号更新 | `{"account_id": 1, "field": "status", "old": "active", "new": "inactive"}` |
| `account.deleted` | 账号删除 | `{"account_id": 1, "username": "deleted_account"}` |

---

## SDK 与工具

### cURL 示例

```bash
# 登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your_password"}'

# 获取账号列表
curl -X GET http://localhost:8080/api/accounts \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# 创建导出任务
curl -X POST http://localhost:8080/api/export/tasks \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"format":"csv","type":"exchange"}'
```

### 使用 Postman/Insomnia

导入以下环境变量：

```json
{
  "base_url": "http://localhost:8080/api",
  "access_token": "YOUR_ACCESS_TOKEN",
  "refresh_token": "YOUR_REFRESH_TOKEN"
}
```

---

## 📞 技术支持

- **API 问题**：[GitHub Issues](https://github.com/xiaohai09/mobile-cloud-disk-management/issues)
- **安全漏洞**：[Security Advisories](https://github.com/xiaohai09/mobile-cloud-disk-management/security/advisories)

---

*最后更新：2026-07-07*
