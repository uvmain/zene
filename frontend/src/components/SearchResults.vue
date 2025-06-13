<script setup lang="ts">
import { useSearch } from '../composables/useSearch'

const { searchInput, searchResultsGenres, searchResultsAlbums, searchResultsArtists, searchResultsTracks } = useSearch()
</script>

<template>
  <div v-if="searchInput.length >= 3" class="mt-2 rounded-lg from-zene-400 to-zene-700 bg-gradient-to-b">
    <div class="flex flex-col gap-2 p-4">
      <h3>
        Search results for "{{ searchInput }}":
      </h3>
      <h4>
        Tracks: {{ searchResultsTracks.length }}
      </h4>
      <div class="flex flex-nowrap gap-6 overflow-hidden">
        <div
          v-for="track in searchResultsTracks"
          :key="track.title"
          class="w-30 flex flex-none flex-col gap-y-1 overflow-hidden"
        >
          <div class="overflow-hidden text-ellipsis whitespace-nowrap text-sm text-gray-300">
            {{ track.track_number }} / {{ track.total_tracks }}
          </div>
          <div class="overflow-hidden text-ellipsis whitespace-nowrap text-sm text-gray-300">
            {{ track.title }}
          </div>
          <Album :album="track" size="lg" />
        </div>
      </div>
      <h4>
        Albums: {{ searchResultsAlbums.length }}
      </h4>
      <div class="flex flex-nowrap gap-6 overflow-hidden">
        <div
          v-for="album in searchResultsAlbums"
          :key="album.album"
          class="w-30 flex flex-none flex-col gap-y-1 overflow-hidden"
        >
          <Album :album="album" size="lg" />
        </div>
      </div>
      <h4>
        Artists: {{ searchResultsArtists.length }}
      </h4>
      <div class="flex flex-wrap gap-6 overflow-hidden">
        <div
          v-for="artist in searchResultsArtists"
          :key="artist.artist"
          class="w-30 flex flex-none flex-col gap-y-1 overflow-hidden"
        >
          <ArtistThumb :artist="artist" />
        </div>
      </div>
      <h4>
        Genres: {{ searchResultsGenres.length }}
      </h4>
      <div v-if="searchResultsGenres.length > 0" class="flex flex-wrap gap-2 overflow-hidden">
        <GenreBottle v-for="genre in searchResultsGenres" :key="genre" :genre="genre.genre" />
      </div>
    </div>
  </div>
</template>
