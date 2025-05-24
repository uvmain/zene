<script setup lang="ts">
import type { Component } from 'vue'

import Home from './components/routes/Home.vue'

const routes: Record<string, Component> = {
  '/': Home,
}

const currentPath = ref('#/')

if (typeof window !== 'undefined') {
  currentPath.value = window.location.hash

  window.addEventListener('hashchange', () => {
    currentPath.value = window.location.hash
  })
}

const currentView = computed(() => {
  return routes[currentPath.value.slice(1) || '/'] as Component
})
</script>

<template>
  <div class="grid grid-cols-[250px_1fr] h-screen from-zene-800 to-zene-700 bg-gradient-to-b text-white">
    <Navbar />
    <main class="overflow-y-auto p-6 space-y-6">
      <HeaderAndSearch />
      <component :is="currentView" />
    </main>
  </div>
</template>

<style>
html, body, #app {
  margin: 0;
  padding: 0;
  border: 0;
  font-family: 'Montserrat', sans-serif;
  min-height: 100%;
  @apply standard;
}
</style>
