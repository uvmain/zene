<script setup lang="ts">
import { useSearch } from '~/composables/useSearch'

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
  <ZButton
    @click="navigateToGenre()"
  >
    <span class="text-nowrap">{{ genre }}</span>
  </ZButton>
</template>
