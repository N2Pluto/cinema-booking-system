<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import AppNavbar from '@/components/layout/AppNavbar.vue'
import Pagination from '@/components/ui/Pagination.vue'

const router = useRouter()
const { fetchMe, token } = useAuth()
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

// ── Types ──────────────────────────────────────────────────────────────────
interface Booking {
  id: string
  user_id: string
  cinema_id: string
  seat_numbers: string[]
  total_amount: number
  status: string
  created_at: string
  movie_name: string
  theater_no: number
  start_time: string
}

interface BookingResult {
  data: Booking[]
  total: number
  page: number
  limit: number
  total_pages: number
}

// ── State ──────────────────────────────────────────────────────────────────
const loading = ref(false)
const result = ref<BookingResult>({ data: [], total: 0, page: 1, limit: 20, total_pages: 1 })

const filterMovie = ref('')
const filterDate = ref('')
const filterStatus = ref('')
const page = ref(1)

// ── Lifecycle ──────────────────────────────────────────────────────────────
onMounted(async () => {
  await fetchMe()
  await load()
})

// ── Methods ────────────────────────────────────────────────────────────────
async function load() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    params.set('page', String(page.value))
    params.set('limit', '20')
    if (filterMovie.value) params.set('movie', filterMovie.value)
    if (filterDate.value) params.set('date', filterDate.value)
    if (filterStatus.value) params.set('status', filterStatus.value)

    const res = await fetch(`${API_URL}/api/admin/bookings?${params}`, {
      headers: { Authorization: `Bearer ${token.value}` },
    })
    if (res.ok) result.value = await res.json()
  } finally {
    loading.value = false
  }
}

function applyFilter() {
  page.value = 1
  load()
}

function clearFilter() {
  filterMovie.value = ''
  filterDate.value = ''
  filterStatus.value = ''
  page.value = 1
  load()
}

async function changePage(p: number) {
  page.value = p
  await load()
}

// ── Helpers ────────────────────────────────────────────────────────────────
const statusStyles: Record<string, string> = {
  SUCCESS: 'text-[#64c896] bg-[#64c896]/10 border-[#64c896]/20',
  PENDING: 'text-yellow-400 bg-yellow-400/10 border-yellow-400/20',
  TIMEOUT: 'text-[#606882] bg-white/5 border-white/10',
  CANCELLED: 'text-[#e94560] bg-[#e94560]/10 border-[#e94560]/20',
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleString('th-TH', {
    month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit',
  })
}

function formatCurrency(amount: number) {
  return amount.toLocaleString('th-TH') + ' ฿'
}
</script>

