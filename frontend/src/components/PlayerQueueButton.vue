<script setup lang="ts">
const router = useRouter()

const canNavigateToQueue = computed(() => {
  return router.currentRoute.value.path !== '/queue'
})

function navigateToQueue() {
  router.push('/queue')
}
</script>

<template>
  <abbr title="Queue">
    <button
      id="queue"
      class="h-10 w-10 flex cursor-pointer items-center justify-center border-none bg-white/0 text-muted font-semibold outline-none lg:h-12 lg:w-12 sm:h-10 sm:w-10"
      @click="navigateToQueue()"
    >
      <icon-nrk-media-playlist
        :class="{
          'footer-icon-disabled': !canNavigateToQueue,
          'footer-icon': canNavigateToQueue,
        }"
      />
    </button>
  </abbr>
</template>
