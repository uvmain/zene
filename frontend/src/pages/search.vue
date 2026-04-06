<script setup lang="ts">
import type { SearchResult } from '~/types'
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicSong } from '~/types/subsonicSong'
import { fetchSearchResults } from '~/logic/backendFetch'

const searchInput = ref<string>('')
const searchResults = ref<SearchResult | null>(null)
const searchInputRef = useTemplateRef('searchInputRef')

async function getSearchResults(): Promise<SearchResult> {
  if (!searchInput.value || searchInput.value.length < 3) {
    return Promise.resolve({} as SearchResult)
  }
  return fetchSearchResults(searchInput.value)
}

function trackToAlbum(track: SubsonicSong): SubsonicAlbum {
  return {
    id: track.albumId,
    coverArt: track.coverArt,
    year: track.year,
    artist: track.artist,
    title: track.album,
  } as SubsonicAlbum
}

watch(searchInput, async (newValue) => {
  if (newValue && newValue.length >= 3) {
    searchResults.value = await getSearchResults()
  }
})

onMounted(() => {
  searchInputRef.value?.focus()
})
</script>

<template>
  <div class="flex flex-col gap-4">
    <div aria-label="Search input" class="relative lg:(max-w-md w-1/2)">
      <span class="text-muted pl-3 flex h-full items-center inset-y-0 left-0 justify-center absolute">
        <icon-nrk-search class="lg:text-xl" />
      </span>
      <input
        id="search-input"
        ref="searchInputRef"
        v-model="searchInput"
        placeholder="Type here to search"
        type="text"
        class="py-2 pl-10 border-1 border-primary2 background-2 lg:pr-full focus:outline-none focus:border-primary2 dark:border-opacity-60 focus:border-solid focus:shadow-primary2 hover:shadow-lg"
        @change="getSearchResults()"
        @input="getSearchResults()"
        @keydown.escape="searchInput = ''"
      >
    </div>
    <div v-if="searchInput.length >= 3" class="corner-cut background-2">
      <div class="p-4 flex flex-col gap-4">
        <div aria-label="Search results for tracks" class="flex flex-col gap-2">
          <div class="text-lg font-bold">
            Tracks: {{ searchResults?.songs?.length ?? 0 }}
          </div>
          <div class="flex flex-wrap gap-6">
            <Album v-for="track in searchResults?.songs" :key="track.path" :album="trackToAlbum(track)" :track-title="track.title" :show-artist="true" :show-date="false" />
          </div>
        </div>
        <hr class="border-primary2 opacity-50" />
        <div aria-label="Search results for albums" class="flex flex-col gap-2">
          <div class="text-lg font-bold">
            Albums: {{ searchResults?.albums?.length ?? 0 }}
          </div>
          <div class="flex flex-wrap gap-6 overflow-hidden">
            <Album v-for="album in searchResults?.albums" :key="album.name" :album="album" />
          </div>
        </div>
        <hr class="border-primary2 opacity-50" />
        <div aria-label="Search results for artists" class="flex flex-col gap-2">
          <div class="text-lg font-bold">
            Artists: {{ searchResults?.artists?.length ?? 0 }}
          </div>
          <div class="flex flex-wrap gap-6 overflow-hidden">
            <ArtistThumb v-for="artist in searchResults?.artists" :key="artist.id" :artist="artist" />
          </div>
        </div>
        <hr class="border-primary2 opacity-50" />
        <div aria-label="Search results for genres" class="flex flex-col gap-2">
          <div class="text-lg font-bold">
            Genres: {{ searchResults?.genres?.length ?? 0 }}
          </div>
          <div class="flex flex-wrap gap-2 overflow-hidden">
            <GenreBottle v-for="genre in searchResults?.genres.filter(g => g.value !== '')" :key="genre.value" :genre="genre.value" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
