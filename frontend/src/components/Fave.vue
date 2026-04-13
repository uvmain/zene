<script setup lang="ts">
import { postStarToggle } from '~/logic/backendFetch'

const props = defineProps({
  musicbrainzId: { type: String, required: true },
})

const model = defineModel<string | undefined>()

function toggleStarred() {
  if (model.value) {
    postStarToggle(props.musicbrainzId, false)
    model.value = undefined
  }
  else {
    postStarToggle(props.musicbrainzId, true)
    model.value = new Date().toDateString()
  }
}
</script>

<template>
  <div class="flex cursor-pointer transition-all duration-200 items-center justify-center hover:scale-115" @click="toggleStarred" @click.stop>
    <icon-nrk-heart-active v-if="model" class="text-primary-400" />
    <icon-nrk-heart v-else class="text-muted opacity-70" />
  </div>
</template>
