<script setup lang="ts">
import { ref } from 'vue'
import { setRating } from '~/logic/backendFetch'

const props = defineProps({
  musicbrainzId: { type: String, required: true },
})

const model = defineModel<number>({ default: 0 })
const hoverRating = ref(0)

function updateRating(newRating: number) {
  if (model.value === newRating) {
    newRating = 0
  }
  setRating(props.musicbrainzId, newRating)
  model.value = newRating
}

function effectiveRating() {
  return hoverRating.value > 0 ? hoverRating.value : model.value
}
</script>

<template>
  <div
    class="flex flex-row gap-1 cursor-pointer transition-all duration-200 items-center justify-center lg:gap-2"
    @mouseleave="hoverRating = 0"
  >
    <div v-for="rating in 5" :key="rating" class="hover:scale-115" @click="updateRating(rating)" @mouseenter="hoverRating = rating">
      <icon-nrk-star-active
        v-if="effectiveRating() >= rating"
        :class="{
          'text-primary-400': rating === 5,
          'text-primary-400/60': rating === 1,
          'text-primary-400/70': rating === 2,
          'text-primary-400/80': rating === 3,
          'text-primary-400/90': rating === 4,
        }"
      />
      <icon-nrk-star v-else class="text-muted opacity-70" />
    </div>
  </div>
</template>
