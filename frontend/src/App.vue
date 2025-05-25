<script setup lang="ts">
</script>

<template>
  <div class="grid grid-cols-[250px_1fr] h-screen from-zene-800 to-zene-700 bg-gradient-to-b text-white">
    <Navbar />
    <main class="overflow-y-auto p-6 space-y-6">
      <HeaderAndSearch />
      <RouterView />
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
