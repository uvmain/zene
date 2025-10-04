<script setup lang="ts">
import type { SubsonicAlbum } from '~/types/subsonicAlbum'
import { fetchAlbums } from '~/composables/backendFetch'
import { getCoverArtUrl } from '~/composables/logic'

const METADATA_COUNT = 20
const isShaking = ref(false)
const albumArray = ref<SubsonicAlbum[]>([])
const index = ref(0)

const indexCount = computed(() => {
  return albumArray.value.length < METADATA_COUNT ? albumArray.value.length : METADATA_COUNT
})

function nextIndex() {
  if (index.value < indexCount.value - 1) {
    index.value += 1
  }
  else {
    index.value = 0
  }
}

function prevIndex() {
  if (index.value > 0) {
    index.value -= 1
  }
  else {
    index.value = indexCount.value - 1
  }
}

async function getRandomAlbums(limit: number) {
  albumArray.value = await fetchAlbums('random', limit, 0)
  index.value = 0
}

const coverArtUrl = computed(() => {
  return getCoverArtUrl(albumArray.value[index.value].coverArt, 600)
})

function handleDiceClick() {
  isShaking.value = true
  setTimeout(() => {
    isShaking.value = false
  }, 200)
  getRandomAlbums(METADATA_COUNT)
}

onBeforeMount(async () => {
  await getRandomAlbums(METADATA_COUNT)
  index.value = 0
})
</script>

<template>
  <section v-if="albumArray.length" class="corner-cut-large overflow-hidden">
    <div
      class="h-full w-full bg-cover bg-center"
      :style="{ backgroundImage: `url(${coverArtUrl})` }"
    >
      <div class="h-full w-full flex items-center justify-between background-grad-2 backdrop-blur-md">
        <Album :album="albumArray[index]" size="md" />
        <div class="corner-cut m-3 mb-auto flex gap-2 background-2 p-3 md:m-6 md:mb-auto md:p-2">
          <icon-nrk-chevron-left
            class="cursor-pointer text-2xl opacity-80 md:text-3xl hover:text-primary2 active:opacity-100"
            @click="prevIndex"
          />
          <icon-nrk-dice-3
            class="cursor-pointer text-2xl opacity-80 md:text-3xl hover:text-primary2 active:opacity-100"
            :class="{ shake: isShaking }"
            @click="handleDiceClick()"
          />
          <icon-nrk-chevron-right
            class="cursor-pointer text-2xl opacity-80 md:text-3xl hover:text-primary2 active:opacity-0"
            @click="nextIndex"
          />
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
@keyframes shake {
  0% {
    transform: rotate(0deg);
  }
  25% {
    transform: rotate(-15deg);
  }
  50% {
    transform: rotate(15deg);
  }
  75% {
    transform: rotate(-15deg);
  }
  100% {
    transform: rotate(0deg);
  }
}

.shake {
  animation: shake 0.2s ease-in-out;
}
</style>
