<script setup lang="ts">
import { computed } from 'vue'
import type { Seat } from '@/composables/useCinemas'
import SeatCell from './SeatCell.vue'

const props = defineProps<{
  seats: Seat[]
  selectedSeats: string[]
  myUserId?: string
}>()

const emit = defineEmits<{ (e: 'toggle', seatNo: string): void }>()

// Group seats by row letter
const rows = computed(() => {
  const map = new Map<string, Seat[]>()
  for (const seat of props.seats) {
    const row = seat.seat_no.replace(/\d+$/, '')
    if (!map.has(row)) map.set(row, [])
    map.get(row)!.push(seat)
  }
  return Array.from(map.entries()).map(([row, seats]) => ({ row, seats }))
})
</script>

<template>
  <div class="select-none">
    <!-- Screen indicator -->
    <div class="flex flex-col items-center mb-6">
      <div class="w-48 h-2 bg-gradient-to-r from-transparent via-[#e94560]/60 to-transparent rounded-full mb-1"></div>
      <p class="text-xs text-[#606882]">จอภาพยนตร์</p>
    </div>

    <!-- Seat grid -->
    <div class="flex flex-col gap-2">
      <div v-for="{ row, seats: rowSeats } in rows" :key="row" class="flex items-center gap-2">
        <span class="w-5 text-xs text-[#606882] text-right shrink-0">{{ row }}</span>
        <div class="flex gap-1.5 flex-wrap">
          <SeatCell
            v-for="seat in rowSeats"
            :key="seat.seat_no"
            :seat="seat"
            :selected="selectedSeats.includes(seat.seat_no)"
            :my-user-id="myUserId"
            @toggle="emit('toggle', $event)"
          />
        </div>
      </div>
    </div>

    <!-- Legend -->
    <div class="flex items-center gap-4 mt-6 justify-center text-xs text-[#606882]">
      <div class="flex items-center gap-1.5">
        <div class="w-4 h-4 rounded bg-white/[0.06] border border-white/[0.15]"></div>ว่าง
      </div>
      <div class="flex items-center gap-1.5">
        <div class="w-4 h-4 rounded bg-[#e94560] border border-[#e94560]"></div>เลือกแล้ว
      </div>
      <div class="flex items-center gap-1.5">
        <div class="w-4 h-4 rounded bg-yellow-500/20 border border-yellow-500/30"></div>ถูกจองชั่วคราว
      </div>
      <div class="flex items-center gap-1.5">
        <div class="w-4 h-4 rounded bg-white/[0.03] border border-white/[0.05]"></div>จองแล้ว
      </div>
    </div>
  </div>
</template>
