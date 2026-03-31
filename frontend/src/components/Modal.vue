<script setup lang="ts">
defineProps({
  showModal: { type: Boolean, required: true },
  modalTitle: { type: String, required: true },
})

defineEmits(['close'])
</script>

<template>
  <teleport to="body">
    <div v-if="showModal" class="bg-gray/5 flex items-center inset-0 justify-center fixed z-50 backdrop-blur-lg">
      <div class="p-6 text-center align-middle border-muted corner-cut-large background-1 flex flex-col gap-4 max-w-lg w-full items-center justify-center relative">
        <div class="text-lg text-muted font-semibold mb-4 max-w-80%">
          {{ modalTitle }}
        </div>
        <ZButton :size10="true" aria-label="Close" hover-text="Close" class="right-4 top-4 absolute" @click="$emit('close')">
          <icon-nrk-close class="text-primary size-full" />
        </ZButton>
        <slot name="content" />
      </div>
    </div>
  </teleport>
</template>
