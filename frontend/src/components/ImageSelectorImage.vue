<script setup lang="ts">
const props = defineProps({
  imageUrl: { type: String, required: true },
  label: { type: String, required: true },
  type: { type: String, required: true },
})

defineEmits(['updateArt'])
</script>

<template>
  <div class="relative w-56">
    <img
      class="size-56"
      :src="imageUrl"
      :alt="label"
    />
    <ZButton
      class="absolute bottom-10 right-2"
      aria-label="Choose art"
      @click="$emit('updateArt', props.type)"
    >
      Use This Art
    </ZButton>
    <div class="mt-2 text-primary">
      {{ label }}
    </div>
  </div>
</template>
