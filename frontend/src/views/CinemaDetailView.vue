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
const { loading: bookLoading, error: bookError, lockSeats, confirmBooking, cancelBooking, fetchMyBookings } = useBooking()

const cinema = ref<Cinema | null>(null)
const seats = ref<Seat[]>([])
const selectedSeats = ref<string[]>([])
const pendingBookingId = ref<string | null>(null)
const step = ref<'select' | 'qr' | 'done'>('select')
const loadingCinema = ref(true)
const countdown = ref(300) // 5 minutes
let ws: WebSocket | null = null
let countdownTimer: ReturnType<typeof setInterval> | null = null

const cinemaId = route.params.id as string

const totalPrice = computed(() =>
  selectedSeats.value.reduce((sum, sno) => {
    const seat = seats.value.find(s => s.seat_no === sno)
    return sum + (seat?.price ?? 0)
  }, 0)
)

const countdownDisplay = computed(() => {
  const m = Math.floor(countdown.value / 60).toString().padStart(2, '0')
  const s = (countdown.value % 60).toString().padStart(2, '0')
  return `${m}:${s}`
})

// Generate a pseudo-random but deterministic 21×21 QR-looking grid
const qrModules = computed(() => {
  const size = 21
  const grid: boolean[][] = Array.from({ length: size }, () => Array(size).fill(false))

  // Finder pattern: 7×7 at (r,c)
  const setFinder = (r: number, c: number) => {
    for (let i = 0; i < 7; i++)
      for (let j = 0; j < 7; j++) {
        const outer = i === 0 || i === 6 || j === 0 || j === 6
        const inner = i >= 2 && i <= 4 && j >= 2 && j <= 4
        grid[r + i]![c + j] = outer || inner
      }
  }
  setFinder(0, 0)
  setFinder(0, 14)
  setFinder(14, 0)

  // Timing strips (row 6 / col 6)
  for (let i = 8; i < 13; i++) {
    grid[6]![i] = i % 2 === 0
    grid[i]![6] = i % 2 === 0
  }

  // Data modules: deterministic pseudo-random from booking ID
  const seed = pendingBookingId.value ?? 'demo'
  let hash = 5381
  for (let i = 0; i < seed.length; i++) hash = ((hash << 5) + hash) + seed.charCodeAt(i)

  for (let r = 0; r < size; r++) {
    for (let c = 0; c < size; c++) {
      const inFinderTL = r < 8 && c < 8
      const inFinderTR = r < 8 && c > 12
      const inFinderBL = r > 12 && c < 8
      if (inFinderTL || inFinderTR || inFinderBL || r === 6 || c === 6) continue
      const v = Math.imul(hash, r * 31 + c + 1)
      grid[r]![c] = (v >>> 0) % 3 !== 0
    }
  }
  return grid
})

onMounted(async () => {
  await fetchMe()
  await loadCinema()
  connectWS()
  await restorePendingBooking()
})

