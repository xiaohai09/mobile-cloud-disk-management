import request from './axios'
import { unwrapApiData, type ApiResponse } from './response'

// ==================== 类型定义 ====================

export type WebhookEventType = 'task.success' | 'task.failure' | 'exchange.hit' | 'system.alert'

export interface WebhookEndpoint {
  id: number
  user_id: number
  name: string
  url: string
  events: WebhookEventType[]
  secret?: string
  headers?: Record<string, string>
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CreateWebhookEndpointRequest {
  name: string
  url: string
  events?: WebhookEventType[]
  secret?: string
  headers?: Record<string, string>
  is_active?: boolean
}

export interface UpdateWebhookEndpointRequest {
  name?: string
  url?: string
  events?: WebhookEventType[]
  secret?: string
  headers?: Record<string, string>
  is_active?: boolean
}

export interface WebhookEndpointsResponse {
  endpoints: WebhookEndpoint[]
  total: number
}

// ==================== API 函数 ====================

// 获取 Webhook 列表
export function getWebhookEndpoints(): Promise<WebhookEndpointsResponse> {
  const fallback: WebhookEndpointsResponse = { endpoints: [], total: 0 }
  return request<WebhookEndpointsResponse | ApiResponse<WebhookEndpointsResponse>>({
    url: '/api/webhooks',
    method: 'get'
  }).then((res) => unwrapApiData(res, fallback))
}

// 创建 Webhook
export function createWebhookEndpoint(data: CreateWebhookEndpointRequest): Promise<WebhookEndpoint> {
  const fallback = { id: 0 } as WebhookEndpoint
  return request<WebhookEndpoint | ApiResponse<WebhookEndpoint>>({
    url: '/api/webhooks',
    method: 'post',
    data
  }).then((res) => unwrapApiData(res, fallback))
}

// 更新 Webhook
export function updateWebhookEndpoint(id: number, data: UpdateWebhookEndpointRequest): Promise<{ id: number; updated: boolean }> {
  return request({
    url: `/api/webhooks/${id}`,
    method: 'put',
    data
  })
}

// 删除 Webhook
export function deleteWebhookEndpoint(id: number): Promise<{ id: number; deleted: boolean }> {
  return request({
    url: `/api/webhooks/${id}`,
    method: 'delete'
  })
}

// 测试 Webhook
export function testWebhookEndpoint(id: number): Promise<{ id: number; success: boolean }> {
  return request({
    url: `/api/webhooks/${id}/test`,
    method: 'post'
  })
}
