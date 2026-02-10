<script setup lang="ts">
import type { SearchResult } from '~/types'
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import type { SubsonicSong } from '~/types/subsonicSong'
import { closeSearch, getSearchResults, searchInput } from '~/logic/search'

const searchResults = ref<SearchResult | null>(null)

watch(searchInput, async (newValue) => {
  if (newValue && newValue.length >= 3) {
    searchResults.value = await getSearchResults()
  }
})

function trackToAlbum(track: SubsonicSong): SubsonicAlbum {
  return {
    id: track.albumId,
    coverArt: track.coverArt,
    year: track.year,
    artist: track.artist,
  } as SubsonicAlbum
}
</script>

<template>
  <div v-if="searchInput.length >= 3" class="corner-cut-large mt-2 background-2 bg-gradient-to-br">
    <div class="flex flex-col gap-2 p-4">
      <div class="flex justify-between">
        <h3>
          Search results for "{{ searchInput }}":
        </h3>
        <ZButton size="sm" class="my-auto" @click="closeSearch()">
          Close
        </ZButton>
      </div>
      <h4>
        Tracks: {{ searchResults?.songs?.length ?? 0 }}
      </h4>
      <div class="flex flex-nowrap gap-6 overflow-hidden">
        <div
          v-for="track in searchResults?.songs"
          :key="track.path"
          class="w-auto flex flex-none flex-col gap-y-1 overflow-hidden"
        >
          <div class="overflow-hidden text-ellipsis whitespace-nowrap text-sm text-primary">
            {{ track.title }}
          </div>
          <Album :album="trackToAlbum(track)" />
        </div>
      </div>
      <h4>
        Albums: {{ searchResults?.albums?.length ?? 0 }}
      </h4>
      <div class="flex flex-nowrap gap-6 overflow-hidden">
        <div
          v-for="album in searchResults?.albums"
          :key="album.name"
          class="w-auto flex flex-none flex-col gap-y-1 overflow-hidden"
        >
          <Album :album="album" />
        </div>
      </div>
      <h4>
        Artists: {{ searchResults?.artists?.length ?? 0 }}
      </h4>
      <div class="flex flex-wrap gap-6 overflow-hidden">
        <div
          v-for="artist in searchResults?.artists"
          :key="artist.id"
          class="flex flex-none flex-col gap-y-1 overflow-hidden"
        >
          <ArtistThumb :artist="artist" />
        </div>
      </div>
      <h4>
        Genres: {{ searchResults?.genres?.length ?? 0 }}
      </h4>
      <div class="flex flex-wrap gap-2 overflow-hidden">
        <GenreBottle v-for="genre in searchResults?.genres.filter(g => g.value !== '')" :key="genre.value" :genre="genre.value" />
      </div>
    </div>
  </div>
</template>
