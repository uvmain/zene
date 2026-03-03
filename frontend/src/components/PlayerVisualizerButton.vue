<script setup lang="ts">
const router = useRouter()

const canNavigateToVisualizer = computed(() => {
  return router.currentRoute.value.path !== '/visualizer'
})

function navigateToVisualizer() {
  router.push('/visualizer')
}
</script>

<template>
  <abbr title="Visualizer">
    <button
      id="visualizer"
      class="h-10 w-10 flex cursor-pointer items-center justify-center border-none bg-white/0 text-muted font-semibold outline-none lg:h-12 lg:w-12 sm:h-10 sm:w-10"
      @click="navigateToVisualizer()"
    >
      <icon-nrk-media-soundwave
        :class="{
          'footer-icon-disabled': !canNavigateToVisualizer,
          'footer-icon': canNavigateToVisualizer,
        }"
      />
    </button>
  </abbr>
</template>
