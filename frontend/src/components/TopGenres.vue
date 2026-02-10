<script setup lang="ts">
import type { SubsonicGenre } from '~/types/subsonicGenres'
import { fetchGenres } from '~/logic/backendFetch'

const genres = ref<SubsonicGenre[]>()

async function getGenres() {
  genres.value = await fetchGenres(50)
}

onBeforeMount(async () => {
  await getGenres()
})
</script>

<template>
  <div>
    <RefreshHeader title="Top Genres" @refreshed="getGenres()" />
    <div class="flex flex-wrap justify-center gap-2 overflow-hidden lg:justify-start" :style="`max-height: calc(${(28 * 2) + 12}px);`">
      <GenreBottle v-for="genre in genres?.filter(g => g.value !== '')" :key="genre.value" :genre="genre.value" />
    </div>
  </div>
</template>
