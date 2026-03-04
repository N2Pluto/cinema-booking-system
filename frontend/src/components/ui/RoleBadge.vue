<script setup lang="ts">
defineProps<{ role: string }>()
</script>

<template>
  <span
    class="px-3 py-0.5 rounded-full text-xs font-semibold uppercase tracking-wide border"
    :class="
      role?.toLowerCase() === 'admin'
        ? 'bg-[#e94560]/15 text-[#e94560] border-[#e94560]/30'
        : 'bg-[#64c896]/15 text-[#64c896] border-[#64c896]/30'
    "
  >
    {{ role }}
  </span>
</template>
