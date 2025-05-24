<script setup lang="ts">
import type { AlbumMetadata, ArtistMetadata, TrackMetadataWithImageUrl } from '../types'
import { useDark, useSessionStorage, useToggle } from '@vueuse/core'
import dayjs from 'dayjs'
import { backendFetchRequest } from '../composables/fetchFromBackend'

const inputText = useSessionStorage<string>('searchInput', '')
const isDark = useDark()
const toggleDark = useToggle(isDark)
const searchResults = ref<TrackMetadataWithImageUrl[]>([])
const searchResultsGenres = ref<any[]>([])
const searchResultsArtists = ref<ArtistMetadata[]>([])

async function search() {
  if (!inputText.value || inputText.value.length < 3) {
    searchResults.value = []
    return
  }
  const response = await backendFetchRequest(`search?search=${inputText.value}`)
  const json = await response.json()
  const albumMetadata: TrackMetadataWithImageUrl[] = []
  json.forEach((metadata: any) => {
    const metadataInstance = {
      id: metadata.id,
      file_id: metadata.file_id,
      filename: metadata.filename,
      format: metadata.format,
      duration: metadata.duration,
      size: metadata.size,
      bitrate: metadata.bitrate,
      title: metadata.title,
      artist: metadata.artist,
      album: metadata.album,
      album_artist: metadata.album_artist ?? metadata.artist,
      track_number: metadata.track_number,
      total_tracks: metadata.total_tracks,
      disc_number: metadata.disc_number,
      total_discs: metadata.total_discs,
      musicbrainz_artist_id: metadata.musicbrainz_artist_id,
      musicbrainz_album_id: metadata.musicbrainz_album_id,
      musicbrainz_track_id: metadata.musicbrainz_track_id,
      label: metadata.label,
      genres: metadata.genre.split(';').filter((genre: string) => genre !== ''),
      release_date: dayjs(metadata.release_date).format('YYYY'),
      image_url: `/api/albums/${metadata.musicbrainz_album_id}/art?size=lg`,
    }
    albumMetadata.push(metadataInstance)
  })
  searchResults.value = albumMetadata
  getGenres()
  getArtists()
}

const searchResultsAlbums = computed(() => {
  const uniqueAlbums = new Map<string, AlbumMetadata>()
  searchResults.value.forEach((album: TrackMetadataWithImageUrl) => {
    if (!uniqueAlbums.has(album.musicbrainz_album_id)) {
      uniqueAlbums.set(album.musicbrainz_album_id, {
        artist: album.artist,
        album_artist: album.album_artist ?? album.artist,
        album: album.album,
        musicbrainz_track_id: album.musicbrainz_track_id,
        musicbrainz_album_id: album.musicbrainz_album_id,
        musicbrainz_artist_id: album.musicbrainz_artist_id,
        genres: album.genres,
        release_date: album.release_date,
        image_url: album.image_url,
      })
    }
  })
  return Array.from(uniqueAlbums.values())
})

const searchResultsTracks = computed(() => {
  return searchResults.value
})

async function getGenres() {
  const response = await backendFetchRequest(`genres?search=${inputText.value}`) as any
  const json = await response.json()
  if (!json || json.length === 0) {
    searchResultsGenres.value = []
    return
  }
  searchResultsGenres.value = json.slice(0, 12)
}

async function getArtists() {
  const response = await backendFetchRequest(`artists?search=${inputText.value}`)
  const json = await response.json()
  if (!json || json.length === 0) {
    searchResultsArtists.value = []
    return
  }
  searchResultsArtists.value = json
}
</script>

<template>
  <header>
    <div class="flex p-2">
      <div class="flex flex-grow justify-center">
        <div class="relative w-1/2">
          <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400">
            <icon-tabler-search class="text-xl" />
          </span>
          <input
            id="search-input"
            v-model="inputText"
            placeholder="Type here to search"
            type="text"
            class="block w-full border border-zene-400 rounded-lg bg-gray-800 px-10 py-2 text-white focus:border-zene-200 focus:border-solid focus:shadow-zene-400 hover:shadow-lg focus:outline-none"
            @change="search()"
            @input="search()"
            @keydown.escape="inputText = ''"
          >
        </div>
      </div>
      <div id="user-and-settings" class="flex gap-4">
        <div class="hover:cursor-pointer" @click="toggleDark()">
          <icon-tabler-sun v-if="isDark" class="text-2xl" />
          <icon-tabler-moon-stars v-else class="text-2xl" />
        </div>
        <div class="hover:cursor-pointer">
          <icon-tabler-user class="text-2xl" />
        </div>
      </div>
    </div>
    <div v-if="inputText.length >= 3" class="mt-2 rounded-lg from-zene-400 to-zene-700 bg-gradient-to-b">
      <div class="flex flex-col gap-2 p-4">
        <h3>
          Search results for "{{ inputText }}":
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
        <div class="flex flex-nowrap gap-6">
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
        <div class="flex flex-wrap gap-6">
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
        <div v-if="searchResultsGenres.length > 0" class="flex flex-wrap gap-2">
          <GenreBottle v-for="genre in searchResultsGenres" :key="genre" :genre="genre.genre" />
        </div>
      </div>
    </div>
  </header>
</template>
