<template>
  <div class="export-center">
    <page-header
      title="数据导出"
      :subtitle="`共 ${total || 0} 条记录`"
    >
      <template #extra>
        <el-button
          type="primary"
          @click="showExportDialog = true"
        >
          <el-icon><Download /></el-icon>
          导出数据
        </el-button>
      </template>
    </page-header>

    <el-card class="export-history-card">
      <template #header>
        <div class="card-header">
          <span>导出历史</span>
          <el-button
            text
            :loading="loading"
            @click="loadHistory"
          >
            <el-icon><Refresh /></el-icon>
          </el-button>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="history"
        empty-text="暂无导出记录"
      >
        <el-table-column
          prop="id"
          label="ID"
          width="80"
        />
        <el-table-column
          prop="type"
          label="类型"
          width="120"
        >
          <template #default="{ row }">
            <el-tag size="small">
              {{ typeLabels[row.type] || row.type }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="format"
          label="格式"
          width="100"
        >
          <template #default="{ row }">
            <el-tag
              type="info"
              size="small"
            >
              {{ row.format?.toUpperCase() }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="status"
          label="状态"
          width="100"
        >
          <template #default="{ row }">
            <el-tag
              :type="statusTagType(row.status)"
              size="small"
            >
              {{ statusLabels[row.status] || row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="file_size"
          label="文件大小"
          width="120"
        >
          <template #default="{ row }">
            {{ formatFileSize(row.file_size) }}
          </template>
        </el-table-column>
        <el-table-column
          prop="created_at"
          label="创建时间"
          min-width="180"
        >
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column
          prop="completed_at"
          label="完成时间"
          min-width="180"
        >
          <template #default="{ row }">
            {{ row.completed_at ? formatDateTime(row.completed_at) : '-' }}
          </template>
        </el-table-column>
        <el-table-column
          label="操作"
          width="100"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'completed'"
              type="primary"
              link
              size="small"
              @click="downloadExport"
            >
              下载
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 导出对话框 -->
    <el-dialog
      v-model="showExportDialog"
      title="导出数据"
      width="520px"
    >
      <el-form
        :model="exportForm"
        label-width="100px"
      >
        <el-form-item
          label="数据类型"
          required
        >
          <el-select
            v-model="exportForm.type"
            placeholder="请选择数据类型"
            style="width: 100%"
          >
            <el-option
              label="任务日志"
              value="task_logs"
            />
            <el-option
              label="云朵统计"
              value="cloud_stats"
            />
            <el-option
              label="抢兑记录"
              value="exchange_records"
            />
            <el-option
              label="账号列表"
              value="accounts"
            />
            <el-option
              label="全部数据"
              value="all"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="导出格式">
          <el-radio-group v-model="exportForm.format">
            <el-radio value="csv">
              CSV
            </el-radio>
            <el-radio value="json">
              JSON
            </el-radio>
            <el-radio value="xlsx">
              XLSX
            </el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="账号筛选">
          <el-input
            v-model="exportForm.account_id"
            placeholder="账号ID（可选）"
            type="number"
          />
        </el-form-item>

        <el-form-item label="开始日期">
          <el-date-picker
            v-model="exportForm.start_date"
            type="date"
            placeholder="开始日期"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="结束日期">
          <el-date-picker
            v-model="exportForm.end_date"
            type="date"
            placeholder="结束日期"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="状态筛选">
          <el-select
            v-model="exportForm.status"
            placeholder="全部状态"
            clearable
            style="width: 100%"
          >
            <el-option
              label="成功"
              value="success"
            />
            <el-option
              label="失败"
              value="failed"
            />
            <el-option
              label="运行中"
              value="running"
            />
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showExportDialog = false">
          取消
        </el-button>
        <el-button
          type="primary"
          :loading="exporting"
          @click="handleExport"
        >
          开始导出
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Download, Refresh } from '@element-plus/icons-vue'
import PageHeader from '@/components/PageHeader.vue'
import { exportData, getExportHistory, type ExportHistoryItem, type ExportType, type ExportFormat } from '@/api/export'

const loading = ref(false)
const exporting = ref(false)
const showExportDialog = ref(false)
const history = ref<ExportHistoryItem[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

const typeLabels: Record<string, string> = {
  task_logs: '任务日志',
  cloud_stats: '云朵统计',
  exchange_records: '抢兑记录',
  accounts: '账号列表',
  all: '全部数据'
}

const statusLabels: Record<string, string> = {
  pending: '等待中',
  processing: '处理中',
  completed: '已完成',
  failed: '失败'
}

const statusTagType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'info',
    processing: 'warning',
    completed: 'success',
    failed: 'danger'
  }
  return map[status] || 'info'
}

const exportForm = ref({
  type: 'task_logs' as ExportType,
  format: 'csv' as ExportFormat,
  account_id: undefined as number | undefined,
  start_date: '',
  end_date: '',
  status: ''
})

const loadHistory = async () => {
  loading.value = true
  try {
    const res = await getExportHistory()
    history.value = res.exports
    total.value = res.total
  } catch (error) {
    ElMessage.error('加载导出历史失败')
  } finally {
    loading.value = false
  }
}

const handleExport = async () => {
  if (!exportForm.value.type) {
    ElMessage.warning('请选择数据类型')
    return
  }

  exporting.value = true
  try {
    const blob = await exportData({
      type: exportForm.value.type,
      format: exportForm.value.format,
      account_id: exportForm.value.account_id,
      start_date: exportForm.value.start_date || undefined,
      end_date: exportForm.value.end_date || undefined,
      status: exportForm.value.status || undefined
    })

    // Download the file
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `export_${exportForm.value.type}_${Date.now()}.${exportForm.value.format}`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)

    ElMessage.success('导出成功')
    showExportDialog.value = false
    await loadHistory()
  } catch (error) {
    ElMessage.error('导出失败')
  } finally {
    exporting.value = false
  }
}

const downloadExport = () => {
  ElMessage.info('下载功能需要服务端支持文件访问')
}

const formatFileSize = (bytes?: number) => {
  if (!bytes) return '-'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
}

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const handleSizeChange = () => {
  page.value = 1
  loadHistory()
}

const handlePageChange = () => {
  loadHistory()
}

onMounted(() => {
  loadHistory()
})
</script>

<style scoped>
.export-center {
  padding: 20px;
}

.export-history-card {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 768px) {
  .export-center {
    padding: 12px;
  }

  .export-history-card {
    margin-top: 12px;
  }

  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  :deep(.el-table) {
    font-size: 12px;
  }

  :deep(.el-table .cell) {
    padding: 8px 4px;
  }

  :deep(.el-pagination) {
    flex-wrap: wrap;
    gap: 8px;
  }
}
</style>
