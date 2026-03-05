<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import RoleBadge from '@/components/ui/RoleBadge.vue'

const router = useRouter()
const { user, logout } = useAuth()

function handleLogout() {
  logout()
  router.push({ name: 'login' })
}
</script>

<template>
  <nav class="sticky top-0 z-10 flex items-center justify-between px-8 py-4 bg-white/[0.03] border-b border-white/[0.08] backdrop-blur-md">
    <!-- Brand -->
    <button
      class="flex items-center gap-2.5 text-lg font-bold tracking-tight text-white hover:opacity-80 transition"
      @click="router.push({ name: 'home' })"
      style="outline: none; cursor: pointer;"
    >
      <span class="w-2.5 h-2.5 rounded-full bg-[#e94560]"></span>
      Cinema Booking
    </button>

    <!-- User info -->
    <div v-if="user" class="flex items-center gap-3">
      <div class="w-8 h-8 rounded-full bg-gradient-to-br from-[#e94560] to-[#302b63] flex items-center justify-center text-sm font-bold shrink-0">
        {{ user.display_name?.charAt(0).toUpperCase() }}
      </div>
      <div class="flex flex-col leading-tight">
        <span class="text-sm font-semibold text-white">{{ user.display_name }}</span>
        <span class="text-xs text-[#606882]">{{ user.email }}</span>
      </div>
      <RoleBadge :role="user.role" />
      <button
        @click="handleLogout"
        class="px-4 py-1.5 text-sm font-medium text-[#e94560] bg-[#e94560]/10 border border-[#e94560]/30 rounded-lg hover:bg-[#e94560]/20 transition-colors cursor-pointer"
      >
        ออกจากระบบ
      </button>
    </div>

    <button
      v-else
      @click="handleLogout"
      class="px-4 py-1.5 text-sm font-medium text-[#e94560] bg-[#e94560]/10 border border-[#e94560]/30 rounded-lg hover:bg-[#e94560]/20 transition-colors cursor-pointer"
    >
      ออกจากระบบ
    </button>
  </nav>
</template>
