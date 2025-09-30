<script setup lang="ts">
defineProps({
  title: { type: String, required: true },
})

const emits = defineEmits(['refreshed', 'titleClick'])

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
  <div class="flex flex-row items-center justify-center gap-x-2 py-4 md:justify-start">
    <ZButton @click="emits('titleClick')">
      {{ title }}
    </ZButton>
    <icon-nrk-refresh class="cursor-pointer text-sm hover:text-primary2" :class="{ spin: refreshed }" @click="refresh()" />
  </div>
</template>

<style scoped>
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
