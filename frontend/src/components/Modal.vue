<script setup lang="ts">
defineProps({
  showModal: { type: Boolean, required: true },
  modalText: { type: String, required: true },
})
</script>

<template>
  <teleport to="body">
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center backdrop-blur-lg">
      <div class="relative max-w-md w-full flex flex-col items-center justify-center gap-4 background-3 p-6 text-center align-middle">
        <div class="mb-4 text-lg text-muted font-semibold">
          {{ modalText }}
        </div>
        <slot name="buttons" />
      </div>
    </div>
  </teleport>
</template>
