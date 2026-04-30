<script setup lang="ts">
import type { SubsonicPlaylistsResponse, SubsonicResponse } from '~/types/subsonic'
import type { SubsonicPlaylist } from '~/types/subsonicPlaylists'
import { getServerUrl, openSubsonicFetchRequest } from '~/logic/backendFetch'
import { onImageError } from '~/logic/common'

const showModal = ref(false)
const isSubmitting = ref(false)
const submitError = ref('')
const showSuccess = ref(false)
const newPlaylistName = ref('')
const playlists = ref<SubsonicPlaylist[]>([])

const router = useRouter()

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
    playlist.coverArt = getServerUrl(`/share/img/${playlist.coverArt}?size=200`)
  })
}

function navigateToPlaylist(playlistId: number) {
  router.push(`/playlists/${playlistId}`)
}

onBeforeMount(getPlaylists)
</script>

<template>
  <div class="p-4 space-y-4">
    <ZButton @click="showModal = true">
      <span class="text-nowrap">Create New Playlist</span>
    </ZButton>

    <div class="mt-8">
      <div v-if="playlists.length === 0" class="text-primary">
        No playlists found.
      </div>
      <div class="flex flex-wrap gap-6 justify-center lg:justify-start">
        <div
          v-for="(playlist, index) in playlists"
          :key="playlist.id" class="mx-auto flex flex-col gap-4 max-w-60dvw transition duration-150 items-center justify-center hover:scale-101"
          @click="navigateToPlaylist(playlist.id)"
        >
          <img
            :src="playlist.coverArt"
            alt="Playlist Cover"
            class="rounded-md size-40 object-cover"
            :loading="index < 20 ? 'eager' : 'lazy'"
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
      <div v-if="showModal" class="flex items-center inset-0 justify-center fixed z-50 backdrop-blur-lg">
        <div class="p-6 background-3 max-w-md w-full shadow-lg relative">
          <button class="z-button right-2 top-2 absolute" aria-label="Close" @click="showModal = false">
            X
          </button>
          <h2 class="text-xl text-primary font-bold mb-4">
            Add New Playlist
          </h2>
          <form class="space-y-4" @submit.prevent="createNewPlaylist">
            <div>
              <label for="playlist-name" class="text-muted font-medium mb-1">Playlist Name</label>
              <input
                id="playlist-name"
                v-model="newPlaylistName"
                type="text"
                class="px-3 py-2 border w-auto"
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
                <svg xmlns="http://www.w3.org/2000/svg" class="text-primary-500 h-5 w-5 inline" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-7.5 7.5a1 1 0 01-1.414 0l-3.5-3.5a1 1 0 111.414-1.414L8 11.086l6.793-6.793a1 1 0 011.414 0z" clip-rule="evenodd" /></svg>
              </span>
              <span v-else>Add Playlist</span>
            </button>
            <div v-if="submitError" class="text-sm text-red-600 mt-2">
              {{ submitError }}
            </div>
          </form>
        </div>
      </div>
    </teleport>
  </div>
</template>
