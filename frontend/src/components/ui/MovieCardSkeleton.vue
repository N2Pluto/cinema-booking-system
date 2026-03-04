<template>
  <div class="bg-white/[0.03] border border-white/[0.07] rounded-2xl overflow-hidden">
    <div
      class="w-full aspect-[2/3] animate-shimmer"
      style="background: linear-gradient(90deg, rgba(255,255,255,0.03) 0%, rgba(255,255,255,0.07) 50%, rgba(255,255,255,0.03) 100%); background-size: 200% 100%;"
    ></div>
    <div class="p-3 flex flex-col gap-2">
      <div
        class="h-2.5 rounded w-3/5 animate-shimmer"
        style="background: rgba(255,255,255,0.06);"
      ></div>
      <div
        class="h-2.5 rounded w-[90%] animate-shimmer"
        style="background: rgba(255,255,255,0.06);"
      ></div>
    </div>
  </div>
</template>
