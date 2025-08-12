<script setup lang="ts">
import type { Genre, SubsonicGenresResponse } from '~/types/subsonicGenres'
import { useBackendFetch } from '~/composables/useBackendFetch'

const { openSubsonicFetchRequest } = useBackendFetch()

const genres = ref<Genre[]>()

async function getGenres() {
  const response = await openSubsonicFetchRequest('getGenres.view')
  const json = await response.json() as SubsonicGenresResponse
  const allGenres = json['subsonic-response'].genres.genre
  genres.value = allGenres.slice(0, 30)
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
