<script setup lang="ts">
import { nowPlayingRoute } from '~/logic/playbackQueue'

const router = useRouter()

const canNavigateToQueue = computed(() => {
  return router.currentRoute.value.path !== nowPlayingRoute.value
})

function navigateToQueue() {
  router.push(nowPlayingRoute.value)
}
</script>

<template>
  <button
    title="Now Playing"
    class="text-muted font-semibold outline-none border-none bg-white/0 flex cursor-pointer items-center justify-center"
    @click="navigateToQueue()"
  >
    <icon-nrk-bulleted-list
      :class="{
        'footer-icon-disabled': !canNavigateToQueue,
        'footer-icon': canNavigateToQueue,
      }"
    />
  </button>
</template>
