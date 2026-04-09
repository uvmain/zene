<template>
  <main class="text-muted background-1 flex flex-row h-screen">
    <Navbar />
    <div class="flex flex-1 flex-col h-screen min-h-0">
      <div class="p-2 flex flex-1 min-h-0 overflow-y-auto lg:p-6">
        <RouterView class="flex-1 min-h-0" />
      </div>
      <FooterPlayer class="flex-shrink-0" />
    </div>
  </main>
</template>
