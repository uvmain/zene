<script setup lang="ts">
import type { GenreMetadata } from '~/types'
import { useBackendFetch } from '~/composables/useBackendFetch'

const router = useRouter()
const { backendFetchRequest } = useBackendFetch()

const genres = ref<GenreMetadata[]>()

async function getArtists() {
  const response = await backendFetchRequest('genres')
  const json = await response.json() as GenreMetadata[]
  genres.value = json
}

onBeforeMount(async () => {
  await getArtists()
})
</script>

<template>
  <div>
    <h2 class="py-2 text-lg font-semibold">
      Genres
    </h2>
    <div class="flex flex-wrap gap-6">
      <div v-if="genres" class="flex flex-wrap cursor-pointer gap-2">
        <div v-for="genre in genres" :key="genre.genre">
          <GenreBottle :genre="genre.genre" class="cursor-pointer" @click="() => router.push(`/genres/${genre.genre}`)" />
        </div>
      </div>
    </div>
  </div>
</template>
