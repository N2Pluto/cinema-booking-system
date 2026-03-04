import { ref, reactive } from 'vue'
import { useAuth } from './useAuth'

export interface Seat {
  seat_no: string
  status: 'AVAILABLE' | 'LOCKED' | 'BOOKED'
  price: number
  user_id?: string
}

export interface Cinema {
  id: string
  movie_name: string
  theater_no: number
  start_time: string
  end_time: string
  price: number
  seats: Seat[]
  created_at: string
}

export interface CinemasResult {
  data: Cinema[]
  total: number
  page: number
  limit: number
  total_pages: number
}

const API_URL = () => import.meta.env.VITE_API_URL || 'http://localhost:8080'

export function useCinemas() {
  const { token } = useAuth()
  const loading = ref(false)
  const error = ref<string | null>(null)
  const result = reactive<CinemasResult>({
    data: [], total: 0, page: 1, limit: 10, total_pages: 1,
  })

  async function fetchCinemas(params: {
    page?: number
    limit?: number
    orderBy?: string
    startDate?: string
    endDate?: string
    admin?: boolean
  } = {}) {
    loading.value = true
    error.value = null
    try {
      const q = new URLSearchParams()
      if (params.page) q.set('page', String(params.page))
      if (params.limit) q.set('limit', String(params.limit))
      if (params.orderBy) q.set('orderBy', params.orderBy)
      if (params.startDate) q.set('startDate', params.startDate)
      if (params.endDate) q.set('endDate', params.endDate)

      const endpoint = params.admin ? '/api/admin/cinema' : '/api/cinema'
      const res = await fetch(`${API_URL()}${endpoint}?${q}`, {
        headers: { Authorization: `Bearer ${token.value}` },
      })
      if (!res.ok) throw new Error(await res.text())
      const data: CinemasResult = await res.json()
      Object.assign(result, data)
    } catch (e: any) {
      error.value = e.message
    } finally {
      loading.value = false
    }
  }

  async function fetchCinema(id: string): Promise<Cinema | null> {
    try {
      const res = await fetch(`${API_URL()}/api/seats/${id}`, {
        headers: { Authorization: `Bearer ${token.value}` },
      })
      if (!res.ok) return null
      return await res.json()
    } catch {
      return null
    }
  }

  return { loading, error, result, fetchCinemas, fetchCinema }
}
