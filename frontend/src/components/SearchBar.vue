<script setup lang="ts">
import { useSessionStorage } from '@vueuse/core'

const inputText = useSessionStorage<string>('searchInput', '')
</script>

<template>
  <div class="flex rounded shadow-md hover:shadow-lg">
    <input
      id="search-input"
      v-model="inputText"
      placeholder="Type here to search"
      type="text"
    >
  </div>
</template>
