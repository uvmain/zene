<script setup lang="ts">
import type { SubsonicGenre } from '~/types/subsonicGenres'
import { fetchGenres } from '~/composables/backendFetch'

const router = useRouter()

const genres = ref<SubsonicGenre[]>()

onBeforeMount(async () => {
  genres.value = await fetchGenres()
})
</script>

<template>
  <div>
    <h2 class="py-2 text-lg font-semibold">
      Genres
    </h2>
    <div class="flex flex-wrap gap-6">
      <div v-if="genres" class="flex flex-wrap cursor-pointer gap-2">
        <div v-for="genre in genres.filter(g => g.value !== '')" :key="genre.value">
          <GenreBottle :genre="genre.value" class="cursor-pointer" @click="() => router.push(`/genres/${genre.value}`)" />
        </div>
      </div>
    </div>
  </div>
</template>
