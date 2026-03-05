<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { useBooking, type Booking } from '@/composables/useBooking'
import AppNavbar from '@/components/layout/AppNavbar.vue'

const router = useRouter()
const { fetchMe } = useAuth()
const { loading, fetchMyBookings } = useBooking()
const bookings = ref<Booking[]>([])

onMounted(async () => {
  await fetchMe()
  bookings.value = await fetchMyBookings()
})

const statusLabel: Record<string, string> = {
  PENDING: 'รอยืนยัน',
  SUCCESS: 'สำเร็จ',
  TIMEOUT: 'หมดเวลา',
  CANCELLED: 'ยกเลิก',
}
const statusClass: Record<string, string> = {
  PENDING: 'bg-yellow-500/15 text-yellow-400 border-yellow-500/30',
  SUCCESS: 'bg-[#64c896]/15 text-[#64c896] border-[#64c896]/30',
  TIMEOUT: 'bg-white/5 text-[#606882] border-white/10',
  CANCELLED: 'bg-white/5 text-[#606882] border-white/10',
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleString('th-TH', {
    month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit',
  })
}
</script>

<template>
  <div class="min-h-screen bg-[#0d0d1a] text-white">
    <AppNavbar />

    <main class="max-w-3xl mx-auto px-6 py-10">
      <button @click="router.back()" class="text-sm text-[#606882] hover:text-white mb-6 flex items-center gap-1 transition-colors">
        ← กลับ
      </button>

      <h2 class="text-xl font-bold mb-6">การจองของฉัน</h2>

      <div v-if="loading" class="text-center py-20 text-[#606882]">กำลังโหลด...</div>

      <div v-else-if="bookings.length === 0" class="text-center py-20 text-[#606882]">
        ยังไม่มีประวัติการจอง
      </div>

      <div v-else class="flex flex-col gap-4">
        <div
          v-for="b in bookings"
          :key="b.id"
          class="bg-white/[0.04] border border-white/[0.08] rounded-2xl p-5 flex items-start justify-between gap-4"
        >
          <div class="flex-1 min-w-0">
            <p class="text-sm font-semibold text-white mb-1 font-mono truncate">{{ b.id }}</p>
            <p class="text-sm text-[#a8b2d8] mb-1">
              ที่นั่ง: <span class="text-white">{{ b.seat_numbers.join(', ') }}</span>
            </p>
            <p class="text-xs text-[#606882]">{{ formatDate(b.created_at) }}</p>
          </div>
          <div class="flex flex-col items-end gap-2 shrink-0">
            <span
              class="px-2.5 py-0.5 rounded-full text-xs font-semibold border"
              :class="statusClass[b.status] || statusClass.CANCELLED"
            >
              {{ statusLabel[b.status] || b.status }}
            </span>
            <span class="text-[#e94560] font-bold text-sm">฿{{ b.total_amount.toLocaleString() }}</span>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
