<script setup lang="ts">
import { onKeyStroke } from '@vueuse/core'

defineProps({
  showModal: { type: Boolean, required: true },
  modalTitle: { type: String, required: false },
})

const emits = defineEmits(['close'])

const modalContainer = useTemplateRef('modal-container')

onKeyStroke('Escape', (e) => {
  e.preventDefault()
  emits('close')
})

function handleClickOutside(event: MouseEvent) {
  if (modalContainer.value && !modalContainer.value.contains(event.target as Node)) {
    emits('close')
  }
}

onMounted(() => document.addEventListener('mousedown', handleClickOutside))
onUnmounted(() => document.removeEventListener('mousedown', handleClickOutside))
</script>

<template>
  <teleport to="body">
    <div v-if="showModal" class="p-2 bg-gray/5 flex items-center inset-0 justify-center fixed z-50 backdrop-blur-lg lg:p-4">
      <div ref="modal-container" class="m-auto p-6 text-center align-middle border-muted corner-cut-large background-1 flex flex-col gap-4 items-center justify-center relative">
        <div v-if="modalTitle" class="text-lg text-muted font-semibold mb-4 max-w-80%">
          {{ modalTitle }}
        </div>
        <div class="opacity-50 right-1 top-1 absolute hover:opacity-100 lg:(right-2 top-2)">
          <ZButton title="Close modal" size-10 @click="$emit('close')">
            <icon-nrk-close class="text-primary size-full" />
          </ZButton>
        </div>
        <slot name="content" />
      </div>
    </div>
  </teleport>
</template>
