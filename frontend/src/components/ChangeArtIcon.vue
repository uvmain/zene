<script setup lang="ts">
defineEmits(['click'])
</script>

<template>
  <div title="Change art" class="flex cursor-pointer transition-all duration-200 items-center justify-center hover:scale-115" @click="$emit('click')" @click.stop>
    <icon-nrk-gallery />
  </div>
</template>
