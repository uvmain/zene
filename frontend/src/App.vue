<template>
  <main class="text-muted background-1 flex h-screen lg:grid lg:grid-cols-[200px_1fr]">
    <Navbar />
    <div class="flex flex-1 flex-col h-screen min-h-0">
      <div class="p-2 flex flex-1 flex-col min-h-0 overflow-y-auto space-y-4 lg:p-6 lg:space-y-6">
        <Header class="flex-shrink-0" />
        <RouterView class="flex-1 min-h-0" />
      </div>
      <FooterPlayer class="flex-shrink-0" />
    </div>
  </main>
</template>

<style>
html, body, #app {
  margin: 0;
  padding: 0;
  border: 0;
  font-family: 'Jura', sans-serif;
  min-height: 100dvh;
  font-weight: 600;
  scrollbar-color: var(--colors-zshade-600) var(--colors-zshade-200);
  .dark & {
    scrollbar-color: var(--colors-zshade-300) var(--colors-zshade-800);
  }
}

button, input, select, textarea {
  font-family: inherit;
}

input, textarea {
  @apply text-muted;
  border: 1px solid var(--colors-zshade-300);
  background: var(--colors-zshade-100);
  .dark & {
    border: 1px solid var(--colors-zshade-700);
    background: var(--colors-zshade-800);
  }
}
</style>
