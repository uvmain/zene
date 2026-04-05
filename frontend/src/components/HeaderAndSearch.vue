<script setup lang="ts">
import { useDark } from '@vueuse/core'
import { openMobileNav } from '~/logic/navbar'
import { getSearchResults, searchInput } from '~/logic/search'

const isDark = useDark()

function toggleDark() {
  if (!document.startViewTransition) {
    // Fallback for browsers that don't support view transitions
    isDark.value = !isDark.value
    return
  }
  document.startViewTransition(() => {
    isDark.value = !isDark.value
  })
}
</script>

<template>
  <header>
    <div class="p-2 flex gap-2 lg:gap-4">
      <div class="flex items-center justify-center lg:hidden">
        <icon-nrk-list class="lg:text-2xl" @click="openMobileNav()" />
      </div>

      <div class="flex items-center justify-center lg:flex-grow">
        <div class="relative lg:(max-w-md w-1/2)">
          <span class="text-muted pl-3 flex h-full items-center inset-y-0 left-0 justify-center absolute">
            <icon-nrk-search class="lg:text-xl" />
          </span>
          <input
            id="search-input"
            v-model="searchInput"
            placeholder="Type here to search"
            type="text"
            class="py-2 pl-10 border-1 border-primary2 background-2 lg:pr-full focus:outline-none focus:border-primary2 dark:border-opacity-60 focus:border-solid focus:shadow-primary2 hover:shadow-lg"
            @change="getSearchResults()"
            @input="getSearchResults()"
            @keydown.escape="searchInput = ''"
          >
        </div>
      </div>
      <div id="user-and-settings" class="text-muted flex gap-2 items-center justify-center lg:gap-4">
        <abbr :title="isDark ? 'Light mode' : 'Dark mode'" class="icon" @click="toggleDark()">
          <icon-fluent-dark-theme-24-regular class="lg:text-2xl" />
        </abbr>
        <SettingsDropDown class="icon" />
        <abbr title="User" class="icon">
          <icon-nrk-user-loggedin class="lg:text-2xl" />
        </abbr>
      </div>
    </div>
    <SearchResults />
  </header>
</template>

<style scoped lang="css">
.icon {
  @apply flex items-center justify-center hover:cursor-pointer hover:text-primary opacity-70 hover:opacity-100;
}
</style>
