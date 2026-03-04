import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    // ── Public ──────────────────────────────────────────────────────────────
    { path: '/login', name: 'login', component: () => import('@/views/LoginView.vue') },
    { path: '/auth/callback', name: 'callback', component: () => import('@/views/CallbackView.vue') },

    // ── User ────────────────────────────────────────────────────────────────
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/cinema/:id',
      name: 'cinema-detail',
      component: () => import('@/views/CinemaDetailView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/bookings',
      name: 'bookings',
      component: () => import('@/views/BookingsView.vue'),
      meta: { requiresAuth: true },
    },

    // ── Admin ────────────────────────────────────────────────────────────────
    {
      path: '/admin/cinemas',
      name: 'admin-cinemas',
      component: () => import('@/views/admin/AdminCinemasView.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
    },
    {
      path: '/admin/cinemas/create',
      name: 'admin-cinemas-create',
      component: () => import('@/views/admin/AdminCreateShowtimeView.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
    },
    {
      path: '/admin/audit-logs',
      name: 'admin-audit-logs',
      component: () => import('@/views/admin/AdminAuditLogsView.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
    },
  ],
})

router.beforeEach(async (to) => {
  const { isLoggedIn, user, fetchMe } = useAuth()

  if (to.meta.requiresAuth && !isLoggedIn.value) {
    return { name: 'login' }
  }

  if (to.name === 'login' && isLoggedIn.value) {
    return { name: 'home' }
  }

  if (to.meta.requiresAdmin) {
    // Ensure user is loaded before checking admin role
    if (!user.value) {
      await fetchMe()
    }
    if (!user.value?.is_admin) {
      return { name: 'home' }
    }
  }
})

export default router
