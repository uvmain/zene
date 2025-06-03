<script setup lang="ts">
defineProps({
  genre: { type: String, required: true },
})
</script>

<template>
  <span class="rounded-full bg-zene-400 px-3 py-1 text-sm text-white">
    {{ genre }}
  </span>
</template>
