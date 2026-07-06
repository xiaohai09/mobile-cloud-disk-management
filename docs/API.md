# Caiyun Refactored - API 文档

## 基础信息

- Base URL: `/api`
- 认证方式: Bearer JWT / Cookie `auth_token`
- 响应格式: JSON

## 数据导出 API

### 导出数据

```
GET /api/export
```

**查询参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `type` | string | 是 | 导出类型：`accounts`、`tasks`、`exchange`、`records` |
| `format` | string | 否 | 导出格式：`csv` 或 `json`，默认 `csv` |
| `account_id` | int | 否 | 筛选账号 ID |
| `start_date` | string | 否 | 开始日期 `YYYY-MM-DD` |
| `end_date` | string | 否 | 结束日期 `YYYY-MM-DD` |
| `status` | string | 否 | 状态筛选 |

**响应**

```json
{
  "id": 1,
  "file_path": "/exports/accounts_1_20260706.csv",
  "file_size": 1024,
  "status": "pending"
}
```

### 获取导出历史

```
GET /api/export/history
```

**查询参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `page` | int | 否 | 页码，默认 1 |
| `page_size` | int | 否 | 每页条数，默认 20，最大 100 |

**响应**

```json
{
  "data": [],
  "total": 0,
  "page": 1,
  "page_size": 20
}
```

## Webhook API

### 获取端点列表

```
GET /api/webhooks
```

**查询参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `page` | int | 否 | 页码，默认 1 |
| `page_size` | int | 否 | 每页条数，默认 20 |

**响应**

```json
{
  "endpoints": [],
  "total": 0,
  "page": 1,
  "page_size": 20
}
```

### 创建端点

```
POST /api/webhooks
```

**请求体**

```json
{
  "name": "飞书通知",
  "url": "https://open.feishu.cn/open-apis/bot/v2/hook/xxx",
  "events": ["task.success", "task.failed"],
  "secret": "optional-signing-secret",
  "headers": {
    "X-Custom": "value"
  },
  "is_active": true
}
```

### 更新端点

```
PUT /api/webhooks/:id
```

**请求体**

```json
{
  "name": "新名称",
  "url": "https://example.com/hook",
  "events": ["task.success"],
  "secret": "new-secret",
  "headers": {},
  "is_active": false
}
```

### 删除端点

```
DELETE /api/webhooks/:id
```

### 测试端点

```
POST /api/webhooks/:id/test
```

**响应**

```json
{
  "success": true,
  "message": "test event sent"
}
```

## 平台适配 API

```
GET /api/platforms
```

**响应**

```json
{
  "platforms": [
    {
      "id": "caiyun",
      "name": "移动云盘",
      "description": "中国移动云盘"
    }
  ]
}
```

## 通用错误码

| HTTP 状态码 | 说明 |
|-------------|------|
| 400 | 请求参数错误 |
| 401 | 未认证或 token 过期 |
| 403 | 无权限访问 |
| 404 | 资源不存在 |
| 429 | 请求过于频繁 |
| 500 | 服务器内部错误 |

## 安全说明

- 所有 `/api/*` 业务接口均要求认证 + CSRF
- 导出/Webhook 接口额外经过限流中间件
- 审计日志自动记录请求/响应（敏感字段已脱敏）
- Webhook 推送支持 HMAC-SHA256 签名验证
