<script setup lang="ts">
import type { SubsonicRadioStationsResponse, SubsonicResponse } from '~/types/subsonic'
import type { SubsonicRadioStation } from '~/types/subsonicRadioStations'
import { openSubsonicFetchRequest } from '~/composables/backendFetch'

const showModal = ref(false)
const isSubmitting = ref(false)
const submitError = ref('')
const showSuccess = ref(false)

const newStreamUrl = ref('')
const newStreamName = ref('')
const newStreamHomepageUrl = ref('')

const radioStations = ref<SubsonicRadioStation[]>([])

async function createNewRadioStation() {
  if (isSubmitting.value)
    return
  isSubmitting.value = true
  submitError.value = ''
  showSuccess.value = false
  const formData = new FormData()
  formData.append('streamUrl', newStreamUrl.value)
  formData.append('name', newStreamName.value)
  formData.append('homepageUrl', newStreamHomepageUrl.value)
  const response = await openSubsonicFetchRequest<SubsonicResponse>('createInternetRadioStation', {
    body: formData,
  })
  if (response?.error?.message) {
    submitError.value = response.error.message
    isSubmitting.value = false
    return
  }
  showSuccess.value = true
  getRadioStations()
  setTimeout(() => {
    showModal.value = false
    showSuccess.value = false
    isSubmitting.value = false
  }, 1000)
}

async function getRadioStations() {
  const response = await openSubsonicFetchRequest<SubsonicRadioStationsResponse>('getInternetRadioStations')
  radioStations.value = response?.internetRadioStations?.internetRadioStation
}

onBeforeMount(getRadioStations)
</script>

<template>
  <div class="p-4 space-y-4">
    <button class="z-button" @click="showModal = true">
      Add New Radio Station
    </button>

    <div class="mt-8">
      <h2 class="mb-4 text-lg font-bold">
        Internet Radio Stations
      </h2>
      <div v-if="radioStations.length === 0" class="text-primary">
        No radio stations found.
      </div>
      <ul v-else class="space-y-4">
        <li v-for="station in radioStations" :key="station.id" class="border background-2 p-4 shadow">
          <div class="text-lg text-primary font-semibold">
            {{ station.name }}
          </div>
          <div class="mt-1 text-sm">
            <span class="font-medium">Stream URL: </span>
            <a :href="station.streamUrl" target="_blank" class="text-muted hover:underline">{{ station.streamUrl }}</a>
          </div>
          <div class="mt-1 text-sm">
            <span class="font-medium">Homepage: </span>
            <a :href="station.homepageUrl" target="_blank" class="text-muted hover:underline">{{ station.homepageUrl }}</a>
          </div>
        </li>
      </ul>
    </div>

    <teleport to="body">
      <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center backdrop-blur-lg">
        <div class="relative max-w-md w-full background-3 p-6 shadow-lg">
          <button class="z-button absolute right-2 top-2" aria-label="Close" @click="showModal = false">
            X
          </button>
          <h2 class="mb-4 text-xl text-primary font-bold">
            Add New Radio Station
          </h2>
          <form class="space-y-4" @submit.prevent="createNewRadioStation">
            <div>
              <label for="stream-url" class="mb-1 block text-muted font-medium">Stream URL</label>
              <input id="stream-url" v-model="newStreamUrl" type="text" class="w-auto border px-3 py-2" placeholder="Enter stream URL" required />
            </div>
            <div>
              <label for="stream-name" class="mb-1 block text-muted font-medium">Station Name</label>
              <input id="stream-name" v-model="newStreamName" type="text" class="w-auto border px-3 py-2" placeholder="Enter station name" required />
            </div>
            <div>
              <label for="homepage-url" class="mb-1 block text-muted font-medium">Homepage URL</label>
              <input id="homepage-url" v-model="newStreamHomepageUrl" type="text" class="w-auto border px-3 py-2" placeholder="Enter homepage URL" />
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
