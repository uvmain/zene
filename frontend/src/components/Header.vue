<script setup lang="ts">
import { useDark, useSessionStorage, useToggle } from '@vueuse/core'

defineProps({
  showAdd: { type: Boolean, default: false },
  showEdit: { type: Boolean, default: false },
})

const inputText = useSessionStorage<string>('searchInput', '')

const isDark = useDark()
const toggleDark = useToggle(isDark)
</script>

<template>
  <header class="flex p-2">
    <div class="flex flex-grow justify-center">
      <input
        id="search-input"
        v-model="inputText"
        placeholder="Type here to search"
        type="text"
        class="block w-1/2 border border-zene-400 rounded rounded-lg bg-gray-800 px-4 py-2 text-white focus-border-1 focus:border-zene-200 focus:border-solid focus:shadow-zene-400 hover:shadow-lg focus:outline-none"
      >
    </div>
    <div id="user-and-settings" class="flex gap-4">
      <div class="hover:cursor-pointer" @click="toggleDark()">
        <icon-tabler-sun v-if="isDark" class="text-2xl" />
        <icon-tabler-moon-stars v-else class="text-2xl" />
      </div>
      <div class="hover:cursor-pointer">
        <icon-tabler-user class="text-2xl" />
      </div>
    </div>
  </header>
</template>
