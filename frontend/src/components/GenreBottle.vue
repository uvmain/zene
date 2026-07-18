<script setup lang="ts">
const props = defineProps({
  genre: { type: String, required: true },
})

const router = useRouter()

const buttonTitle = computed(() => {
  return `Navigate to ${props.genre} genre`
})

function navigateToGenre(genre: string) {
  router.push(`/genres/${genre}`)
}
</script>

<template>
  <ZButton :title="buttonTitle" @click="navigateToGenre(genre)">
    <span class="text-nowrap">{{ genre }}</span>
  </ZButton>
</template>
