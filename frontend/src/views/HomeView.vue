<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const { user, fetchMe, logout } = useAuth()

onMounted(async () => {
  await fetchMe()
})

function handleLogout() {
  logout()
  router.push({ name: 'login' })
}
</script>

<template>
  <div class="home-page">
    <nav class="navbar">
      <div class="nav-brand">
        <span class="brand-dot"></span>
        Cinema Booking
      </div>
      <button class="logout-btn" @click="handleLogout">ออกจากระบบ</button>
    </nav>

    <main class="content">
      <div class="welcome-card" v-if="user">
        <div class="avatar">{{ user.email?.charAt(0).toUpperCase() }}</div>
        <h2>ยินดีต้อนรับ 👋</h2>
        <p class="email">{{ user.email }}</p>
        <span class="role-badge" :class="user.role?.toLowerCase()">
          {{ user.role }}
        </span>
      </div>

      <div class="section-title">
        <h3>ภาพยนตร์ที่กำลังฉาย</h3>
        <p>เลือกรอบที่ต้องการแล้วจองได้เลย</p>
      </div>

      <div class="placeholder-grid">
        <div class="movie-card placeholder" v-for="n in 6" :key="n">
          <div class="poster-placeholder"></div>
          <div class="info-placeholder">
            <div class="line short"></div>
            <div class="line long"></div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
.home-page {
  min-height: 100vh;
  background: #0d0d1a;
  color: #fff;
}

.navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 32px;
  background: rgba(255,255,255,0.03);
  border-bottom: 1px solid rgba(255,255,255,0.08);
  backdrop-filter: blur(10px);
  position: sticky;
  top: 0;
  z-index: 10;
}

.nav-brand {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 18px;
  font-weight: 700;
  letter-spacing: -0.3px;
}

.brand-dot {
  width: 10px;
  height: 10px;
  background: #e94560;
  border-radius: 50%;
}

.logout-btn {
  padding: 8px 18px;
  background: rgba(233, 69, 96, 0.15);
  color: #e94560;
  border: 1px solid rgba(233, 69, 96, 0.3);
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.logout-btn:hover {
  background: rgba(233, 69, 96, 0.25);
}

.content {
  max-width: 1100px;
  margin: 0 auto;
  padding: 40px 24px;
}

.welcome-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 32px;
  background: linear-gradient(135deg, rgba(233,69,96,0.1), rgba(48,43,99,0.3));
  border: 1px solid rgba(233,69,96,0.2);
  border-radius: 20px;
  margin-bottom: 40px;
  text-align: center;
}

.avatar {
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, #e94560, #302b63);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  font-weight: 700;
}

h2 {
  margin: 0;
  font-size: 22px;
}

.email {
  color: #a8b2d8;
  margin: 0;
  font-size: 14px;
}

.role-badge {
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.role-badge.user {
  background: rgba(100,200,150,0.15);
  color: #64c896;
  border: 1px solid rgba(100,200,150,0.3);
}

.role-badge.admin {
  background: rgba(233,69,96,0.15);
  color: #e94560;
  border: 1px solid rgba(233,69,96,0.3);
}

.section-title {
  margin-bottom: 24px;
}

.section-title h3 {
  font-size: 20px;
  margin: 0 0 4px;
}

.section-title p {
  color: #606882;
  font-size: 14px;
  margin: 0;
}

.placeholder-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 20px;
}

.movie-card.placeholder {
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.07);
  border-radius: 16px;
  overflow: hidden;
}

.poster-placeholder {
  width: 100%;
  aspect-ratio: 2/3;
  background: linear-gradient(
    90deg,
    rgba(255,255,255,0.03) 0%,
    rgba(255,255,255,0.07) 50%,
    rgba(255,255,255,0.03) 100%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

.info-placeholder {
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.line {
  height: 10px;
  border-radius: 4px;
  background: rgba(255,255,255,0.06);
  animation: shimmer 1.5s infinite;
}

.line.short { width: 60%; }
.line.long  { width: 90%; }

@keyframes shimmer {
  0%   { background-position: -200% 0; }
  100% { background-position: 200% 0; }
}
</style>