onUnmounted(() => {
  ws?.close()
  if (countdownTimer) clearInterval(countdownTimer)
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
  const wsBase = (import.meta.env.VITE_API_URL || 'http://localhost:8080').replace(/^http/, 'ws')
  ws = new WebSocket(`${wsBase}/api/ws/seats/${cinemaId}?token=${token.value}`)
  ws.onmessage = (ev) => {
    try {
      const msg = JSON.parse(ev.data)
      if (msg.type === 'seat_update' && Array.isArray(msg.seats)) {
        seats.value = msg.seats
        selectedSeats.value = selectedSeats.value.filter(sno => {
          const s = msg.seats.find((x: Seat) => x.seat_no === sno)
          return s?.status === 'AVAILABLE'
        })
      }
    } catch {}
  }
  ws.onclose = () => setTimeout(connectWS, 3000)
}

function toggleSeat(seatNo: string) {
  const idx = selectedSeats.value.indexOf(seatNo)
  if (idx >= 0) selectedSeats.value.splice(idx, 1)
  else selectedSeats.value.push(seatNo)
}

const LOCK_TTL_SEC = 300

function startCountdown(remainingSeconds?: number) {
  countdown.value = remainingSeconds ?? LOCK_TTL_SEC
  if (countdownTimer) clearInterval(countdownTimer)
  countdownTimer = setInterval(async () => {
    countdown.value--
    if (countdown.value <= 0) {
      clearInterval(countdownTimer!)
      pendingBookingId.value = null
      selectedSeats.value = []
      step.value = 'select'
    }
  }, 1000)
}

/** หลังรีเฟรช: ถ้ามีการจอง PENDING ของรอบนี้อยู่ ให้กลับมาแสดงหน้า QR + countdown ที่เหลือ */
async function restorePendingBooking() {
  if (!user.value) return
  const list = await fetchMyBookings()
  const created = (b: { cinema_id: string; status: string; created_at: string }) =>
    b.cinema_id === cinemaId && b.status === 'PENDING' ? new Date(b.created_at).getTime() : 0
  const pending = list
    .filter((b) => b.cinema_id === cinemaId && b.status === 'PENDING')
    .sort((a, b) => created(b) - created(a))[0]
  if (!pending) return
  pendingBookingId.value = String(pending.id)
  selectedSeats.value = [...(pending.seat_numbers || [])]
  step.value = 'qr'
  const elapsed = (Date.now() - new Date(pending.created_at).getTime()) / 1000
  const remaining = Math.max(0, Math.floor(LOCK_TTL_SEC - elapsed))
  startCountdown(remaining)
}

async function handleLock() {
  if (selectedSeats.value.length === 0) return
  const booking = await lockSeats(cinemaId, selectedSeats.value)
  if (booking) {
    pendingBookingId.value = booking.id
    startCountdown()
    step.value = 'qr'
  }
}

async function handleConfirm() {
  if (!pendingBookingId.value) return
  const booking = await confirmBooking(pendingBookingId.value)
  if (booking) {
    if (countdownTimer) clearInterval(countdownTimer)
    step.value = 'done'
  }
}

async function handleCancel() {
  if (countdownTimer) clearInterval(countdownTimer)
  const bookingIdToCancel = pendingBookingId.value

  if (bookingIdToCancel) {
    const ok = await cancelBooking(bookingIdToCancel)
    if (!ok) {
      bookError.value = 'ยกเลิกไม่สำเร็จ กรุณาลองใหม่หรือรอปล่อยอัตโนมัติ 5 นาที'
      return
    }
  }

  pendingBookingId.value = null
  selectedSeats.value = []
  if (!bookingIdToCancel) bookError.value = null

  await loadCinema()
}

async function goBackToSelect() {
  step.value = 'select'
  bookError.value = null
  await loadCinema()
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

        <!-- ─── Step: select seats ─────────────────────────────────────── -->
        <div v-if="step === 'select'">
          <SeatGrid
            :seats="seats"
            :selected-seats="selectedSeats"
            :my-user-id="user?.user_id"
            :readonly="user?.role === 'ADMIN'"
            @toggle="toggleSeat"
          />

          <!-- Admin: view-only notice -->
          <div v-if="user?.role === 'ADMIN'" class="mt-6 p-3 bg-amber-500/10 border border-amber-500/30 rounded-xl text-center">
            <p class="text-amber-400 text-sm">👁 Admin สามารถดูแผนที่นั่งได้เท่านั้น ไม่สามารถจองได้</p>
          </div>

          <!-- User: booking panel -->
          <div v-else-if="selectedSeats.length > 0" class="mt-8 p-4 bg-white/[0.04] border border-white/[0.1] rounded-2xl flex items-center justify-between gap-4">
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

        <!-- ─── Step: QR payment ──────────────────────────────────────── -->
        <div v-else-if="step === 'qr'" class="max-w-sm mx-auto">
          <!-- หลังรีเซ็ต: ยังอยู่หน้า QR แสดงข้อความยกเลิก + ปุ่มกลับ (qr ต้องไม่หาย) -->
          <div v-if="!pendingBookingId" class="bg-white rounded-2xl overflow-hidden shadow-xl shadow-black/40 p-8 text-center">
            <p class="text-[#1a237e] font-semibold mb-2">ยกเลิกการจองแล้ว</p>
            <p class="text-[#606882] text-sm mb-6">ที่นั่งถูกปล่อยแล้ว สามารถเลือกจองใหม่ได้</p>
            <button
              @click="goBackToSelect"
              class="w-full py-3 rounded-xl bg-[#e94560] text-white font-semibold hover:bg-[#d63651] transition-colors text-sm"
            >
              เลือกที่นั่งใหม่
            </button>
          </div>
          <!-- กำลังจองอยู่: แสดง QR + ปุ่มยกเลิก/ยืนยัน -->
          <template v-else>
          <div class="bg-white rounded-2xl overflow-hidden shadow-xl shadow-black/40">
            <!-- Header bar -->
            <div class="bg-[#1a237e] px-5 py-3 flex items-center justify-between">
              <div class="flex items-center gap-2">
                <div class="w-6 h-6 bg-white rounded-full flex items-center justify-center">
                  <span class="text-[#1a237e] text-xs font-black">P</span>
                </div>
                <span class="text-white font-semibold text-sm">PromptPay</span>
              </div>
              <span
                class="text-xs font-mono font-semibold px-2 py-0.5 rounded-full"
                :class="countdown <= 60 ? 'bg-red-500/30 text-red-200 animate-pulse' : 'bg-white/20 text-white'"
              >
                ⏱ {{ countdownDisplay }}
              </span>
            </div>

            <!-- QR code -->
            <div class="px-6 pt-5 pb-3 flex flex-col items-center">
              <p class="text-[#1a237e] text-xs font-medium mb-3 uppercase tracking-wide">สแกนเพื่อชำระเงิน</p>

              <!-- Mock QR SVG -->
              <div class="p-2 border-4 border-[#1a237e] rounded-xl bg-white">
                <svg width="180" height="180" :viewBox="`0 0 ${21*9+4} ${21*9+4}`">
                  <rect width="100%" height="100%" fill="white" />
                  <template v-for="(row, r) in qrModules" :key="r">
                    <template v-for="(dark, c) in row" :key="`${r}-${c}`">
                      <rect
                        v-if="dark"
                        :x="c * 9 + 2"
                        :y="r * 9 + 2"
                        width="8"
                        height="8"
                        rx="1.5"
                        fill="#0d0d2a"
                      />
                    </template>
                  </template>
                </svg>
              </div>

              <!-- Amount -->
              <div class="mt-4 text-center">
                <p class="text-[#606882] text-xs">ยอดชำระ</p>
                <p class="text-[#1a237e] text-2xl font-black mt-0.5">฿{{ totalPrice.toLocaleString() }}</p>
                <p class="text-[#606882] text-xs mt-1">ที่นั่ง: <span class="text-[#1a237e] font-semibold">{{ selectedSeats.join(', ') }}</span></p>
              </div>

              <!-- PromptPay logo placeholder -->
              <div class="mt-3 flex items-center gap-1 text-[#606882] text-[10px]">
                <div class="flex gap-0.5">
                  <div class="w-3 h-3 rounded-full bg-[#e53935]"></div>
                  <div class="w-3 h-3 rounded-full bg-[#fb8c00] -ml-1"></div>
                </div>
                <span>ธนาคารไทยพาณิชย์ • DEMO</span>
              </div>
            </div>

            <!-- Note -->
            <div class="bg-amber-50 px-5 py-2.5 text-center">
              <p class="text-amber-700 text-xs">⚠️ QR Code นี้เป็นตัวอย่างสำหรับการ demo เท่านั้น</p>
            </div>
          </div>

          <!-- Action buttons (outside white card) -->
          <div class="mt-4 flex gap-3">
            <button
              @click="handleCancel"
              :disabled="bookLoading"
              class="flex-1 py-3 rounded-xl border border-white/20 text-[#a8b2d8] hover:bg-white/5 transition-colors text-sm disabled:opacity-50"
            >
              {{ bookLoading ? 'กำลังยกเลิก...' : 'ยกเลิก' }}
            </button>
            <button
              @click="handleConfirm"
              :disabled="bookLoading"
              class="flex-1 py-3 rounded-xl bg-[#e94560] text-white font-semibold hover:bg-[#d63651] transition-colors text-sm disabled:opacity-50"
            >
              {{ bookLoading ? 'กำลังยืนยัน...' : '✓ ยืนยัน QR ชำระเงินแล้ว' }}
            </button>
          </div>
          </template>
          <p v-if="bookError" class="mt-3 text-sm text-[#e94560] text-center">{{ bookError }}</p>
        </div>

        <!-- ─── Step: done ─────────────────────────────────────────────── -->
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
