<script setup lang="ts">
defineProps({
  genre: { type: String, required: true },
})
</script>

<template>
  <span class="rounded-full bg-zenegray-200 px-4 py-1 text-xs text-zenegray-50 font-semibold">
    {{ genre }}
  </span>
</template>
