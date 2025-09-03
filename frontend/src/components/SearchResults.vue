<script setup lang="ts">
import type { SearchResult } from '~/types'
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import { fetchAlbum } from '~/composables/backendFetch'
import { useSearch } from '../composables/useSearch'

const { searchInput, getSearchResults } = useSearch()

const searchResults = ref<SearchResult | null>(null)
const trackAlbums = ref<SubsonicAlbum[]>([])

watch(searchInput, async (newValue) => {
  if (newValue.length >= 3) {
    searchResults.value = await getSearchResults()
    searchResults.value?.songs.forEach(async (track) => {
      const album = await fetchAlbum(track.albumId)
      if (!trackAlbums.value.some(a => a.id === album.id)) {
        trackAlbums.value.push(album)
      }
    })
  }
})
</script>

<template>
  <div v-if="searchInput.length >= 3" class="mt-2 rounded-lg from-zene-400 to-zene-700 bg-gradient-to-b">
    <div class="flex flex-col gap-2 p-4">
      <h3>
        Search results for "{{ searchInput }}":
      </h3>
      <h4>
        Tracks: {{ searchResults?.songs.length || 0 }}
      </h4>
      <div class="flex flex-nowrap gap-6 overflow-hidden">
        <div
          v-for="track in searchResults?.songs"
          :key="track.title"
          class="w-30 flex flex-none flex-col gap-y-1 overflow-hidden"
        >
          <div class="overflow-hidden text-ellipsis whitespace-nowrap text-sm text-gray-300">
            {{ track.track }}
          </div>
          <div class="overflow-hidden text-ellipsis whitespace-nowrap text-sm text-gray-300">
            {{ track.title }}
          </div>
          <Album
            v-if="trackAlbums.find(a => a.id === track.albumId)"
            :album="trackAlbums.find(a => a.id === track.albumId)!"
            size="lg"
          />
        </div>
      </div>
      <h4>
        Albums: {{ searchResults?.albums.length || 0 }}
      </h4>
      <div class="flex flex-nowrap gap-6 overflow-hidden">
        <div
          v-for="album in searchResults?.albums"
          :key="album.name"
          class="w-30 flex flex-none flex-col gap-y-1 overflow-hidden"
        >
          <Album :album="album" size="lg" />
        </div>
      </div>
      <h4>
        Artists: {{ searchResults?.artists.length || 0 }}
      </h4>
      <div class="flex flex-wrap gap-6 overflow-hidden">
        <div
          v-for="artist in searchResults?.artists"
          :key="artist.name"
          class="w-30 flex flex-none flex-col gap-y-1 overflow-hidden"
        >
          <ArtistThumb :artist="artist" />
        </div>
      </div>
      <h4>
        Genres: {{ searchResults?.genres.length || 0 }}
      </h4>
      <div class="flex flex-wrap gap-2 overflow-hidden">
        <GenreBottle v-for="genre in searchResults?.genres" :key="genre.value" :genre="genre.value" />
      </div>
    </div>
  </div>
</template>
