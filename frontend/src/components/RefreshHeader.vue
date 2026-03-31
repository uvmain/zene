<script setup lang="ts">
defineProps({
  title: { type: String, required: true },
})

const emits = defineEmits(['refreshed', 'titleClick', 'close'])

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
  <div class="py-6 flex flex-row gap-x-2 items-center justify-center lg:justify-start">
    <ZButton @click="emits('titleClick')">
      <span class="text-sm lg:text-base">{{ title }}</span>
    </ZButton>
    <abbr title="Refresh">
      <button class="flex items-center justify-center" @click="refresh()">
        <icon-nrk-refresh class="text-sm cursor-pointer hover:text-primary2" :class="{ spin: refreshed }" />
      </button>
    </abbr>
  </div>
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
