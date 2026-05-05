<script setup lang="ts">
defineProps({
  primary: { type: Boolean, default: false },
  size10: { type: Boolean, default: false },
  size12: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
  hoverText: { type: String, required: false },
  red: { type: Boolean, default: false },
  green: { type: Boolean, default: false },
})

defineEmits(['click'])
</script>

<template>
  <abbr :title="hoverText" class="align-middle decoration-none flex items-center">
    <button
      class="group z-button align-middle inline-flex items-center relative"
      :disabled="disabled"
      :class="{
        'size-10': size10,
        'size-12': size12,
        'cursor-not-allowed opacity-50': disabled,
        'border-primary-400!': primary,
        'border-red-400/70!': red,
        'border-green-400/70!': green,
      }"
      @click="$emit('click')"
    >
      <slot />
    </button>
  </abbr>
</template>
