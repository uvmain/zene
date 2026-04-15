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
  <div class="flex flex-col gap-4">
    <div class="flex flex-row gap-x-4 items-center justify-between">
      <div class="flex flex-row gap-x-2 items-center">
        <h2 class="text-lg font-semibold uppercase lg:text-xl">
          Genres
        </h2>
        <Refresher @refreshed="getGenres()" />
      </div>
      <hr class="mx-2 border-t border-primary-400/20 flex-1 lg:mx-4" />
      <ZButton>
        Top Genres
      </ZButton>
    </div>
    <div class="mb-2 flex flex-wrap gap-2 justify-center overflow-hidden lg:justify-start" :style="`max-height: calc(${(28 * 2) + 12}px);`">
      <GenreBottle v-for="genre in genres?.filter(g => g.value !== '')" :key="genre.value" :genre="genre.value" />
    </div>
  </div>
</template>
