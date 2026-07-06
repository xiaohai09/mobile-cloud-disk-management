<template>
  <div class="webhook-center">
    <page-header
      title="Webhook 管理"
      :subtitle="`共 ${total || 0} 个端点`"
    >
      <template #extra>
        <el-button
          type="primary"
          @click="showCreateDialog = true"
        >
          <el-icon><Plus /></el-icon>
          新建端点
        </el-button>
      </template>
    </page-header>

    <el-card class="webhook-list-card">
      <template #header>
        <div class="card-header">
          <span>Webhook 端点</span>
          <el-button
            text
            :loading="loading"
            @click="loadEndpoints"
          >
            <el-icon><Refresh /></el-icon>
          </el-button>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="endpoints"
        empty-text="暂无 Webhook 端点"
      >
        <el-table-column
          prop="id"
          label="ID"
          width="80"
        />
        <el-table-column
          prop="name"
          label="名称"
          min-width="160"
        />
        <el-table-column
          prop="url"
          label="URL"
          min-width="240"
        >
          <template #default="{ row }">
            <el-link
              :href="row.url"
              target="_blank"
              type="primary"
            >
              {{ row.url }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column
          prop="events"
          label="事件"
          min-width="220"
        >
          <template #default="{ row }">
            <el-tag
              v-for="event in row.events"
              :key="event"
              size="small"
              class="event-tag"
            >
              {{ eventLabels[event] || event }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="is_active"
          label="状态"
          width="100"
        >
          <template #default="{ row }">
            <el-tag
              :type="row.is_active ? 'success' : 'info'"
              size="small"
            >
              {{ row.is_active ? '启用' : '停用' }}
            </el-tag>
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
          label="操作"
          width="220"
          fixed="right"
        >
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              size="small"
              @click="openEditDialog(row)"
            >
              编辑
            </el-button>
            <el-button
              type="success"
              link
              size="small"
              @click="testEndpoint(row)"
            >
              测试
            </el-button>
            <el-button
              type="danger"
              link
              size="small"
              @click="deleteEndpoint(row)"
            >
              删除
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

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="showDialog"
      :title="editingEndpoint ? '编辑 Webhook 端点' : '新建 Webhook 端点'"
      width="560px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
      >
        <el-form-item
          label="名称"
          prop="name"
        >
          <el-input
            v-model="form.name"
            placeholder="例如：飞书通知"
          />
        </el-form-item>

        <el-form-item
          label="URL"
          prop="url"
        >
          <el-input
            v-model="form.url"
            placeholder="https://example.com/webhook"
          />
        </el-form-item>

        <el-form-item label="事件">
          <el-select
            v-model="form.events"
            multiple
            filterable
            placeholder="请选择事件"
            style="width: 100%"
          >
            <el-option
              label="任务成功"
              value="task.success"
            />
            <el-option
              label="任务失败"
              value="task.failure"
            />
            <el-option
              label="抢兑命中"
              value="exchange.hit"
            />
            <el-option
              label="系统告警"
              value="system.alert"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="密钥">
          <el-input
            v-model="form.secret"
            type="textarea"
            :rows="2"
            placeholder="可选，用于签名验证"
          />
        </el-form-item>

        <el-form-item label="自定义 Headers">
          <el-input
            v-model="headersText"
            type="textarea"
            :rows="3"
            placeholder="{&quot;X-Custom&quot;: &quot;value&quot;}"
          />
        </el-form-item>

        <el-form-item label="启用">
          <el-switch v-model="form.is_active" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showDialog = false; resetForm()">
          取消
        </el-button>
        <el-button
          type="primary"
          :loading="submitting"
          @click="submitForm"
        >
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import PageHeader from '@/components/PageHeader.vue'
import {
  getWebhookEndpoints,
  createWebhookEndpoint,
  updateWebhookEndpoint,
  deleteWebhookEndpoint,
  testWebhookEndpoint,
  type WebhookEndpoint,
  type WebhookEventType
} from '@/api/webhook'

const loading = ref(false)
const submitting = ref(false)
const showDialog = ref(false)
const showCreateDialog = ref(false)
const editingEndpoint = ref<WebhookEndpoint | null>(null)
const endpoints = ref<WebhookEndpoint[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const formRef = ref<FormInstance | null>(null)

const eventLabels: Record<string, string> = {
  'task.success': '任务成功',
  'task.failure': '任务失败',
  'exchange.hit': '抢兑命中',
  'system.alert': '系统告警'
}

const form = reactive({
  name: '',
  url: '',
  events: [] as WebhookEventType[],
  secret: '',
  is_active: true
})

const headersText = ref('')

const rules: FormRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  url: [
    { required: true, message: '请输入 URL', trigger: 'blur' },
    { type: 'url', message: '请输入有效的 URL', trigger: 'blur' }
  ]
}

const resetForm = () => {
  editingEndpoint.value = null
  form.name = ''
  form.url = ''
  form.events = []
  form.secret = ''
  form.is_active = true
}

const loadEndpoints = async () => {
  loading.value = true
  try {
    const res = await getWebhookEndpoints()
    endpoints.value = res.endpoints
    total.value = res.total
  } catch (error) {
    ElMessage.error('加载 Webhook 端点失败')
  } finally {
    loading.value = false
  }
}

const openEditDialog = (endpoint: WebhookEndpoint) => {
  editingEndpoint.value = endpoint
  form.name = endpoint.name
  form.url = endpoint.url
  form.events = [...endpoint.events]
  form.secret = endpoint.secret || ''
  form.is_active = endpoint.is_active
  headersText.value = JSON.stringify(endpoint.headers || {}, null, 2)
  showDialog.value = true
  showCreateDialog.value = false
}

const submitForm = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      let headers: Record<string, string> = {}
      if (headersText.value.trim()) {
        try {
          headers = JSON.parse(headersText.value)
        } catch {
          ElMessage.error('Headers JSON 格式错误')
          return
        }
      }

      if (editingEndpoint.value) {
        await updateWebhookEndpoint(editingEndpoint.value.id, {
          name: form.name,
          url: form.url,
          events: form.events,
          secret: form.secret || undefined,
          headers,
          is_active: form.is_active
        })
        ElMessage.success('更新成功')
      } else {
        await createWebhookEndpoint({
          name: form.name,
          url: form.url,
          events: form.events,
          secret: form.secret || undefined,
          headers,
          is_active: form.is_active
        })
        ElMessage.success('创建成功')
      }

      showDialog.value = false
      await loadEndpoints()
    } catch (error) {
      ElMessage.error(editingEndpoint.value ? '更新失败' : '创建失败')
    } finally {
      submitting.value = false
    }
  })
}

const testEndpoint = async (endpoint: WebhookEndpoint) => {
  try {
    await testWebhookEndpoint(endpoint.id)
    ElMessage.success('测试事件已发送')
  } catch (error) {
    ElMessage.error('测试失败')
  }
}

const deleteEndpoint = async (endpoint: WebhookEndpoint) => {
  try {
    await ElMessageBox.confirm('确定要删除这个 Webhook 端点吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteWebhookEndpoint(endpoint.id)
    ElMessage.success('删除成功')
    await loadEndpoints()
  } catch {
    // 用户取消
  }
}

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const handleSizeChange = () => {
  page.value = 1
  loadEndpoints()
}

const handlePageChange = () => {
  loadEndpoints()
}

onMounted(() => {
  loadEndpoints()
})
</script>

<style scoped>
.webhook-center {
  padding: 20px;
}

.webhook-list-card {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.event-tag {
  margin-right: 4px;
  margin-bottom: 4px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 768px) {
  .webhook-center {
    padding: 12px;
  }

  .webhook-list-card {
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
