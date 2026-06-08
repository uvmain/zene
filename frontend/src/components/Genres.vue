<script setup lang="ts">
const props = defineProps({
  genreStrings: { type: Array as PropType<string[]>, required: true },
  rowLimit: { type: Number, required: false, default: 0 },
})

const router = useRouter()

const genreBottles = useTemplateRef('genreBottles')

const firstBottleHeight = computed(() => {
  if (genreBottles.value && genreBottles.value.length > 0) {
    const firstBottle = genreBottles.value[0] as HTMLElement
    return firstBottle.offsetHeight
  }
  return 0
})

const containerStyle = computed(() => {
  if (props.rowLimit <= 0 || firstBottleHeight.value <= 0) {
    return undefined
  }

  const gap = 8
  const totalHeight = (props.rowLimit * firstBottleHeight.value) + ((props.rowLimit - 1) * gap)

  return {
    height: `${totalHeight}px`,
  }
})
</script>

<template>
  <div
    v-if="genreStrings.length > 0"
    class="flex flex-wrap gap-8px overflow-y-hidden"
    :style="containerStyle"
  >
    <div v-for="genre in genreStrings" ref="genreBottles" :key="genre">
      <GenreBottle :genre="genre" class="cursor-pointer" @click="() => router.push(`/genres/${genre}`)" />
    </div>
  </div>
</template>
