<script setup lang="ts">
import type { SubsonicPlaylistsResponse, SubsonicResponse } from '~/types/subsonic'
import type { SubsonicPlaylist } from '~/types/subsonicPlaylists'
import { openSubsonicFetchRequest } from '~/composables/backendFetch'
import { getCoverArtUrl, onImageError } from '~/composables/logic'

const router = useRouter()

const showModal = ref(false)
const isSubmitting = ref(false)
const submitError = ref('')
const showSuccess = ref(false)

const newPlaylistName = ref('')

const playlists = ref<SubsonicPlaylist[]>([])

async function createNewPlaylist() {
  if (isSubmitting.value)
    return
  isSubmitting.value = true
  submitError.value = ''
  showSuccess.value = false
  const formData = new FormData()
  formData.append('name', newPlaylistName.value)
  const response = await openSubsonicFetchRequest<SubsonicResponse>('createPlaylist', {
    body: formData,
  })
  if (response?.error?.message) {
    submitError.value = response.error.message
    isSubmitting.value = false
    return
  }
  showSuccess.value = true
  getPlaylists()
  setTimeout(() => {
    showModal.value = false
    showSuccess.value = false
    isSubmitting.value = false
  }, 1000)
}

async function getPlaylists() {
  const formData = new FormData()
  formData.append('includeEpisodes', true.toString())
  const response = await openSubsonicFetchRequest<SubsonicPlaylistsResponse>('getPlaylists', {
    body: formData,
  })
  playlists.value = response?.playlists?.playlist
  playlists.value.forEach((playlist) => {
    playlist.coverArt = `/share/img/${playlist.coverArt}?size=200`
  })
}

function navigateToPlaylist(playlistId: string) {
  router.push(`/playlists/${playlistId}`)
}

onBeforeMount(getPlaylists)
</script>

<template>
  <div class="p-4 space-y-4">
    <button class="z-button" @click="showModal = true">
      Create New Playlist
    </button>

    <div class="mt-8">
      <div v-if="playlists.length === 0" class="text-zgray-200">
        No playlists found.
      </div>
      <div class="flex flex-wrap justify-center gap-6 md:justify-start">
        <div
          v-for="playlist in playlists"
          :key="playlist.id" class="mx-auto max-w-60dvw flex flex-col items-center justify-center gap-4 transition duration-150 hover:scale-101"
          @click="navigateToPlaylist(`${playlist.id}`)"
        >
          <img
            :src="playlist.coverArt"
            alt="Playlist Cover"
            class="size-40 object-cover"
            loading="lazy"
            width="200"
            height="200"
            @error="onImageError"
          />
          <div class="my-auto flex flex-col gap-4">
            <div class="text-2xl font-bold">
              {{ playlist.name }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <teleport to="body">
      <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-60">
        <div class="relative max-w-md w-full bg-zgray-600 p-6 shadow-lg">
          <button class="z-button absolute right-2 top-2" aria-label="Close" @click="showModal = false">
            X
          </button>
          <h2 class="mb-4 text-xl text-zgray-200 font-bold">
            Add New Playlist
          </h2>
          <form class="space-y-4" @submit.prevent="createNewPlaylist">
            <div>
              <label for="playlist-name" class="mb-1 block text-zgray-200 font-medium">Playlist Name</label>
              <input
                id="playlist-name"
                v-model="newPlaylistName"
                type="text"
                class="w-auto border px-3 py-2"
                placeholder="Enter playlist name"
                required
              />
            </div>
            <button
              type="submit"
              class="z-button"
              :disabled="isSubmitting || showSuccess"
            >
              <span v-if="isSubmitting && !showSuccess">Adding...</span>
              <span v-else-if="showSuccess">
                <svg xmlns="http://www.w3.org/2000/svg" class="inline h-5 w-5 text-green-600" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-7.5 7.5a1 1 0 01-1.414 0l-3.5-3.5a1 1 0 111.414-1.414L8 11.086l6.793-6.793a1 1 0 011.414 0z" clip-rule="evenodd" /></svg>
              </span>
              <span v-else>Add Radio Station</span>
            </button>
            <div v-if="submitError" class="mt-2 text-sm text-red-600">
              {{ submitError }}
            </div>
          </form>
        </div>
      </div>
    </teleport>
  </div>
</template>
