<script setup lang="ts">
defineProps({
  genre: { type: String, required: true },
})

const router = useRouter()

function navigateToGenre(genre: string) {
  router.push(`/genres/${genre}`)
}
</script>

<template>
  <ZButton
    @click="navigateToGenre(genre)"
  >
    <span class="text-nowrap">{{ genre }}</span>
  </ZButton>
</template>
