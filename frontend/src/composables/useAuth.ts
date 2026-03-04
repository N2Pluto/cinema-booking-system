import { ref, computed } from 'vue'

const TOKEN_KEY = 'cinema_token'

interface User {
  user_id: string
  email: string
  role: string
}

const user = ref<User | null>(null)

export function useAuth() {
  const token = computed(() => localStorage.getItem(TOKEN_KEY))
  const isLoggedIn = computed(() => !!token.value)

  function saveToken(t: string) {
    localStorage.setItem(TOKEN_KEY, t)
  }

  function logout() {
    localStorage.removeItem(TOKEN_KEY)
    user.value = null
  }

  async function fetchMe(): Promise<User | null> {
    const t = token.value
    if (!t) return null

    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const res = await fetch(`${apiUrl}/auth/me`, {
      headers: { Authorization: `Bearer ${t}` },
    })

    if (!res.ok) {
      logout()
      return null
    }

    const data = await res.json()
    user.value = data.user_id ? data : data.user
    return user.value
  }

  return { token, isLoggedIn, user, saveToken, logout, fetchMe }
}
