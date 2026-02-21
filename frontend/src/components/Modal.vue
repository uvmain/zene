<script setup lang="ts">
defineProps({
  showModal: { type: Boolean, required: true },
  modalTitle: { type: String, required: true },
})

defineEmits(['close'])
</script>

<template>
  <teleport to="body">
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-gray/5 backdrop-blur-lg">
      <div class="corner-cut-large relative max-w-lg w-full flex flex-col items-center justify-center gap-4 border-muted background-1 p-6 text-center align-middle">
        <div class="mb-4 max-w-80% text-lg text-muted font-semibold">
          {{ modalTitle }}
        </div>
        <ZButton :size10="true" aria-label="Close" hover-text="Close" class="absolute right-4 top-4" @click="$emit('close')">
          <icon-nrk-close class="size-full text-primary" />
        </ZButton>
        <slot name="content" />
      </div>
    </div>
  </teleport>
</template>
