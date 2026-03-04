<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { useCinemas, type Cinema, type Seat } from '@/composables/useCinemas'
import { useBooking } from '@/composables/useBooking'
import AppNavbar from '@/components/layout/AppNavbar.vue'
import SeatGrid from '@/components/cinema/SeatGrid.vue'

const route = useRoute()
const router = useRouter()
const { token, user, fetchMe } = useAuth()
const { fetchCinema } = useCinemas()
const { loading: bookLoading, error: bookError, lockSeats, confirmBooking } = useBooking()

const cinema = ref<Cinema | null>(null)
const seats = ref<Seat[]>([])
const selectedSeats = ref<string[]>([])
const pendingBookingId = ref<string | null>(null)
const step = ref<'select' | 'confirm' | 'done'>('select')
const loadingCinema = ref(true)
let ws: WebSocket | null = null

const cinemaId = route.params.id as string

const totalPrice = computed(() =>
  selectedSeats.value.reduce((sum, sno) => {
    const seat = seats.value.find(s => s.seat_no === sno)
    return sum + (seat?.price ?? 0)
  }, 0)
)

onMounted(async () => {
  await fetchMe()
  await loadCinema()
  connectWS()
})

onUnmounted(() => {
  ws?.close()
})

async function loadCinema() {
  loadingCinema.value = true
  const data = await fetchCinema(cinemaId)
  if (data) {
    cinema.value = data as unknown as Cinema
    seats.value = (data as any).seats ?? []
  }
  loadingCinema.value = false
}

function connectWS() {
  const wsBase = (import.meta.env.VITE_API_URL || 'http://localhost:8080')
    .replace(/^http/, 'ws')
  ws = new WebSocket(`${wsBase}/api/ws/seats/${cinemaId}?token=${token.value}`)
  ws.onmessage = (ev) => {
    try {
      const msg = JSON.parse(ev.data)
      if (msg.type === 'seat_update' && Array.isArray(msg.seats)) {
        seats.value = msg.seats
        // Deselect seats that are no longer available
        selectedSeats.value = selectedSeats.value.filter(sno => {
          const s = msg.seats.find((x: Seat) => x.seat_no === sno)
          return s?.status === 'AVAILABLE'
        })
      }
    } catch {}
  }
  ws.onclose = () => {
    // Reconnect after 3s
    setTimeout(connectWS, 3000)
  }
}

function toggleSeat(seatNo: string) {
  const idx = selectedSeats.value.indexOf(seatNo)
  if (idx >= 0) {
    selectedSeats.value.splice(idx, 1)
  } else {
    selectedSeats.value.push(seatNo)
  }
}

async function handleLock() {
  if (selectedSeats.value.length === 0) return
  const booking = await lockSeats(cinemaId, selectedSeats.value)
  if (booking) {
    pendingBookingId.value = booking.id
    step.value = 'confirm'
  }
}

async function handleConfirm() {
  if (!pendingBookingId.value) return
  const booking = await confirmBooking(pendingBookingId.value)
  if (booking) {
    step.value = 'done'
  }
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleString('th-TH', {
    weekday: 'short', month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit',
  })
}
</script>

