<script setup lang="ts">
import type { SubsonicSong } from '~/types/subsonicSong'
import { fetchShare } from '~/logic/shareFetch'
import { shareToken } from '~/stores/main'

const route = useRoute('/shares/[token]')
const shareTokenParam = computed(() => route.params.token as string)

const tracks = ref<SubsonicSong[]>([])
const description = ref<string>('')

onBeforeMount(async () => {
  console.log('Fetching share details for token:', shareTokenParam.value)
  shareToken.value = shareTokenParam.value
  const response = await fetchShare(shareToken.value)
  if (response) {
    description.value = response.description || ''
    tracks.value = response.entry || []
  }
})

onUnmounted(() => {
  shareToken.value = ''
})
</script>

<template>
  <div class="mx-auto p-4 container">
    <div v-if="description" class="mb-4 text-center text-lg text-muted">
      {{ description }}
    </div>
    <Tracks v-if="tracks" :tracks="tracks" :is-share="true" />
  </div>
</template>
