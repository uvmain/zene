<script setup lang="ts">
import type { Genre, SubsonicGenresResponse } from '~/types/subsonicGenres'
import { useBackendFetch } from '~/composables/useBackendFetch'

const router = useRouter()
const { openSubsonicFetchRequest } = useBackendFetch()

const genres = ref<Genre[]>()

async function getGenres() {
  const response = await openSubsonicFetchRequest('getGenres.view')
  const json = await response.json() as SubsonicGenresResponse
  genres.value = json['subsonic-response'].genres.genre
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
