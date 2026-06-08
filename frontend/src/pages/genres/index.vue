<script setup lang="ts">
import type { SubsonicGenre } from '~/types/subsonicGenres'
import { fetchGenres } from '~/logic/backendFetch'
import { getStoredKV, setStoredKV } from '~/stores/keyValueIdbStore'

const genres = ref<SubsonicGenre[]>()

async function getGenres() {
  const storedGenres = await getStoredKV('genres')
  if (storedGenres) {
    genres.value = JSON.parse(storedGenres) as SubsonicGenre[]
  }
  const fetchedGenres = await fetchGenres()
  if (fetchedGenres && fetchedGenres.length > 0 && JSON.stringify(fetchedGenres) !== JSON.stringify(genres.value)) {
    genres.value = fetchedGenres
    await setStoredKV('genres', JSON.stringify(genres.value))
  }
}

onBeforeMount(async () => {
  await getGenres()
})
</script>

<template>
  <Genres :genre-strings="genres ? genres.map((g) => g.value) : []" />
</template>
