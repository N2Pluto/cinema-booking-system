<script setup lang="ts">
import type { Seat } from '@/composables/useCinemas'

const props = defineProps<{
  seat: Seat
  selected: boolean
  myUserId?: string
}>()

const emit = defineEmits<{ (e: 'toggle', seatNo: string): void }>()

function handleClick() {
  if (props.seat.status !== 'AVAILABLE' && !props.selected) return
  emit('toggle', props.seat.seat_no)
}
</script>

<template>
  <button
    :title="`${seat.seat_no} • ฿${seat.price}`"
    :disabled="seat.status === 'BOOKED' || seat.status === 'LOCKED'"
    @click="handleClick"
    class="w-8 h-8 rounded-md text-xs font-semibold transition-all border"
    :class="{
      // Available + not selected
      'bg-white/[0.06] border-white/[0.15] text-[#a8b2d8] hover:bg-[#e94560]/20 hover:border-[#e94560]/50 hover:text-white cursor-pointer':
        seat.status === 'AVAILABLE' && !selected,
      // Selected
      'bg-[#e94560] border-[#e94560] text-white cursor-pointer shadow-[0_0_10px_rgba(233,69,96,0.4)]':
        selected,
      // Locked by someone else
      'bg-yellow-500/20 border-yellow-500/30 text-yellow-400 cursor-not-allowed':
        seat.status === 'LOCKED' && !selected,
      // Booked
      'bg-white/[0.03] border-white/[0.05] text-white/20 cursor-not-allowed':
        seat.status === 'BOOKED',
    }"
  >
    {{ seat.seat_no }}
  </button>
</template>
