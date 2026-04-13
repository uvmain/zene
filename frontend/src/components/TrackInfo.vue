<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { formatTimeFromSeconds, getCoverArtUrl, onImageError } from '~/logic/common'

const props = defineProps({
  track: { type: Object as PropType<SubsonicSong>, required: true },
})

const showModal = defineModel<boolean>({
  default: false,
})

const router = useRouter()

const modelTrack = computed(() => {
  return props.track
})

const coverArtUrl = computed(() => {
  return getCoverArtUrl(props.track.musicBrainzId)
})

function navigateAlbum() {
  showModal.value = false
  router.push(`/albums/${modelTrack.value.albumId}`)
}
</script>

<template>
  <Modal :show-modal="showModal" @close="showModal = false">
    <template #content>
      <div class="m-auto p-4 corner-cut background-1 lg:p-8">
        <div class="flex flex-col gap-8 lg:flex-row">
          <img
            v-if="coverArtUrl"
            :src="coverArtUrl"
            alt="Album Art"
            class="mx-auto h-auto max-w-30vw w-full cursor-pointer shadow-lg"
            @error="onImageError"
            @click="navigateAlbum"
          />
          <div class="flex flex-col gap-4 lg:w-2/3">
            <h1 class="text-3xl text-primary font-bold mb-2">
              {{ track.title }}
            </h1>
            <RouterLink
              class="text-xl text-muted mb-1 no-underline cursor-pointer hover:underline hover:underline-white"
              :to="`/artists/${track.artistId}`"
            >
              Artist: {{ track.artist }}
            </RouterLink>
            <RouterLink
              class="text-lg text-muted mb-1 no-underline cursor-pointer hover:underline hover:underline-white"
              :to="`/albums/${track.albumId}`"
            >
              Album: {{ track.album }}
            </RouterLink>
            <p class="text-muted mb-1">
              Duration: {{ formatTimeFromSeconds(track.duration) }}
            </p>
            <p v-if="track" class="text-muted mb-4">
              Released: {{ track.year }}
            </p>
            <div class="flex flex-row gap-4 items-center justify-center lg:gap-8">
              <PlayButton class="flex justify-start" :track="modelTrack" />
              <Fave v-model="modelTrack.starred" :musicbrainz-id="modelTrack.musicBrainzId" />
              <Rating v-model="modelTrack.userRating" :musicbrainz-id="modelTrack.musicBrainzId" />
            </div>
          </div>
        </div>
      </div>
    </template>
  </Modal>
</template>

<style lang="css" scoped>
dialog::backdrop {
  @apply backdrop-blur-lg;
}
</style>
