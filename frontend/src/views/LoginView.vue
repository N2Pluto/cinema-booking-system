<script setup lang="ts">
const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

function loginWithGoogle() {
  window.location.href = `${apiUrl}/auth/google`
}
</script>

<template>
  <div class="login-page">
    <div class="login-card">
      <div class="logo-area">
        <svg class="cinema-icon" viewBox="0 0 64 64" fill="none">
          <rect width="64" height="64" rx="16" fill="#1a1a2e"/>
          <rect x="8" y="20" width="48" height="28" rx="4" fill="#16213e"/>
          <rect x="12" y="24" width="40" height="20" rx="2" fill="#0f3460"/>
          <circle cx="22" cy="34" r="6" fill="#e94560"/>
          <path d="M26 34 L20 30 L20 38 Z" fill="white"/>
          <rect x="32" y="28" width="14" height="2" rx="1" fill="#a8b2d8"/>
          <rect x="32" y="32" width="10" height="2" rx="1" fill="#a8b2d8"/>
          <rect x="32" y="36" width="12" height="2" rx="1" fill="#a8b2d8"/>
          <rect x="6" y="16" width="4" height="6" rx="1" fill="#e94560"/>
          <rect x="54" y="16" width="4" height="6" rx="1" fill="#e94560"/>
          <rect x="6" y="44" width="4" height="6" rx="1" fill="#e94560"/>
          <rect x="54" y="44" width="4" height="6" rx="1" fill="#e94560"/>
        </svg>
      </div>

      <h1>Cinema Booking</h1>
      <p class="subtitle">จองตั๋วหนังออนไลน์ง่ายๆ ที่นี่ที่เดียว</p>

      <button class="google-btn" @click="loginWithGoogle">
        <svg class="google-icon" viewBox="0 0 24 24">
          <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
          <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
          <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
          <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
        </svg>
        เข้าสู่ระบบด้วย Google
      </button>

      <p class="note">ระบบใช้ Google OAuth 2.0 — ปลอดภัย ไม่เก็บรหัสผ่าน</p>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #0f0c29, #302b63, #24243e);
}

.login-card {
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 24px;
  padding: 48px 40px;
  width: 100%;
  max-width: 400px;
  text-align: center;
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.5);
}

.logo-area {
  margin-bottom: 24px;
}

.cinema-icon {
  width: 72px;
  height: 72px;
}

h1 {
  font-size: 28px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 8px;
  letter-spacing: -0.5px;
}

.subtitle {
  color: #a8b2d8;
  font-size: 15px;
  margin: 0 0 36px;
}

.google-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  width: 100%;
  padding: 14px 24px;
  background: #fff;
  color: #1f1f1f;
  border: none;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.google-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.4);
}

.google-btn:active {
  transform: translateY(0);
}

.google-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.note {
  margin-top: 20px;
  font-size: 12px;
  color: #606882;
}
</style>
