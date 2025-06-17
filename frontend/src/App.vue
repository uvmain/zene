<script setup lang="ts">
import { useAuth } from './composables/useAuth'

const { userLoginState } = useAuth()
</script>

<template>
  <div v-if="userLoginState" class="grid grid-cols-[250px_1fr] h-screen from-zene-800 to-zene-700 bg-gradient-to-b text-white">
    <Navbar />
    <main class="flex flex-col overflow-y-auto">
      <div class="flex flex-col overflow-y-auto p-6 space-y-6">
        <HeaderAndSearch />
        <RouterView />
      </div>
      <FooterPlayer />
    </main>
  </div>
  <div v-else>
    <Login />
  </div>
</template>

<style>
html, body, #app {
  margin: 0;
  padding: 0;
  border: 0;
  font-family: 'Montserrat', sans-serif;
  min-height: 100vh;
  @apply standard;
}
</style>
