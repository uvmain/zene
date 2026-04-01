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
      class="text-muted font-semibold outline-none border-none bg-white/0 flex cursor-pointer items-center justify-center"
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
