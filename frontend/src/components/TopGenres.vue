<script setup lang="ts">
import { useBackendFetch } from '~/composables/useBackendFetch'

const { backendFetchRequest } = useBackendFetch()

const topGenres = ref<any[]>([])

async function getGenres() {
  const formData = new FormData()
  formData.append('limit', '30')
  const response = await backendFetchRequest('genres', {
    method: 'POST',
    body: formData,
  })
  const json = await response.json()
  topGenres.value = json
}

onBeforeMount(async () => {
  await getGenres()
})
</script>

<template>
  <div>
    <RefreshHeader title="Top Genres" @refreshed="getGenres()" />
    <div class="flex flex-wrap justify-center gap-2 md:justify-start">
      <GenreBottle v-for="genre in topGenres" :key="genre" :genre="genre.genre" />
    </div>
  </div>
</template>
