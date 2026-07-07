import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/store/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/Login.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/Register.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/forgot-password',
      name: 'ForgotPassword',
      component: () => import('@/views/ForgotPassword.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/',
      name: 'Layout',
      component: () => import('@/views/Layout.vue'),
      meta: { requiresAuth: true },
      redirect: '/dashboard',
      children: [
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('@/views/Dashboard.vue'),
          meta: { title: '首页' }
        },
        {
          path: 'accounts',
          name: 'Accounts',
          component: () => import('@/views/AccountManage.vue'),
          meta: { title: '账号' }
        },
        {
          path: 'logs',
          name: 'Logs',
          component: () => import('@/views/TaskLogs.vue'),
          meta: { title: '日志' }
        },
        {
          path: 'admin',
          name: 'Admin',
          component: () => import('@/views/AdminPanel.vue'),
          meta: { title: '管理', requiresAdmin: true }
        },
        {
          path: 'exchange',
          name: 'Exchange',
          component: () => import('@/views/ExchangeCenter.vue'),
          meta: { title: '兑换' }
        },
        {
          path: 'exchange/records',
          name: 'ExchangeRecords',
          component: () => import('@/views/ExchangeRecords.vue'),
          meta: { title: '抢兑记录' }
        },
        {
          path: 'export',
          name: 'Export',
          component: () => import('@/views/ExportCenter.vue'),
          meta: { title: '数据导出' }
        },
        {
          path: 'webhooks',
          name: 'Webhooks',
          component: () => import('@/views/WebhookCenter.vue'),
          meta: { title: 'Webhook' }
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFound.vue'),
      meta: { requiresAuth: false }
    }
  ]
})

function safeRedirect(value: unknown): string {
  const redirect = Array.isArray(value) ? value[0] : value
  if (typeof redirect !== 'string') {
    return '/'
  }
  if (!redirect.startsWith('/') || redirect.startsWith('//') || redirect.includes('\\')) {
    return '/'
  }
  return redirect
}

router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()
  const isProtectedRoute = to.meta.requiresAuth
  const isAuthPage = to.name === 'Login' || to.name === 'Register'
  const shouldProbeSession = !authStore.hasCheckedSession
    && (isProtectedRoute || isAuthPage)

  // 周期性重验证：已认证但上次验证超过 5 分钟，重新探测会话，避免 Cookie 过期后仍保持登录态。
  const shouldRevalidate = isProtectedRoute
    && authStore.isAuthenticated
    && (!authStore.lastVerified || Date.now() - authStore.lastVerified > 5 * 60 * 1000)

  if (shouldProbeSession || shouldRevalidate) {
    try {
      await authStore.refreshProfile()
    } catch {
      // 未登录或 Cookie 已失效时保持匿名状态；/api/auth/me 的 401 不弹全局错误。
    }
  }

  const isAuthenticated = authStore.isAuthenticated
  const userRole = authStore.user?.role

  if (to.meta.requiresAuth && !isAuthenticated) {
    // 未登录访问受保护页面 → 携带 redirect，登录后可回到原页面。
    next({ name: 'Login', query: to.fullPath !== '/' ? { redirect: to.fullPath } : undefined })
  } else if (to.meta.requiresAdmin && userRole !== 'admin') {
    next('/')
  } else if ((to.name === 'Login' || to.name === 'Register') && isAuthenticated) {
    // 已登录访问登录/注册页：优先跳到 redirect，否则首页。
    next(safeRedirect(to.query.redirect))
  } else {
    next()
  }
})

export default router
