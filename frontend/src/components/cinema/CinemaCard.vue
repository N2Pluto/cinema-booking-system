<script setup lang="ts">
import { computed } from 'vue'
import type { Cinema } from '@/composables/useCinemas'

const props = defineProps<{ cinema: Cinema }>()

const availableSeats = computed(() =>
  props.cinema.seats?.filter(s => s.status === 'AVAILABLE').length ?? 0
)
const totalSeats = computed(() => props.cinema.seats?.length ?? 0)

function formatDate(iso: string) {
  return new Date(iso).toLocaleString('th-TH', {
    month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit',
  })
}
</script>

<template>
  <div class="bg-white/[0.04] border border-white/[0.08] rounded-2xl overflow-hidden hover:border-[#e94560]/40 hover:bg-white/[0.07] transition-all cursor-pointer group">
    <!-- Poster placeholder -->
    <div class="w-full aspect-[2/3] bg-gradient-to-br from-[#302b63] to-[#0f0c29] flex items-end p-3 relative">
      <span class="absolute top-3 right-3 text-xs font-semibold px-2 py-0.5 rounded-full"
        :class="availableSeats > 0
          ? 'bg-[#64c896]/20 text-[#64c896] border border-[#64c896]/30'
          : 'bg-[#e94560]/20 text-[#e94560] border border-[#e94560]/30'">
        {{ availableSeats > 0 ? `${availableSeats} ที่ว่าง` : 'เต็มแล้ว' }}
      </span>
      <div class="w-10 h-10 rounded-full bg-[#e94560]/20 flex items-center justify-center text-[#e94560] text-xl font-bold">
        {{ cinema.theater_no }}
      </div>
    </div>

    <!-- Info -->
    <div class="p-3">
      <p class="text-sm font-semibold text-white truncate mb-1">{{ cinema.movie_name }}</p>
      <p class="text-xs text-[#606882] mb-2">{{ formatDate(cinema.start_time) }}</p>
      <div class="flex items-center justify-between">
        <span class="text-xs text-[#a8b2d8]">{{ totalSeats }} ที่นั่ง</span>
        <span class="text-xs font-semibold text-[#e94560]">฿{{ cinema.price }}</span>
      </div>
    </div>
  </div>
</template>
