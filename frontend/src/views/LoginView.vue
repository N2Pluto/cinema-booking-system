<script setup lang="ts">
const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

function loginWithGoogle() {
  window.location.href = `${apiUrl}/auth/google`
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-[#0f0c29] via-[#302b63] to-[#24243e]">
    <div class="bg-white/5 backdrop-blur-xl border border-white/10 rounded-3xl px-10 py-12 w-full max-w-sm text-center shadow-[0_25px_50px_rgba(0,0,0,0.5)]">

      <!-- Logo -->
      <div class="mb-6 flex justify-center">
        <svg class="w-18 h-18" viewBox="0 0 64 64" fill="none">
          <rect width="64" height="64" rx="16" fill="#1a1a2e"/>
          <rect x="8" y="20" width="48" height="28" rx="4" fill="#16213e"/>
          <rect x="12" y="24" width="40" height="20" rx="2" fill="#0f3460"/>
          <circle cx="22" cy="34" r="6" fill="#e94560"/>
          <path d="M26 34 L20 30 L20 38 Z" fill="white"/>
          <rect x="32" y="28" width="14" height="2" rx="1" fill="#a8b2d8"/>
          <rect x="32" y="32" width="10" height="2" rx="1" fill="#a8b2d8"/>
          <rect x="32" y="36" width="12" height="2" rx="1" fill="#a8b2d8"/>
          <rect x="6"  y="16" width="4" height="6" rx="1" fill="#e94560"/>
          <rect x="54" y="16" width="4" height="6" rx="1" fill="#e94560"/>
          <rect x="6"  y="44" width="4" height="6" rx="1" fill="#e94560"/>
          <rect x="54" y="44" width="4" height="6" rx="1" fill="#e94560"/>
        </svg>
      </div>

      <h1 class="text-3xl font-bold text-white tracking-tight mb-2">Cinema Booking</h1>
      <p class="text-[#a8b2d8] text-sm mb-9">จองตั๋วหนังออนไลน์ง่ายๆ ที่นี่ที่เดียว</p>

      <!-- Google Login Button -->
      <button
        @click="loginWithGoogle"
        class="flex items-center justify-center gap-3 w-full py-3.5 bg-white text-[#1f1f1f] rounded-xl font-semibold text-sm shadow-[0_4px_12px_rgba(0,0,0,0.3)] hover:-translate-y-0.5 hover:shadow-[0_8px_20px_rgba(0,0,0,0.4)] active:translate-y-0 transition-all cursor-pointer"
      >
        <svg class="w-5 h-5 shrink-0" viewBox="0 0 24 24">
          <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
          <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
          <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
          <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
        </svg>
        เข้าสู่ระบบด้วย Google
      </button>
    </div>
  </div>
</template>
