<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { useCinemas } from '@/composables/useCinemas'
import AppNavbar from '@/components/layout/AppNavbar.vue'
import CinemaCard from '@/components/cinema/CinemaCard.vue'
import MovieCardSkeleton from '@/components/ui/MovieCardSkeleton.vue'
import Pagination from '@/components/ui/Pagination.vue'

const router = useRouter()
const { fetchMe, user } = useAuth()
const { loading, result, fetchCinemas } = useCinemas()
const page = ref(1)

onMounted(async () => {
  await fetchMe()
  await fetchCinemas({ page: page.value, limit: 12 })
})

async function changePage(p: number) {
  page.value = p
  await fetchCinemas({ page: p, limit: 12 })
  window.scrollTo({ top: 0, behavior: 'smooth' })
}
</script>

<template>
  <div class="min-h-screen bg-[#0d0d1a] text-white">
    <AppNavbar />

    <main class="max-w-[1100px] mx-auto px-6 py-10">
      <!-- Admin shortcut -->
      <div v-if="user?.is_admin" class="mb-6 flex gap-3">
        <router-link
          to="/admin/cinemas"
          class="px-4 py-2 text-sm font-medium bg-[#e94560]/10 text-[#e94560] border border-[#e94560]/30 rounded-lg hover:bg-[#e94560]/20 transition-colors"
        >
          จัดการรอบหนัง
        </router-link>
        <router-link
          to="/admin/audit-logs"
          class="px-4 py-2 text-sm font-medium bg-white/5 text-[#a8b2d8] border border-white/10 rounded-lg hover:bg-white/10 transition-colors"
        >
          Audit Logs
        </router-link>
      </div>

      <div class="flex items-center justify-between mb-6">
        <div>
          <h3 class="text-xl font-semibold mb-1">ภาพยนตร์ที่กำลังฉาย</h3>
          <p class="text-sm text-[#606882]">เลือกรอบที่ต้องการแล้วจองได้เลย</p>
        </div>
        <router-link
          to="/bookings"
          class="text-sm text-[#a8b2d8] hover:text-white transition-colors"
        >
          การจองของฉัน →
        </router-link>
      </div>

      <!-- Skeleton -->
      <div v-if="loading" class="grid gap-5" style="grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));">
        <MovieCardSkeleton v-for="n in 12" :key="n" />
      </div>

      <!-- Cinema cards -->
      <template v-else>
        <div v-if="!result.data?.length" class="text-center py-20 text-[#606882]">
          ยังไม่มีรอบหนัง
        </div>
        <div v-else class="grid gap-5 mb-8" style="grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));">
          <router-link
            v-for="cinema in result.data"
            :key="cinema.id"
            :to="`/cinema/${cinema.id}`"
          >
            <CinemaCard :cinema="cinema" />
          </router-link>
        </div>

        <Pagination
          :page="result.page"
          :total-pages="result.total_pages"
          :total="result.total"
          :limit="result.limit"
          @change="changePage"
        />
      </template>
    </main>
  </div>
</template>
