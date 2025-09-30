<script setup lang="ts">
defineProps({
  primary: { type: Boolean, default: false },
  size12: { type: Boolean, default: false },
})

defineEmits(['click'])
</script>

<template>
  <div
    class="corner-cut inline-flex items-center align-middle hover:bg-primary1"
    :class="{
      'bg-primary2': primary,
      'bg-zshade-600 dark:bg-zshade-100': !primary,
    }"
  >
    <button
      class="z-button m-1px"
      :class="{
        'size-12': size12,
      }"
      @click="$emit('click')"
    >
      <slot />
    </button>
  </div>
</template>
