<script setup lang="ts">
import { useDark, useToggle } from '@vueuse/core'
import { useNavbar } from '../composables/useNavbar'
import { useSearch } from '../composables/useSearch'

const { search, searchInput } = useSearch()
const { toggleMobileNav } = useNavbar()

const isDark = useDark()
const toggleDark = useToggle(isDark)
</script>

<template>
  <header>
    <div class="flex gap-2 p-2 md:gap-4">
      <!-- Mobile hamburger menu -->
      <div class="flex items-center justify-center md:hidden">
        <icon-tabler-menu-2 class="text-2xl" @click="toggleMobileNav()" />
      </div>

      <div class="flex flex-grow items-center justify-center">
        <div class="relative max-w-xs w-full md:max-w-md md:w-1/2">
          <span class="absolute inset-y-0 left-0 h-full flex items-center justify-center pl-3 text-gray-400">
            <icon-tabler-search class="text-xl" />
          </span>
          <input
            id="search-input"
            v-model="searchInput"
            placeholder="Type here to search"
            type="text"
            class="border border-zene-400 rounded-lg bg-gray-800 py-2 pl-10 text-white focus:border-zene-200 focus:border-solid focus:shadow-zene-400 hover:shadow-lg focus:outline-none"
            @change="search()"
            @input="search()"
            @keydown.escape="searchInput = ''"
          >
        </div>
      </div>
      <div id="user-and-settings" class="flex gap-4">
        <div class="items-center justify-center hover:cursor-pointer" @click="toggleDark()">
          <icon-tabler-sun v-if="isDark" class="text-2xl" />
          <icon-tabler-moon-stars v-else class="text-2xl" />
        </div>
        <SettingsDropDown class="items-center justify-center" />
        <div class="items-center justify-center hover:cursor-pointer">
          <icon-tabler-user class="text-2xl" />
        </div>
      </div>
    </div>
    <SearchResults />
  </header>
</template>
