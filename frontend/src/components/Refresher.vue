<script setup lang="ts">
const emits = defineEmits(['refreshed'])

const refreshed = ref(false)

async function refresh() {
  refreshed.value = true
  emits('refreshed')

  setTimeout(() => {
    refreshed.value = false
  }, 1000)
}
</script>

<template>
  <abbr title="Refresh">
    <button class="flex rotate-30 items-center justify-center" @click="refresh()">
      <icon-nrk-refresh class="text-sm cursor-pointer hover:text-secondary-500" :class="{ spin: refreshed }" />
    </button>
  </abbr>
</template>

<style scoped lang="css">
@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.spin {
  animation: spin 0.3s linear;
}
</style>
