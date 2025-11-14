<script setup lang="ts">
interface Option {
  label: string
  emitName: string
}

defineProps({
  options: { type: Array as PropType<Option[]>, required: true },
})

const emits = defineEmits(['setOrder'])
</script>

<template>
  <div class="corner-cut absolute left-0 top-0 z-10 w-auto background-2">
    <div
      v-for="option in options"
      :key="option.label"
      class="cursor-pointer px-4 py-2 hover:background-3"
      @click="emits('setOrder', option.emitName)"
    >
      {{ option.label }}
    </div>
  </div>
</template>
