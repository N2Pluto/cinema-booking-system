<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import AppNavbar from '@/components/layout/AppNavbar.vue'

const router = useRouter()
const { token } = useAuth()
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const loading = ref(false)
const error = ref<string | null>(null)
const success = ref(false)

const form = ref({
  movie_name: '',
  theater_no: 1,
  start_time: '',
  end_time: '',
  price: 200,
  rows: 5,
  seats_per_row: 10,
})

async function handleSubmit() {
  error.value = null
  loading.value = true
  try {
    const body = {
      ...form.value,
      theater_no: Number(form.value.theater_no),
      price: Number(form.value.price),
      rows: Number(form.value.rows),
      seats_per_row: Number(form.value.seats_per_row),
      start_time: new Date(form.value.start_time).toISOString(),
      end_time: new Date(form.value.end_time).toISOString(),
    }
    const res = await fetch(`${API_URL}/api/admin/showtimes`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token.value}`,
      },
      body: JSON.stringify(body),
    })
    if (!res.ok) {
      const data = await res.json()
      throw new Error(data.error || 'เกิดข้อผิดพลาด')
    }
    success.value = true
  } catch (e: any) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-[#0d0d1a] text-white">
    <AppNavbar />

    <main class="max-w-xl mx-auto px-6 py-10">
      <button @click="router.back()" class="text-sm text-[#606882] hover:text-white mb-6 flex items-center gap-1 transition-colors">
        ← กลับ
      </button>

      <h2 class="text-xl font-bold mb-8">เพิ่มรอบหนัง</h2>

      <div v-if="success" class="bg-[#64c896]/10 border border-[#64c896]/30 rounded-2xl p-6 text-center">
        <div class="text-3xl mb-3">✅</div>
        <p class="font-semibold text-[#64c896] mb-4">สร้างรอบหนังสำเร็จ</p>
        <div class="flex gap-3">
          <button
            @click="success = false; Object.assign(form, { movie_name: '', start_time: '', end_time: '' })"
            class="flex-1 py-2.5 border border-white/20 rounded-xl text-[#a8b2d8] hover:bg-white/5 text-sm transition-colors"
          >
            เพิ่มอีก
          </button>
          <router-link
            to="/admin/cinemas"
            class="flex-1 py-2.5 bg-[#e94560] text-white rounded-xl font-semibold text-sm text-center hover:bg-[#d63651] transition-colors"
          >
            ดูรายการ
          </router-link>
        </div>
      </div>

      <form v-else @submit.prevent="handleSubmit" class="flex flex-col gap-5">
        <!-- Movie name -->
        <div>
          <label class="block text-sm text-[#a8b2d8] mb-1.5">ชื่อภาพยนตร์ *</label>
          <input
            v-model="form.movie_name"
            required
            placeholder="เช่น Avengers: Endgame"
            class="w-full bg-white/[0.05] border border-white/[0.1] rounded-xl px-4 py-2.5 text-white placeholder-[#606882] focus:outline-none focus:border-[#e94560]/60 transition-colors"
          />
        </div>

        <!-- Theater & Price row -->
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm text-[#a8b2d8] mb-1.5">โรงหนัง *</label>
            <input
              v-model="form.theater_no"
              type="number" min="1" required
              class="w-full bg-white/[0.05] border border-white/[0.1] rounded-xl px-4 py-2.5 text-white focus:outline-none focus:border-[#e94560]/60 transition-colors"
            />
          </div>
          <div>
            <label class="block text-sm text-[#a8b2d8] mb-1.5">ราคา (บาท) *</label>
            <input
              v-model="form.price"
              type="number" min="1" required
              class="w-full bg-white/[0.05] border border-white/[0.1] rounded-xl px-4 py-2.5 text-white focus:outline-none focus:border-[#e94560]/60 transition-colors"
            />
          </div>
        </div>

        <!-- Start & End time -->
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm text-[#a8b2d8] mb-1.5">เริ่ม *</label>
            <input
              v-model="form.start_time"
              type="datetime-local" required
              class="w-full bg-white/[0.05] border border-white/[0.1] rounded-xl px-4 py-2.5 text-white focus:outline-none focus:border-[#e94560]/60 transition-colors [color-scheme:dark]"
            />
          </div>
          <div>
            <label class="block text-sm text-[#a8b2d8] mb-1.5">สิ้นสุด *</label>
            <input
              v-model="form.end_time"
              type="datetime-local" required
              class="w-full bg-white/[0.05] border border-white/[0.1] rounded-xl px-4 py-2.5 text-white focus:outline-none focus:border-[#e94560]/60 transition-colors [color-scheme:dark]"
            />
          </div>
        </div>

        <!-- Seat layout -->
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm text-[#a8b2d8] mb-1.5">จำนวนแถว (A-Z) *</label>
            <input
              v-model="form.rows"
              type="number" min="1" max="26" required
              class="w-full bg-white/[0.05] border border-white/[0.1] rounded-xl px-4 py-2.5 text-white focus:outline-none focus:border-[#e94560]/60 transition-colors"
            />
          </div>
          <div>
            <label class="block text-sm text-[#a8b2d8] mb-1.5">ที่นั่งต่อแถว *</label>
            <input
              v-model="form.seats_per_row"
              type="number" min="1" max="50" required
              class="w-full bg-white/[0.05] border border-white/[0.1] rounded-xl px-4 py-2.5 text-white focus:outline-none focus:border-[#e94560]/60 transition-colors"
            />
          </div>
        </div>

        <p class="text-xs text-[#606882]">
          จำนวนที่นั่งทั้งหมด: <strong class="text-[#a8b2d8]">{{ form.rows * form.seats_per_row }}</strong> ที่นั่ง ({{ form.rows }} แถว × {{ form.seats_per_row }})
        </p>

        <p v-if="error" class="text-sm text-[#e94560]">{{ error }}</p>

        <button
          type="submit"
          :disabled="loading"
          class="w-full py-3 bg-[#e94560] text-white font-semibold rounded-xl hover:bg-[#d63651] transition-colors disabled:opacity-50 mt-2"
        >
          {{ loading ? 'กำลังสร้าง...' : 'สร้างรอบหนัง' }}
        </button>
      </form>
    </main>
  </div>
</template>