<template>
  <div class="min-h-screen bg-[#0d0d1a] text-white">
    <AppNavbar />

    <main class="max-w-4xl mx-auto px-6 py-10">
      <button @click="router.back()" class="text-sm text-[#606882] hover:text-white mb-6 flex items-center gap-1 transition-colors">
        ← กลับ
      </button>

      <div v-if="loadingCinema" class="text-center py-20 text-[#606882]">กำลังโหลด...</div>

      <template v-else-if="cinema">
        <!-- Header -->
        <div class="mb-8">
          <h1 class="text-2xl font-bold mb-1">{{ (cinema as any).movie_name }}</h1>
          <div class="flex flex-wrap gap-4 text-sm text-[#606882]">
            <span>โรง {{ (cinema as any).theater_no }}</span>
            <span>{{ formatDate((cinema as any).start_time) }} – {{ formatDate((cinema as any).end_time) }}</span>
            <span>฿{{ (cinema as any).price }} / ที่นั่ง</span>
          </div>
        </div>

        <!-- Step: select seats -->
        <div v-if="step === 'select'">
          <SeatGrid
            :seats="seats"
            :selected-seats="selectedSeats"
            :my-user-id="user?.user_id"
            @toggle="toggleSeat"
          />

          <div v-if="selectedSeats.length > 0" class="mt-8 p-4 bg-white/[0.04] border border-white/[0.1] rounded-2xl flex items-center justify-between gap-4">
            <div>
              <p class="text-sm text-[#a8b2d8]">ที่นั่งที่เลือก: <span class="text-white font-semibold">{{ selectedSeats.join(', ') }}</span></p>
              <p class="text-lg font-bold text-[#e94560] mt-1">รวม ฿{{ totalPrice.toLocaleString() }}</p>
            </div>
            <button
              @click="handleLock"
              :disabled="bookLoading"
              class="px-6 py-2.5 bg-[#e94560] text-white font-semibold rounded-xl hover:bg-[#d63651] transition-colors disabled:opacity-50 shrink-0"
            >
              {{ bookLoading ? 'กำลังจอง...' : 'จองที่นั่ง' }}
            </button>
          </div>
          <p v-if="bookError" class="mt-3 text-sm text-[#e94560]">{{ bookError }}</p>
        </div>

        <!-- Step: confirm payment -->
        <div v-else-if="step === 'confirm'" class="max-w-md mx-auto">
          <div class="bg-white/[0.04] border border-white/[0.1] rounded-2xl p-6 text-center">
            <div class="w-14 h-14 rounded-full bg-[#e94560]/15 flex items-center justify-center mx-auto mb-4 text-2xl">🎬</div>
            <h2 class="text-xl font-bold mb-2">ยืนยันการจอง</h2>
            <p class="text-[#a8b2d8] text-sm mb-4">ที่นั่ง: <strong class="text-white">{{ selectedSeats.join(', ') }}</strong></p>
            <p class="text-2xl font-bold text-[#e94560] mb-6">฿{{ totalPrice.toLocaleString() }}</p>
            <p class="text-xs text-[#606882] mb-6">ที่นั่งถูกล็อกไว้ 5 นาที กรุณายืนยันก่อนหมดเวลา</p>
            <div class="flex gap-3">
              <button @click="step = 'select'" class="flex-1 py-2.5 rounded-xl border border-white/20 text-[#a8b2d8] hover:bg-white/5 transition-colors text-sm">
                ยกเลิก
              </button>
              <button
                @click="handleConfirm"
                :disabled="bookLoading"
                class="flex-1 py-2.5 rounded-xl bg-[#e94560] text-white font-semibold hover:bg-[#d63651] transition-colors text-sm disabled:opacity-50"
              >
                {{ bookLoading ? 'กำลังดำเนินการ...' : 'ชำระเงิน' }}
              </button>
            </div>
            <p v-if="bookError" class="mt-3 text-sm text-[#e94560]">{{ bookError }}</p>
          </div>
        </div>

        <!-- Step: done -->
        <div v-else-if="step === 'done'" class="max-w-md mx-auto text-center">
          <div class="bg-white/[0.04] border border-[#64c896]/30 rounded-2xl p-8">
            <div class="text-5xl mb-4">🎉</div>
            <h2 class="text-2xl font-bold mb-2 text-[#64c896]">จองสำเร็จ!</h2>
            <p class="text-[#a8b2d8] mb-6">ที่นั่ง {{ selectedSeats.join(', ') }} ถูกจองเรียบร้อยแล้ว</p>
            <div class="flex gap-3">
              <button @click="router.push('/')" class="flex-1 py-2.5 rounded-xl border border-white/20 text-[#a8b2d8] hover:bg-white/5 transition-colors text-sm">
                กลับหน้าแรก
              </button>
              <router-link
                to="/bookings"
                class="flex-1 py-2.5 rounded-xl bg-[#e94560] text-white font-semibold text-sm text-center hover:bg-[#d63651] transition-colors"
              >
                ดูการจองของฉัน
              </router-link>
            </div>
          </div>
        </div>
      </template>
    </main>
  </div>
</template>
