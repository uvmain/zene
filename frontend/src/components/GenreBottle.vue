<script setup lang="ts">
defineProps({
  genre: { type: String, required: true },
})
</script>

<template>
  <span class="bg-zene-400 rounded-full px-4 py-1 text-xs text-white font-semibold">
    {{ genre }}
  </span>
</template>
