<script setup lang="ts">
import { useSearch } from '../composables/useSearch'

const props = defineProps({
  genre: { type: String, required: true },
})

const { closeSearch } = useSearch()
const router = useRouter()

function navigateToGenre() {
  closeSearch()
  router.push(`/genres/${props.genre}`)
}
</script>

<template>
  <button class="cursor-pointer rounded-full border-none bg-zene-400 px-3 py-1 text-sm text-white outline-none transition-all duration-200 hover:from-zene-400 hover:to-zene-200 hover:bg-gradient-to-b hover:shadow-sm hover:shadow-zene-400" @click="navigateToGenre()">
    {{ genre }}
  </button>
</template>
