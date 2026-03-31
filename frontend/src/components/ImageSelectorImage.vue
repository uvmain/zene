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
  <div class="w-56 relative">
    <img
      class="bg-bluegray rounded-md size-56 aspect-square"
      :src="imageUrl"
      :alt="label"
      width="224"
      height="224"
      @load="loaded = true"
    />
    <Loading v-if="!loaded" class="size-56 translate-x--1/2 translate-y--1/2 left-1/2 top-1/2 absolute" />
    <ZButton
      class="opacity-70 bottom-10 right-2 absolute hover:opacity-100"
      aria-label="Choose art"
      @click="$emit('updateArt', props.type)"
    >
      Use This Art
    </ZButton>
    <div class="text-muted mt-2">
      {{ label }}
    </div>
  </div>
</template>
