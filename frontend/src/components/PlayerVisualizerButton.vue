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
      class="text-muted font-semibold outline-none border-none bg-white/0 flex cursor-pointer items-center justify-center"
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
