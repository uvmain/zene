<script setup lang="ts">
defineProps({
  genre: { type: String, required: true },
})
</script>

<template>
  <span class="rounded-full bg-zene-400 px-4 py-1 text-xs text-white font-semibold">
    {{ genre }}
  </span>
</template>
