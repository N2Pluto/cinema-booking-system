import { ref } from 'vue'
import { useAuth } from './useAuth'

export interface Booking {
  id: string
  user_id: string
  cinema_id: string
  seat_numbers: string[]
  total_amount: number
  status: 'PENDING' | 'SUCCESS' | 'TIMEOUT' | 'CANCELLED'
  created_at: string
}

const API_URL = () => import.meta.env.VITE_API_URL || 'http://localhost:8080'

export function useBooking() {
  const { token } = useAuth()
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function lockSeats(cinemaId: string, seatNumbers: string[]): Promise<Booking | null> {
    loading.value = true
    error.value = null
    try {
      const res = await fetch(`${API_URL()}/api/booking/lock`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token.value}`,
        },
        body: JSON.stringify({ cinema_id: cinemaId, seat_numbers: seatNumbers }),
      })
      if (!res.ok) {
        const body = await res.json()
        throw new Error(body.error || 'Lock failed')
      }
      return await res.json()
    } catch (e: any) {
      error.value = e.message
      return null
    } finally {
      loading.value = false
    }
  }

  async function confirmBooking(bookingId: string): Promise<Booking | null> {
    loading.value = true
    error.value = null
    try {
      const res = await fetch(`${API_URL()}/api/booking/confirm`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token.value}`,
        },
        body: JSON.stringify({ booking_id: bookingId }),
      })
      if (!res.ok) {
        const body = await res.json()
        throw new Error(body.error || 'Confirm failed')
      }
      return await res.json()
    } catch (e: any) {
      error.value = e.message
      return null
    } finally {
      loading.value = false
    }
  }

  async function fetchMyBookings(): Promise<Booking[]> {
    loading.value = true
    error.value = null
    try {
      const res = await fetch(`${API_URL()}/api/booking/me`, {
        headers: { Authorization: `Bearer ${token.value}` },
      })
      if (!res.ok) throw new Error(await res.text())
      const data = await res.json()
      return data.data ?? []
    } catch (e: any) {
      error.value = e.message
      return []
    } finally {
      loading.value = false
    }
  }

  return { loading, error, lockSeats, confirmBooking, fetchMyBookings }
}
