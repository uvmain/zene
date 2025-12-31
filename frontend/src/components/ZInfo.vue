<script setup lang="ts">
defineProps({
  primary: { type: Boolean, default: false },
  text: { type: String, required: true },
})
</script>

<template>
  <button
    class="corner-cut flex inline-flex cursor-default items-center justify-center border-1 border-solid background-2 px-3 py-1 align-middle text-sm text-muted font-semibold opacity-70"
    :class="{
      'border-primary2': primary,
      'border-zshade-600 dark:border-zshade-400': !primary,
    }"
  >
    <span class="text-balance">
      {{ text }}
    </span>
  </button>
</template>
