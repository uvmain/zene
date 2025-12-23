<script setup lang="ts">
defineProps({
  primary: { type: Boolean, default: false },
  size12: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
  hoverText: { type: String, required: false },
})

defineEmits(['click'])
</script>

<template>
  <button
    class="group z-button relative inline-flex items-center align-middle"
    :disabled="disabled"
    :class="{
      'size-12': size12,
      'cursor-not-allowed opacity-50': disabled,
      'border-primary2': primary,
      'border-zshade-600 dark:border-zshade-400': !primary,
    }"
    @click="$emit('click')"
  >
    <slot />
    <div
      v-if="hoverText"
      class="absolute bottom-13 z-100 rounded background-2 p-4px text-primary opacity-0 transition-all duration-200 group-hover:opacity-100"
    >
      {{ hoverText }}
    </div>
  </button>
</template>
