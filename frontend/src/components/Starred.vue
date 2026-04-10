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
  <div class="flex cursor-pointer items-center justify-center" @click="toggleStarred" @click.stop>
    <icon-nrk-star-active v-if="model" class="text-primary-400" />
    <icon-nrk-star v-else class="text-muted opacity-70 hover:opacity-100" />
  </div>
</template>
