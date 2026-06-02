<script setup lang="ts">
type Alignment = 'left' | 'right'

defineProps({
  title: { type: String, required: true },
  options: { type: Array as PropType<(string | number)[]>, required: true },
  align: { type: String as PropType<Alignment>, default: 'left' },
  currentOption: { type: [String, Number], required: false },
})

const emits = defineEmits(['select'])

const show = ref<boolean>(false)
const optionsContainer = useTemplateRef('optionsContainer')

function handleClickOutside(event: MouseEvent) {
  if (optionsContainer.value && !optionsContainer.value.contains(event.target as Node)) {
    show.value = false
  }
}

function handleSelect(optionValue: string | number) {
  show.value = false
  emits('select', optionValue)
}

onMounted(() => document.addEventListener('mousedown', handleClickOutside))
onUnmounted(() => document.removeEventListener('mousedown', handleClickOutside))
</script>

<template>
  <div class="w-fit relative">
    <ZButton @click="show = !show">
      <div class="flex flex-row gap-2 items-center">
        <span class="text-sm">{{ title }}</span>
        <icon-nrk-chevron-down />
      </div>
    </ZButton>
    <div
      v-if="show"
      ref="optionsContainer"
      class="mt-1 border-muted background-2 flex-col shadow absolute z-5 overflow-hidden"
      :class="align === 'right' ? 'right-0' : 'left-0'"
    >
      <div
        v-for="item in options"
        :key="item"
        class="p-2 border-l-4 block cursor-pointer hover:bg-main-500/50"
        :class="{
          'border-transparent': currentOption !== item,
          'border-main-600': currentOption === item,
        }"
        @click="handleSelect(item)"
      >
        {{ item }}
      </div>
    </div>
  </div>
</template>
