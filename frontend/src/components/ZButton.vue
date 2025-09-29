<script setup lang="ts">
defineEmits(['click'])
</script>

<template>
  <div class="bg-zshade-600 corner-cut mx-auto flex items-center align-middle dark:bg-zshade-100">
    <button class="z-button m-1px" @click="$emit('click')">
      <slot />
    </button>
  </div>
</template>
