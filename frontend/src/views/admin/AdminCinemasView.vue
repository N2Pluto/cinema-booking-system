<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { useCinemas, type Cinema } from '@/composables/useCinemas'
import AppNavbar from '@/components/layout/AppNavbar.vue'
import Pagination from '@/components/ui/Pagination.vue'

const router = useRouter()
const { fetchMe } = useAuth()
const { loading, result, fetchCinemas } = useCinemas()

const page = ref(1)
let pollTimer: ReturnType<typeof setInterval> | null = null

onMounted(async () => {
  await fetchMe()
  await load()
  // Poll every 10s
  pollTimer = setInterval(load, 10_000)
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
})

async function load() {
  await fetchCinemas({ page: page.value, limit: 10, admin: true, orderBy: 'desc' })
}

async function changePage(p: number) {
  page.value = p
  await load()
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleString('th-TH', {
    month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit',
  })
}

function availableCount(c: Cinema) {
  return c.seats?.filter(s => s.status === 'AVAILABLE').length ?? 0
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
        <div>
          <h2 class="text-xl font-bold">จัดการรอบหนัง</h2>
          <p class="text-sm text-[#606882] mt-1">รีเฟรชอัตโนมัติทุก 10 วินาที</p>
        </div>
        <div class="flex gap-3">
          <router-link
            to="/admin/audit-logs"
            class="px-4 py-2 text-sm border border-white/[0.1] text-[#a8b2d8] rounded-lg hover:bg-white/5 transition-colors"
          >
            Audit Logs
          </router-link>
          <router-link
            to="/admin/cinemas/create"
            class="px-4 py-2 text-sm bg-[#e94560] text-white rounded-lg hover:bg-[#d63651] transition-colors font-medium"
          >
            + เพิ่มรอบหนัง
          </router-link>
        </div>
      </div>

      <div v-if="loading && !result.data?.length" class="text-center py-20 text-[#606882]">กำลังโหลด...</div>

      <div v-else>
        <div class="bg-white/[0.03] border border-white/[0.08] rounded-2xl overflow-hidden mb-6">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b border-white/[0.08] text-[#606882] text-left">
                <th class="px-4 py-3 font-medium">ชื่อหนัง</th>
                <th class="px-4 py-3 font-medium">โรง</th>
                <th class="px-4 py-3 font-medium">เริ่ม</th>
                <th class="px-4 py-3 font-medium">สิ้นสุด</th>
                <th class="px-4 py-3 font-medium text-right">ที่ว่าง / ทั้งหมด</th>
                <th class="px-4 py-3 font-medium text-right">ราคา</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="c in result.data"
                :key="c.id"
                class="border-b border-white/[0.05] hover:bg-white/[0.03] cursor-pointer transition-colors"
                @click="router.push(`/cinema/${c.id}`)"
              >
                <td class="px-4 py-3 font-medium text-white">{{ c.movie_name }}</td>
                <td class="px-4 py-3 text-[#a8b2d8]">{{ c.theater_no }}</td>
                <td class="px-4 py-3 text-[#a8b2d8]">{{ formatDate(c.start_time) }}</td>
                <td class="px-4 py-3 text-[#a8b2d8]">{{ formatDate(c.end_time) }}</td>
                <td class="px-4 py-3 text-right">
                  <span :class="availableCount(c) > 0 ? 'text-[#64c896]' : 'text-[#e94560]'">
                    {{ availableCount(c) }}
                  </span>
                  <span class="text-[#606882]"> / {{ c.seats?.length ?? 0 }}</span>
                </td>
                <td class="px-4 py-3 text-right text-[#e94560] font-semibold">฿{{ c.price }}</td>
              </tr>
              <tr v-if="result.data.length === 0">
                <td colspan="6" class="text-center py-10 text-[#606882]">ยังไม่มีรอบหนัง</td>
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
