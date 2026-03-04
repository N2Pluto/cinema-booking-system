<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const { saveToken } = useAuth()

onMounted(() => {
  const params = new URLSearchParams(window.location.search)
  const token = params.get('token')

  if (token) {
    saveToken(token)
    router.replace({ name: 'home' })
  } else {
    router.replace({ name: 'login' })
  }
})
</script>

<template>
  <div class="callback-page">
    <div class="spinner"></div>
    <p>กำลังเข้าสู่ระบบ...</p>
  </div>
</template>

<style scoped>
.callback-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 20px;
  background: linear-gradient(135deg, #0f0c29, #302b63, #24243e);
  color: #a8b2d8;
  font-size: 16px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(255,255,255,0.1);
  border-top-color: #e94560;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
