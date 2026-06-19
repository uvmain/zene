<script setup lang="ts">
import type { StructuredLyricLine } from '~/types/subsonicLyrics'

const props = defineProps({
  modalTitle: { type: String, default: 'Lyrics' },
  lyrics: { type: Array as PropType<StructuredLyricLine[]>, default: () => [] },
  currentSeconds: { type: Number, default: 0 },
})

defineEmits(['close'])

const scroller = ref<HTMLElement | null>(null)

const currentIndex = computed(() => {
  return props.lyrics.findIndex((line, index) => {
    const nextLineStart = props.lyrics[index + 1]?.start ?? Infinity
    return (line.start ?? 0) <= props.currentSeconds && nextLineStart > props.currentSeconds
  })
})

function setScroller(element: unknown, index: number) {
  if (currentIndex.value === index) {
    scroller.value = element instanceof HTMLElement ? element : null
  }
}

function scrollToCurrentIndex() {
  if (scroller.value !== null) {
    scroller.value.scrollIntoView({ behavior: 'smooth', block: 'center' })
  }
}

watch(currentIndex, () => {
  scrollToCurrentIndex()
}, { flush: 'post' })

onMounted(() => {
  scrollToCurrentIndex()
})
</script>

<template>
  <Modal :show-modal="true" :modal-title="modalTitle" @close="$emit('close')">
    <template #content>
      <div v-if="lyrics.length > 0">
        <div class="text-lg text-muted text-center flex flex-col gap-2 max-h-85dvh min-w-50dvw overflow-y-scroll lg:text-base">
          <div
            v-for="(line, index) in lyrics" :key="index"
            :ref="(element) => setScroller(element, index)"
            :class="{
              'bg-main-400/20': currentIndex === index,
              '': (line.start ?? 0) >= currentSeconds,
            }"
          >
            <div
              class="p-2 flex flex-row gap-2 items-center justify-center"
            >
              <icon-nrk-media-play
                v-if="currentIndex === index"
                class="footer-icon-on"
              />
              <span>{{ line.value }}</span>
            </div>
          </div>
        </div>
      </div>
      <p v-else class="text-center flex flex-col gap-2">
        'No lyrics available.'
      </p>
    </template>
  </Modal>
</template>
