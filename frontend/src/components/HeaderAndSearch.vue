<script setup lang="ts">
import { useDark, useToggle } from '@vueuse/core'
import { useNavbar } from '~/composables/useNavbar'
import { useSearch } from '~/composables/useSearch'

const { getSearchResults, searchInput } = useSearch()
const { toggleMobileNav } = useNavbar()

const isDark = useDark()
const toggleDark = useToggle(isDark)
</script>

<template>
  <header>
    <div class="flex gap-2 p-2 md:gap-4">
      <!-- Mobile hamburger menu -->
      <div class="flex items-center justify-center md:hidden">
        <icon-nrk-list class="text-2xl" @click="toggleMobileNav()" />
      </div>

      <div class="flex flex-grow items-center justify-center">
        <div class="relative max-w-xs w-full md:max-w-md md:w-1/2">
          <span class="absolute inset-y-0 left-0 h-full flex items-center justify-center pl-3 text-muted">
            <icon-nrk-search class="text-xl" />
          </span>
          <input
            id="search-input"
            v-model="searchInput"
            placeholder="Type here to search"
            type="text"
            class="border-1 border-primary2 background-2 py-2 pl-10 focus:border-primary2 dark:border-opacity-60 focus:border-solid md:pr-full focus:shadow-primary2 hover:shadow-lg focus:outline-none"
            @change="getSearchResults()"
            @input="getSearchResults()"
            @keydown.escape="searchInput = ''"
          >
        </div>
      </div>
      <div id="user-and-settings" class="flex gap-4 text-muted">
        <div class="items-center justify-center hover:cursor-pointer hover:text-primary" @click="toggleDark()">
          <icon-fluent-dark-theme-24-regular class="text-2xl" />
        </div>
        <SettingsDropDown class="items-center justify-center hover:cursor-pointer hover:text-primary" />
        <div class="items-center justify-center hover:cursor-pointer hover:text-primary">
          <icon-nrk-user-loggedin class="text-2xl" />
        </div>
      </div>
    </div>
    <SearchResults />
  </header>
</template>
