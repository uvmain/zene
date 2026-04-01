<script setup lang="ts">
interface Option {
  label: string
  emitValue: string
}

defineProps({
  options: { type: Array as PropType<Option[]>, required: true },
})

const emits = defineEmits(['setOrder'])
</script>

<template>
  <div class="corner-cut background-2 w-auto left-0 top-0 absolute z-10">
    <div
      v-for="option in options"
      :key="option.label"
      class="px-4 py-2 cursor-pointer hover:background-3"
      @click="emits('setOrder', option.emitValue)"
    >
      {{ option.label }}
    </div>
  </div>
</template>
