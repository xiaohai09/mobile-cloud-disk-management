import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export type ThemeMode = 'light' | 'dark' | 'auto'

export const useThemeStore = defineStore('theme', () => {
  const mode = ref<ThemeMode>(getInitialTheme())
  const isDark = ref(false)
  const sidebarCollapsed = ref(false)

  function getInitialTheme(): ThemeMode {
    const saved = localStorage.getItem('theme_mode') as ThemeMode | null
    return saved || 'auto'
  }

  function updateEffectiveTheme() {
    if (mode.value === 'auto') {
      isDark.value = window.matchMedia('(prefers-color-scheme: dark)').matches
    } else {
      isDark.value = mode.value === 'dark'
    }
    applyThemeClass()
  }

  function applyThemeClass() {
    const html = document.documentElement
    if (isDark.value) {
      html.classList.add('dark')
    } else {
      html.classList.remove('dark')
    }
  }

  function setMode(newMode: ThemeMode) {
    mode.value = newMode
    localStorage.setItem('theme_mode', newMode)
    updateEffectiveTheme()
  }

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  watch(
    () => mode.value,
    () => {
      updateEffectiveTheme()
      localStorage.setItem('theme_mode', mode.value)
    }
  )

  updateEffectiveTheme()

  return {
    mode,
    isDark,
    sidebarCollapsed,
    setMode,
    toggleSidebar,
    updateEffectiveTheme,
  }
})
