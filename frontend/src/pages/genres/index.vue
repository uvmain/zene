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
  genres.value = await fetchGenres()
  genresStore.value = genres.value
}

onBeforeMount(async () => {
  await getGenres()
})
</script>

<template>
  <div>
    <h2 class="text-lg font-semibold py-2">
      Genres
    </h2>
    <div class="flex flex-wrap gap-6">
      <div v-if="genres" class="flex flex-wrap gap-2 cursor-pointer">
        <div v-for="genre in genres.filter(g => g.value !== '')" :key="genre.value">
          <GenreBottle :genre="genre.value" class="cursor-pointer" @click="() => router.push(`/genres/${genre.value}`)" />
        </div>
      </div>
    </div>
  </div>
</template>
