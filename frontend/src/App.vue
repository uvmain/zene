<template>
  <div class="text-muted background-1 flex h-screen lg:grid lg:grid-cols-[200px_1fr]">
    <Navbar />
    <main class="flex flex-1 flex-col overflow-y-auto">
      <div class="p-3 flex flex-col overflow-y-auto space-y-4 lg:p-6 lg:space-y-6">
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
  min-height: 100dvh;
  font-weight: 600;
}

button, input, select, textarea {
  font-family: inherit;
}
</style>
