import request from './axios'
import { unwrapApiData, type ApiResponse } from './response'

// ==================== 类型定义 ====================

export type ExportFormat = 'csv' | 'json' | 'xlsx'
export type ExportType = 'task_logs' | 'cloud_stats' | 'exchange_records' | 'accounts' | 'all'

export interface ExportHistoryItem {
  id: number
  user_id: number
  type: ExportType
  format: ExportFormat
  status: string
  file_path?: string
  file_size?: number
  created_at: string
  completed_at?: string
}

export interface ExportHistoryResponse {
  exports: ExportHistoryItem[]
  total: number
}

// ==================== API 函数 ====================

// 导出数据
export function exportData(params: {
  type: ExportType
  format?: ExportFormat
  account_id?: number
  start_date?: string
  end_date?: string
  status?: string
}): Promise<Blob> {
  return request<Blob>({
    url: '/api/export',
    method: 'get',
    params,
    responseType: 'blob'
  })
}

// 获取导出历史
export function getExportHistory(): Promise<ExportHistoryResponse> {
  const fallback: ExportHistoryResponse = { exports: [], total: 0 }
  return request<ExportHistoryResponse | ApiResponse<ExportHistoryResponse>>({
    url: '/api/export/history',
    method: 'get'
  }).then((res) => unwrapApiData(res, fallback))
}
