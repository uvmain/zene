<script setup lang="ts">
import type { SubsonicGenre } from '~/types/subsonicGenres'
import { fetchGenres } from '~/logic/backendFetch'

const genres = ref<SubsonicGenre[]>([])

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
      <hr class="mx-2 border-t border-main-400/20 flex-1 lg:mx-4" />
      <ZButton>
        Top Genres
      </ZButton>
    </div>
    <Genres v-if="genres.length > 0" :genre-strings="genres.map(g => g.value)" :row-limit="2" class="mb-4" />
  </div>
</template>
