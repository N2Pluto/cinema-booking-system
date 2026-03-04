<script setup lang="ts">
defineProps<{
  page: number
  totalPages: number
  total: number
  limit: number
}>()

const emit = defineEmits<{ (e: 'change', page: number): void }>()
</script>

<template>
  <div class="flex items-center justify-between text-sm text-[#606882]">
    <span>ทั้งหมด {{ total }} รายการ</span>
    <div class="flex items-center gap-1">
      <button
        :disabled="page <= 1"
        @click="emit('change', page - 1)"
        class="px-3 py-1.5 rounded-lg border border-white/[0.1] disabled:opacity-30 disabled:cursor-not-allowed hover:bg-white/[0.06] transition-colors"
      >
        ก่อนหน้า
      </button>
      <span class="px-3 py-1.5 text-white font-medium">{{ page }} / {{ totalPages }}</span>
      <button
        :disabled="page >= totalPages"
        @click="emit('change', page + 1)"
        class="px-3 py-1.5 rounded-lg border border-white/[0.1] disabled:opacity-30 disabled:cursor-not-allowed hover:bg-white/[0.06] transition-colors"
      >
        ถัดไป
      </button>
    </div>
  </div>
</template>
