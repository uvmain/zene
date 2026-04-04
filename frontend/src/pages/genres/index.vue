<script setup lang="ts">
import type { SubsonicGenre } from '~/types/subsonicGenres'
import { fetchGenres } from '~/logic/backendFetch'
import { genresStore } from '~/logic/store'

const router = useRouter()

const genres = ref<SubsonicGenre[]>()

async function getGenres() {
  if (genresStore.value.length > 0) {
    genres.value = genresStore.value
  }
  const fetchedGenres = await fetchGenres()
  if (fetchedGenres && fetchedGenres.length > 0 && JSON.stringify(fetchedGenres) !== JSON.stringify(genres.value)) {
    genres.value = fetchedGenres
    genresStore.value = genres.value
  }
}

onBeforeMount(async () => {
  await getGenres()
})
</script>

<template>
  <div>
    <div v-if="genres" class="flex flex-wrap gap-2">
      <div v-for="genre in genres.filter(g => g.value !== '')" :key="genre.value">
        <GenreBottle :genre="genre.value" class="cursor-pointer" @click="() => router.push(`/genres/${genre.value}`)" />
      </div>
    </div>
    <div class="h-4" />
  </div>
</template>
