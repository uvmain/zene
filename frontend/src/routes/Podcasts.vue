<script setup lang="ts">
import type { SubsonicPodcastChannelsResponse, SubsonicResponse } from '~/types/subsonic'
import type { SubsonicPodcastChannel } from '~/types/subsonicPodcasts'
import { openSubsonicFetchRequest } from '~/composables/backendFetch'

const router = useRouter()

const showModal = ref(false)
const isSubmitting = ref(false)
const submitError = ref('')
const showSuccess = ref(false)

const newPodcastUrl = ref('')

const podcasts = ref<SubsonicPodcastChannel[]>([])

async function createNewPodcast() {
  if (isSubmitting.value)
    return
  isSubmitting.value = true
  submitError.value = ''
  showSuccess.value = false
  const formData = new FormData()
  formData.append('url', newPodcastUrl.value)
  const response = await openSubsonicFetchRequest<SubsonicResponse>('createPodcastChannel', {
    body: formData,
  })
  if (response?.error?.message) {
    submitError.value = response.error.message
    isSubmitting.value = false
    return
  }
  showSuccess.value = true
  getPodcasts()
  setTimeout(() => {
    showModal.value = false
    showSuccess.value = false
    isSubmitting.value = false
  }, 1000)
}

async function getPodcasts() {
  const formData = new FormData()
  formData.append('includeEpisodes', true.toString())
  const response = await openSubsonicFetchRequest<SubsonicPodcastChannelsResponse>('getPodcasts', {
    body: formData,
  })
  podcasts.value = response?.podcasts?.channel
  podcasts.value.forEach((podcast) => {
    podcast.coverArt = `/share/img/${podcast.coverArt}?size=400`
  })
}

function navigateToPodcast(podcastId: string) {
  router.push(`/podcasts/${podcastId}`)
}

onBeforeMount(getPodcasts)
</script>

<template>
  <div class="p-4 space-y-4">
    <button class="z-button" @click="showModal = true">
      Add New Podcast Channel
    </button>

    <div class="mt-8">
      <div v-if="podcasts.length === 0" class="text-zgray-200">
        No podcasts found.
      </div>
      <div class="flex flex-wrap justify-center gap-6 md:justify-start">
        <div
          v-for="podcast in podcasts"
          :key="podcast.id" class="mx-auto max-w-60dvw flex flex-row cursor-pointer gap-4 align-top transition duration-150 hover:scale-101"
          @click="navigateToPodcast(podcast.id)"
        >
          <img
            :src="podcast.coverArt"
            alt="Podcast Cover"
            class="size-50 object-cover"
          />
          <div class="my-auto flex flex-col gap-4">
            <div class="text-2xl font-bold">
              {{ podcast.title }}
            </div>
            <div
              class="line-clamp-5 max-h-70 overflow-hidden text-ellipsis whitespace-pre-line text-pretty text-op-80"
              v-html="podcast.description.replaceAll(/\n/g, '<br>')"
            />
            <div v-if="podcast.episode[0].genres?.length > 0" class="flex flex-wrap justify-center gap-2 md:justify-start">
              <GenreBottle v-for="genre in podcast.episode[0].genres" :key="genre.name" :genre="genre.name" />
            </div>
          </div>
        </div>
      </div>
    </div>

    <teleport to="body">
      <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-60">
        <div class="relative max-w-md w-full bg-zgray-800 p-6">
          <button class="z-button" aria-label="Close" @click="showModal = false">
            X
          </button>
          <h2 class="mb-4 text-xl text-zgray-200 font-bold">
            Add New Podcast
          </h2>
          <form class="space-y-4" @submit.prevent="createNewPodcast">
            <div>
              <label for="stream-url" class="mb-1 block text-zgray-200 font-medium">Stream URL</label>
              <input
                id="stream-url"
                v-model="newPodcastUrl"
                type="text"
                class="w-auto border px-3 py-2"
                placeholder="Enter stream URL"
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
