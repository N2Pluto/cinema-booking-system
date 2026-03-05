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
  <div class="min-h-screen flex flex-col items-center justify-center gap-5 bg-gradient-to-br from-[#0f0c29] via-[#302b63] to-[#24243e] text-[#a8b2d8]">
    <div class="w-10 h-10 rounded-full border-3 border-white/10 border-t-[#e94560] animate-spin"></div>
    <p class="text-base">กำลังเข้าสู่ระบบ...</p>
  </div>
</template>
