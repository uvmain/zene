<script setup lang="ts">
import type { SubsonicPodcastChannelsResponse, SubsonicResponse } from '~/types/subsonic'
import type { SubsonicPodcastChannel } from '~/types/subsonicPodcasts'
import { openSubsonicFetchRequest } from '~/logic/backendFetch'

const router = useRouter()

const showAddPodcastModal = ref(false)
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
  setTimeout(() => {
    showAddPodcastModal.value = false
    showSuccess.value = false
    isSubmitting.value = false
  }, 500)
  await getPodcasts()
}

async function getPodcasts() {
  const formData = new FormData()
  formData.append('includeEpisodes', true.toString())
  const response = await openSubsonicFetchRequest<SubsonicPodcastChannelsResponse>('getPodcasts', {
    body: formData,
  })
  podcasts.value = response?.podcasts?.channel ?? []
}

function navigateToPodcast(podcastId: string) {
  router.push(`/podcasts/${podcastId}`)
}

onBeforeMount(getPodcasts)
</script>

<template>
  <div class="mx-auto p-4 lg:max-w-60dvw space-y-4">
    <button class="z-button" @click="showAddPodcastModal = true">
      Add New Podcast Channel
    </button>

    <div class="mt-8">
      <div v-if="podcasts.length === 0" class="text-primary">
        No podcasts found.
      </div>
      <div class="flex flex-wrap gap-6">
        <HeroPodcast
          v-for="podcast in podcasts"
          :key="podcast.id"
          :podcast="podcast"
          @click="navigateToPodcast(podcast.id)"
        />
      </div>
    </div>

    <Modal
      :show-modal="showAddPodcastModal"
      modal-title="Add New Podcast Channel"
      @close="showAddPodcastModal = false; submitError = ''; newPodcastUrl = ''"
    >
      <template #content>
        <form class="space-y-4" @submit.prevent="createNewPodcast">
          <div>
            <label for="stream-url" class="mb-1 block text-muted font-medium">Stream URL</label>
            <input
              id="stream-url"
              v-model="newPodcastUrl"
              type="text"
              class="w-auto border px-3 py-2"
              placeholder="Enter stream URL"
              required
            />
          </div>
          <div v-if="submitError" class="mt-2 text-sm text-red-600">
            {{ submitError }}
          </div>
        </form>
        <ZButton
          aria-label="Add Podcast"
          :disabled="isSubmitting || showSuccess"
          @click="createNewPodcast()"
        >
          <span v-if="isSubmitting && !showSuccess">Adding...</span>
          <span v-else-if="showSuccess">
            <svg xmlns="http://www.w3.org/2000/svg" class="inline h-5 w-5 text-primary1" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-7.5 7.5a1 1 0 01-1.414 0l-3.5-3.5a1 1 0 111.414-1.414L8 11.086l6.793-6.793a1 1 0 011.414 0z" clip-rule="evenodd" /></svg>
          </span>
          <span v-else>Add Podcast</span>
        </ZButton>
      </template>
    </Modal>
  </div>
</template>
