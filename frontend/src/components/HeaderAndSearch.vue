<script setup lang="ts">
import { useDark, useToggle } from '@vueuse/core'
import { useSearch } from '../composables/useSearch'
import { useNavbar } from '../composables/useNavbar'

const { search, searchInput } = useSearch()
const { toggleMobileNav } = useNavbar()

const isDark = useDark()
const toggleDark = useToggle(isDark)
</script>

<template>
  <header>
    <div class="flex p-2">
      <!-- Mobile hamburger menu -->
      <div class="flex items-center md:hidden">
        <button @click="toggleMobileNav()" class="p-2 text-white hover:text-zene-200 transition-colors">
          <icon-tabler-menu-2 class="text-2xl" />
        </button>
      </div>
      
      <div class="flex flex-grow justify-center">
        <div class="relative w-1/2">
          <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400">
            <icon-tabler-search class="text-xl" />
          </span>
          <input
            id="search-input"
            v-model="searchInput"
            placeholder="Type here to search"
            type="text"
            class="block w-full border border-zene-400 rounded-lg bg-gray-800 px-10 py-2 text-white focus:border-zene-200 focus:border-solid focus:shadow-zene-400 hover:shadow-lg focus:outline-none"
            @change="search()"
            @input="search()"
            @keydown.escape="searchInput = ''"
          >
        </div>
      </div>
      <div id="user-and-settings" class="flex gap-4">
        <div class="hover:cursor-pointer" @click="toggleDark()">
          <icon-tabler-sun v-if="isDark" class="text-2xl" />
          <icon-tabler-moon-stars v-else class="text-2xl" />
        </div>
        <SettingsDropDown />
        <div class="hover:cursor-pointer">
          <icon-tabler-user class="text-2xl" />
        </div>
      </div>
    </div>
    <SearchResults />
  </header>
</template>
