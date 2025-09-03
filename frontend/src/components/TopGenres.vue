<script setup lang="ts">
import type { SubsonicGenre } from '~/types/subsonicGenres'
import { fetchGenres } from '~/composables/backendFetch'

const genres = ref<SubsonicGenre[]>()

async function getGenres() {
  genres.value = await fetchGenres(30)
}

onBeforeMount(async () => {
  await getGenres()
})
</script>

<template>
  <div>
    <RefreshHeader title="Top Genres" @refreshed="getGenres()" />
    <div class="flex flex-wrap justify-center gap-2 md:justify-start">
      <GenreBottle v-for="genre in genres" :key="genre.value" :genre="genre.value" />
    </div>
  </div>
</template>
