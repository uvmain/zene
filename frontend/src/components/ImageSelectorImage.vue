<script setup lang="ts">
const props = defineProps({
  imageUrl: { type: String, required: true },
  label: { type: String, required: true },
  type: { type: String, required: true },
})

defineEmits(['updateArt'])

const loaded = ref(false)
</script>

<template>
  <div class="relative w-56">
    <img
      class="aspect-square size-56 rounded-md bg-bluegray"
      :src="imageUrl"
      :alt="label"
      width="224"
      height="224"
      @load="loaded = true"
    />
    <Loading v-if="!loaded" class="absolute left-1/2 top-1/2 size-56 translate-x--1/2 translate-y--1/2" />
    <ZButton
      class="absolute bottom-10 right-2"
      aria-label="Choose art"
      @click="$emit('updateArt', props.type)"
    >
      Use This Art
    </ZButton>
    <div class="mt-2 text-muted">
      {{ label }}
    </div>
  </div>
</template>
