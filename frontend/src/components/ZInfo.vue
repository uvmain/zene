<script setup lang="ts">
defineProps({
  primary: { type: Boolean, default: false },
  text: { type: String, required: true },
})
</script>

<template>
  <button
    class="text-sm text-muted font-semibold px-3 py-1 align-middle border-1 corner-cut border-solid background-1 flex inline-flex cursor-default items-center justify-center"
    :class="{
      'border-secondary-500': primary,
      'border-background-600 dark:border-background-400': !primary,
    }"
  >
    <span class="text-balance">
      {{ text }}
    </span>
  </button>
</template>