<template>
  <div class="min-h-screen bg-[#0d0d1a] text-white">
    <AppNavbar />

    <main class="max-w-6xl mx-auto px-6 py-10">
      <!-- Back -->
      <button
        @click="router.back()"
        class="text-sm text-[#606882] hover:text-white mb-6 flex items-center gap-1 transition-colors"
      >
        ← กลับ
      </button>

      <!-- Header -->
      <div class="flex items-center justify-between mb-8">
        <div>
          <h2 class="text-xl font-bold">Booking Dashboard</h2>
          <p class="text-sm text-[#606882] mt-1">รายการจองทั้งหมดในระบบ</p>
        </div>
        <div class="text-sm text-[#606882]">
          ทั้งหมด <span class="text-white font-semibold">{{ result.total }}</span> รายการ
        </div>
      </div>

      <!-- Filter Bar -->
      <div class="bg-white/[0.03] border border-white/[0.08] rounded-2xl p-5 mb-6 flex flex-wrap gap-4 items-end">
        <!-- Movie filter -->
        <div class="flex flex-col gap-1.5 flex-1 min-w-[160px]">
          <label class="text-xs text-[#606882] font-medium">ชื่อหนัง</label>
          <input
            v-model="filterMovie"
            type="text"
            placeholder="ค้นหาชื่อหนัง..."
            class="px-3 py-2 rounded-lg bg-white/[0.06] border border-white/[0.08] text-sm text-white placeholder-[#606882] focus:outline-none focus:border-[#e94560]/50 transition-colors"
            @keyup.enter="applyFilter"
          />
        </div>

        <!-- Date filter -->
        <div class="flex flex-col gap-1.5 flex-1 min-w-[160px]">
          <label class="text-xs text-[#606882] font-medium">วันที่จอง</label>
          <input
            v-model="filterDate"
            type="date"
            class="px-3 py-2 rounded-lg bg-white/[0.06] border border-white/[0.08] text-sm text-white focus:outline-none focus:border-[#e94560]/50 transition-colors [color-scheme:dark]"
          />
        </div>

        <!-- Status filter -->
        <div class="flex flex-col gap-1.5 min-w-[140px]">
          <label class="text-xs text-[#606882] font-medium">สถานะ</label>
          <select
            v-model="filterStatus"
            class="px-3 py-2 rounded-lg bg-white/[0.06] border border-white/[0.08] text-sm text-white focus:outline-none focus:border-[#e94560]/50 transition-colors"
          >
            <option value="">ทั้งหมด</option>
            <option value="PENDING">PENDING</option>
            <option value="SUCCESS">SUCCESS</option>
            <option value="TIMEOUT">TIMEOUT</option>
            <option value="CANCELLED">CANCELLED</option>
          </select>
        </div>

        <!-- Actions -->
        <div class="flex gap-2 self-end">
          <button
            @click="applyFilter"
            class="px-4 py-2 text-sm font-medium bg-[#e94560] hover:bg-[#e94560]/80 text-white rounded-lg transition-colors"
          >
            ค้นหา
          </button>
          <button
            @click="clearFilter"
            class="px-4 py-2 text-sm font-medium bg-white/[0.06] hover:bg-white/[0.10] text-[#a8b2d8] rounded-lg transition-colors"
          >
            ล้าง
          </button>
        </div>
      </div>

      <!-- Table -->
      <div v-if="loading && result.data.length === 0" class="text-center py-20 text-[#606882]">
        กำลังโหลด...
      </div>

      <div v-else>
        <div class="bg-white/[0.03] border border-white/[0.08] rounded-2xl overflow-hidden mb-6">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b border-white/[0.08] text-[#606882] text-left">
                <th class="px-4 py-3 font-medium">หนัง / โรง</th>
                <th class="px-4 py-3 font-medium">ที่นั่ง</th>
                <th class="px-4 py-3 font-medium">ยอดรวม</th>
                <th class="px-4 py-3 font-medium">สถานะ</th>
                <th class="px-4 py-3 font-medium">วันที่จอง</th>
                <th class="px-4 py-3 font-medium text-[#606882]/60">Booking ID</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="booking in result.data"
                :key="booking.id"
                class="border-b border-white/[0.05] hover:bg-white/[0.02] transition-colors"
              >
                <!-- Movie / Theater -->
                <td class="px-4 py-3">
                  <div class="font-medium text-white">{{ booking.movie_name || '—' }}</div>
                  <div class="text-xs text-[#606882]">
                    โรง {{ booking.theater_no || '—' }}
                    <template v-if="booking.start_time">
                      · {{ formatDate(booking.start_time) }}
                    </template>
                  </div>
                </td>

                <!-- Seats -->
                <td class="px-4 py-3">
                  <div class="flex flex-wrap gap-1">
                    <span
                      v-for="seat in booking.seat_numbers"
                      :key="seat"
                      class="px-1.5 py-0.5 text-xs rounded bg-white/[0.08] text-[#a8b2d8] font-mono"
                    >
                      {{ seat }}
                    </span>
                  </div>
                </td>

                <!-- Amount -->
                <td class="px-4 py-3 font-semibold text-white">
                  {{ formatCurrency(booking.total_amount) }}
                </td>

                <!-- Status -->
                <td class="px-4 py-3">
                  <span
                    class="px-2 py-0.5 rounded-full text-xs font-semibold border"
                    :class="statusStyles[booking.status] || 'text-[#a8b2d8] bg-white/5 border-white/10'"
                  >
                    {{ booking.status }}
                  </span>
                </td>

                <!-- Created At -->
                <td class="px-4 py-3 text-[#606882] whitespace-nowrap">
                  {{ formatDate(booking.created_at) }}
                </td>

                <!-- Booking ID -->
                <td class="px-4 py-3 text-[#606882]/50 text-xs font-mono truncate max-w-[80px]">
                  {{ booking.id.slice(-8) }}
                </td>
              </tr>

              <tr v-if="result.data.length === 0">
                <td colspan="6" class="text-center py-14 text-[#606882]">
                  ไม่พบรายการจองที่ตรงกับเงื่อนไข
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <Pagination
          :page="result.page"
          :total-pages="result.total_pages"
          :total="result.total"
          :limit="result.limit"
          @change="changePage"
        />
      </div>
    </main>
  </div>
</template>
