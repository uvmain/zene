<script setup lang="ts">
const props = defineProps({
  trackUrl: { type: String, required: false, default: '' },
})

const emits = defineEmits(['play', 'pause', 'timeUpdate', 'ended'])

const audioRef = ref<HTMLAudioElement | null>(null)

defineExpose({ audioRef })

watch(
  () => props.trackUrl,
  (newTrack, oldTrack) => {
    const audio = audioRef.value
    if (!audio || newTrack === oldTrack) {
      return
    }
    if (audio) {
      audio.addEventListener(
        'canplaythrough',
        () => {
          audio?.play()
        },
        { once: true },
      )
    }
  },
)

onMounted(() => {
  const audio = audioRef.value
  if (!audio) {
    return
  }
  audio.addEventListener('play', () => emits('play'))
  audio.addEventListener('pause', () => emits('pause'))
  audio.addEventListener('timeupdate', () => emits('timeUpdate', audio.currentTime))
  audio.addEventListener('ended', () => emits('ended'))
})

onUnmounted(() => {
  const audio = audioRef.value
  if (audio) {
    audio.removeEventListener('play', () => emits('play'))
    audio.removeEventListener('pause', () => emits('pause'))
    audio.removeEventListener('timeupdate', () => emits('timeUpdate', audio.currentTime))
    audio.removeEventListener('ended', () => emits('ended'))

    audio.pause()
    audio.removeAttribute('src')
    audio.load()
  }
})
</script>

<template>
  <audio ref="audioRef" :src="trackUrl" preload="metadata" class="hidden" />
</template>
