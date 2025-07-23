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
  <div class="flex flex-row items-center justify-center gap-x-2 py-2 md:justify-start">
    <h2 class="cursor-pointer text-lg font-semibold" @click="emits('titleClick')">
      {{ title }}
    </h2>
    <icon-tabler-refresh class="cursor-pointer text-sm" :class="{ spin: refreshed }" @click="refresh()" />
  </div>
</template>

<style scoped>
@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(-360deg);
  }
}

.spin {
  animation: spin 0.5s linear;
}
</style>
