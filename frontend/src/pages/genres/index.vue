<script setup lang="ts">
import type { SubsonicGenre } from '~/types/subsonicGenres'
import { fetchGenres } from '~/logic/backendFetch'
import { getStoredKV, setStoredKV } from '~/stores/keyValueIdbStore'

const router = useRouter()

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
  <div>
    <div v-if="genres" class="flex flex-wrap gap-2">
      <div v-for="genre in genres.filter(g => g.value !== '')" :key="genre.value">
        <GenreBottle :genre="genre.value" class="cursor-pointer" @click="() => router.push(`/genres/${genre.value}`)" />
      </div>
    </div>
    <div class="h-4" />
  </div>
</template>
