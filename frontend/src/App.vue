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
