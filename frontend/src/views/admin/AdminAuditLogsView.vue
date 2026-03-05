<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import AppNavbar from '@/components/layout/AppNavbar.vue'
import Pagination from '@/components/ui/Pagination.vue'

const router = useRouter()

interface AuditLog {
  id: string
  event_type: string
  user_id: string
  cinema_id: string
  details: string
  timestamp: string
}

interface LogsResult {
  data: AuditLog[]
  total: number
  page: number
  limit: number
  total_pages: number
}

const { fetchMe, token } = useAuth()
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const loading = ref(false)
const result = ref<LogsResult>({ data: [], total: 0, page: 1, limit: 20, total_pages: 1 })
const page = ref(1)

onMounted(async () => {
  await fetchMe()
  await load()
})

async function load() {
  loading.value = true
  try {
    const res = await fetch(`${API_URL}/api/admin/audit-logs?page=${page.value}&limit=20`, {
      headers: { Authorization: `Bearer ${token.value}` },
    })
    if (res.ok) result.value = await res.json()
  } finally {
    loading.value = false
  }
}

async function changePage(p: number) {
  page.value = p
  await load()
}

const eventColors: Record<string, string> = {
  BOOKING_SUCCESS: 'text-[#64c896] bg-[#64c896]/10 border-[#64c896]/20',
  SEAT_LOCKED: 'text-yellow-400 bg-yellow-400/10 border-yellow-400/20',
  SEAT_RELEASED: 'text-[#a8b2d8] bg-white/5 border-white/10',
  TIMEOUT: 'text-[#606882] bg-white/5 border-white/10',
  BOOKING_CANCELLED: 'text-[#e94560] bg-[#e94560]/10 border-[#e94560]/20',
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleString('th-TH', {
    month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit', second: '2-digit',
  })
}
</script>

<template>
  <div class="min-h-screen bg-[#0d0d1a] text-white">
    <AppNavbar />

    <main class="max-w-5xl mx-auto px-6 py-10">
      <button @click="router.back()" class="text-sm text-[#606882] hover:text-white mb-6 flex items-center gap-1 transition-colors">
        ← กลับ
      </button>

      <div class="flex items-center justify-between mb-8">
        <h2 class="text-xl font-bold">Audit Logs</h2>
      </div>

      <div v-if="loading && result.data.length === 0" class="text-center py-20 text-[#606882]">กำลังโหลด...</div>

      <div v-else>
        <div class="bg-white/[0.03] border border-white/[0.08] rounded-2xl overflow-hidden mb-6">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b border-white/[0.08] text-[#606882] text-left">
                <th class="px-4 py-3 font-medium">เวลา</th>
                <th class="px-4 py-3 font-medium">Event</th>
                <th class="px-4 py-3 font-medium">รายละเอียด</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="log in result.data"
                :key="log.id"
                class="border-b border-white/[0.05]"
              >
                <td class="px-4 py-3 text-[#606882] whitespace-nowrap">{{ formatDate(log.timestamp) }}</td>
                <td class="px-4 py-3">
                  <span
                    class="px-2 py-0.5 rounded-full text-xs font-semibold border"
                    :class="eventColors[log.event_type] || 'text-[#a8b2d8] bg-white/5 border-white/10'"
                  >
                    {{ log.event_type }}
                  </span>
                </td>
                <td class="px-4 py-3 text-[#a8b2d8] text-xs max-w-xs truncate">{{ log.details }}</td>
              </tr>
              <tr v-if="result.data.length === 0">
                <td colspan="3" class="text-center py-10 text-[#606882]">ยังไม่มีข้อมูล</td>
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
