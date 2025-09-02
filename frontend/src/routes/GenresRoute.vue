<script setup lang="ts">
import type { SubsonicGenre } from '~/types/subsonicGenres'
import { fetchGenres } from '~/composables/backendFetch'

const router = useRouter()

const genres = ref<SubsonicGenre[]>()

async function getGenres() {
  genres.value = await fetchGenres()
}

onBeforeMount(async () => {
  await getGenres()
})
</script>

<template>
  <div>
    <h2 class="py-2 text-lg font-semibold">
      Genres
    </h2>
    <div class="flex flex-wrap gap-6">
      <div v-if="genres" class="flex flex-wrap cursor-pointer gap-2">
        <div v-for="genre in genres" :key="genre.value">
          <GenreBottle :genre="genre.value" class="cursor-pointer" @click="() => router.push(`/genres/${genre.value}`)" />
        </div>
      </div>
    </div>
  </div>
</template>
