import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './plugins/element-styles'

// 引入移动端适配样式
import './styles/mobile.scss'
import { installElementPlus } from './plugins/element'

import App from './App.vue'
import router from './router'
import { i18n } from './i18n'

const app = createApp(App)

const pinia = createPinia()
app.use(pinia)
app.use(i18n)

// Initialize auth store from localStorage (must be after pinia is installed)
import { useAuthStore } from './store/auth'
import { purgeSensitiveStorage } from './utils/security'
const authStore = useAuthStore()
authStore.initialize()
purgeSensitiveStorage()

app.use(router)
installElementPlus(app)

app.mount('#app')
