<script setup lang="ts">
import { backendFetchRequest } from '../composables/fetchFromBackend'

const topGenres = ref<any[]>([])

async function getGenres() {
  const response = await backendFetchRequest('genres')
  const json = await response.json()
  topGenres.value = json.slice(0, 12)
}

onBeforeMount(async () => {
  await getGenres()
})
</script>

<template>
  <div>
    <h2 class="py-2 text-lg font-semibold">
      Top Genres
    </h2>
    <div class="flex flex-wrap gap-2">
      <GenreBottle v-for="genre in topGenres" :key="genre" :genre="genre.genre" />
    </div>
  </div>
</template>
