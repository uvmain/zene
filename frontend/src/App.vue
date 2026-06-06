<template>
  <main class="text-muted background-1 flex flex-row h-100svh overflow-x-hidden">
    <Navbar />
    <PlayerAudio />
    <div class="flex flex-1 flex-col h-100svh min-h-0 overflow-x-hidden">
      <div class="p-2 flex flex-1 min-h-0 overflow-x-hidden overflow-y-auto lg:p-6 md:p-4">
        <RouterView class="flex-1 min-h-0 min-w-0" />
      </div>
      <Footer class="flex-shrink-0" />
    </div>
  </main>
</template>
