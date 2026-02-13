<script setup lang="ts">
import type { StructuredLyricLine } from '~/types/subsonicLyrics'

defineProps({
  lyrics: { type: Array as PropType<StructuredLyricLine[]>, default: () => [] },
  currentSeconds: { type: Number, default: 0 },
})

defineEmits(['close'])
</script>

<template>
  <div>
    <teleport to="body">
      <div class="fixed inset-0 z-50 overflow-y-scroll bg-white/40 p-4 backdrop-blur-xl dark:bg-black/40">
        <div class="flex justify-center">
          <div v-if="lyrics.length > 0">
            <div class="flex flex-col gap-2 text-center text-lg text-muted lg:text-base">
              <div
                v-for="(line, index) in lyrics" :key="index"
                :class="{
                  'bg-green/20': line.start && line.start <= currentSeconds && (lyrics[index + 1]?.start ?? 0) > currentSeconds,
                  '': (line.start ?? 0) >= currentSeconds,
                }"
              >
                {{ (line.start ?? 0) <= currentSeconds && (lyrics[index + 1]?.start ?? 0) > currentSeconds
                  ? '▶️ '
                  : '' }}
                <span>{{ line.value }}</span>
              </div>
            </div>
          </div>
          <p v-else class="flex flex-col gap-2 text-center">
            'No lyrics available.'
          </p>
        </div>
      </div>
    </teleport>
    <teleport to="body">
      <div class="absolute right-4 top-4 z-51 lg:right-8">
        <ZButton class="" @click="$emit('close')">
          Close
        </ZButton>
      </div>
    </teleport>
  </div>
</template>
