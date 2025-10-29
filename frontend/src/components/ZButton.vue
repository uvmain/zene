<script setup lang="ts">
defineProps({
  primary: { type: Boolean, default: false },
  size12: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
})

defineEmits(['click'])
</script>

<template>
  <div
    class="corner-cut inline-flex items-center align-middle hover:bg-primary1"
    :disabled="disabled"
    :class="{
      'cursor-not-allowed opacity-50': disabled,
      'bg-primary2': primary,
      'bg-zshade-600 dark:bg-zshade-400': !primary,
    }"
  >
    <button
      class="group/button z-button m-1px"
      :class="{
        'size-12': size12,
      }"
      :disabled="disabled"
      @click="$emit('click')"
    >
      <slot />
    </button>
  </div>
</template>
