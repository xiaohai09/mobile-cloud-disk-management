<template>
  <div class="mobile-bottom-nav">
    <div
      v-for="item in menuItems"
      :key="item.path"
      class="nav-item"
      :class="{ active: isActive(item.path) }"
      @click="navigateTo(item.path)"
    >
      <el-icon :size="20">
        <component :is="item.icon" />
      </el-icon>
      <span class="nav-label">{{ item.label }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  DataLine,
  User,
  List,
  Shop,
  Download,
  Bell,
  Setting
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const isAdmin = computed(() => authStore.user?.role === 'admin')

const menuItems = computed(() => {
  const items = [
    { path: '/dashboard', label: '首页', icon: DataLine },
    { path: '/accounts', label: '账号', icon: User },
    { path: '/logs', label: '日志', icon: List },
    { path: '/exchange', label: '兑换', icon: Shop },
    { path: '/export', label: '导出', icon: Download },
    { path: '/webhooks', label: '通知', icon: Bell }
  ]

  if (isAdmin.value) {
    items.push({ path: '/admin', label: '管理', icon: Setting })
  }

  return items
})

const isActive = (path: string) => {
  if (path === '/dashboard') {
    return route.path === '/' || route.path === '/dashboard'
  }
  return route.path.startsWith(path)
}

const navigateTo = (path: string) => {
  router.push(path)
}
</script>

<style scoped>
.mobile-bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 56px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-top: 1px solid rgba(59, 130, 246, 0.15);
  display: none;
  justify-content: space-around;
  align-items: center;
  padding: 0 4px;
  z-index: 1000;
  box-shadow: 0 -4px 20px rgba(59, 130, 246, 0.1);
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  padding: 6px 8px;
  min-width: 56px;
  color: #94a3b8;
  transition: all 0.25s ease;
  cursor: pointer;
  border-radius: 12px;
  user-select: none;
}

.nav-item:active {
  transform: scale(0.95);
}

.nav-item.active {
  color: #2563eb;
}

.nav-label {
  font-size: 11px;
  font-weight: 500;
  white-space: nowrap;
}

@media (max-width: 768px) {
  .mobile-bottom-nav {
    display: flex;
  }
}
</style>
