<template>
  <div class="h-screen flex from-zene-800 to-zene-700 bg-gradient-to-b text-white md:grid md:grid-cols-[200px_1fr]">
    <Navbar />
    <main class="flex flex-1 flex-col overflow-y-auto">
      <div class="flex flex-col overflow-y-auto p-3 space-y-4 md:p-6 md:space-y-6">
        <HeaderAndSearch />
        <RouterView />
      </div>
      <FooterPlayer />
    </main>
  </div>
</template>

<style>
html, body, #app {
  margin: 0;
  padding: 0;
  border: 0;
  font-family: 'Jura', sans-serif;
  min-height: 100vh;
  @apply standard;
}

button, input, select, textarea {
  font-family: inherit;
}
</style>
