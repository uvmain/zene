<script setup lang="ts">
import { useBackendFetch } from '../composables/useBackendFetch'

const { backendFetchRequest } = useBackendFetch()

const topGenres = ref<any[]>([])

async function getGenres() {
  const response = await backendFetchRequest('genres')
  const json = await response.json()
  topGenres.value = json.slice(0, 30)
}

onBeforeMount(async () => {
  await getGenres()
})
</script>

<template>
  <div>
    <RefreshHeader title="Top Genres" @refreshed="getGenres()" />
    <div class="flex flex-wrap justify-center gap-2 md:justify-start">
      <GenreBottle v-for="genre in topGenres" :key="genre" :genre="genre.genre" />
    </div>
  </div>
</template>
