<script setup lang="ts">
import type { SlugWithDimensions } from '../types/main'
import { useElementVisibility } from '@vueuse/core'
import { backendFetchRequest } from '../composables/fetchFromBackend'
import { getThumbnailPath } from '../composables/logic'

const router = useRouter()
const startObserver = ref<HTMLDivElement | null>(null)
const slugs = ref<SlugWithDimensions[]>([])

const startObserverIsVisible = useElementVisibility(startObserver)

async function getSlugs() {
  try {
    const response = await backendFetchRequest('/slugs/with-dimensions')
    slugs.value = await response.json() as SlugWithDimensions[]
  }
  catch (error) {
    console.error('Failed to fetch thumbnails:', error)
    slugs.value = []
  }
}

function navigateToSlug(slug: string) {
  const slugPath = `/${slug}`
  router.push(slugPath)
}

const headerShadowClass = computed(() => {
  return startObserverIsVisible.value ? ' ' : 'shadow-lg'
})

onBeforeMount(async () => {
  await getSlugs()
})
</script>

<template>
  <div class="min-h-screen">
    <Header class="sticky top-0 z-10" :class="headerShadowClass" />
    <div class="flex flex-col items-center p-6">
      <div ref="startObserver" />
      <div class="flex flex-col gap-4 lg:max-w-8/10 lg:flex-row lg:flex-wrap lg:gap-x-2 lg:gap-y-1">
        <div v-for="(slug, index) in slugs" :key="index" class="flex-1 basis-auto">
          <img :src="getThumbnailPath(slug.slug)" :alt="slug.slug" loading="lazy" :width="slug.width" :height="slug.height" class="h-full min-h-20vh w-full cursor-pointer object-cover lg:max-h-25vh lg:max-w-40vw" @click="navigateToSlug(slug.slug)">
        </div>
        <div class="flex grow-2" />
      </div>
    </div>
  </div>
</template>
