<script setup lang="ts">
const emit = defineEmits(['update:modelValue'])
const image = ref<string | null>(null)
const imageUrlInput = ref<string>('')

function onFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    const reader = new FileReader()
    reader.onload = () => {
      image.value = reader.result as string
      imageUrlInput.value = ''
    }
    reader.readAsDataURL(file)
  }
}

async function onUrlChange() {
  try {
    const response = await fetch(imageUrlInput.value)
    if (!response.ok)
      throw new Error('Network response was not ok')
    const blob = await response.blob()
    const reader = new FileReader()
    reader.onloadend = () => {
      image.value = reader.result as string
    }
    reader.readAsDataURL(blob)
  }
  catch (error) {
    console.error(`Failed to load image: ${error}`)
  }
}

function onPaste(event: ClipboardEvent) {
  const items = event.clipboardData?.items
  if (items) {
    for (const item of Array.from(items)) {
      if (item.type.startsWith('image/')) {
        const file = item.getAsFile()
        if (file) {
          const reader = new FileReader()
          reader.onload = () => {
            image.value = reader.result as string
            imageUrlInput.value = ''
          }
          reader.readAsDataURL(file)
        }
      }
    }
  }
}

watch(image, () => {
  emit('update:modelValue', image.value)
})
</script>

<template>
  <div class="text max-w-100vw flex flex-col items-center gap-2 p-6" @paste="onPaste">
    <input
      id="fileInput"
      type="file"
      accept="image/*"
      class="hidden"
      @change="onFileChange"
    />
    <label
      for="fileInput"
      class="block w-28rem cursor-pointer border-1 border-gray-300 rounded-lg border-solid bg-gray-50 p-2 text-center text-sm text-gray-700 hover:bg-gray-200 focus:outline-none"
    >
      Browse
    </label>
    <input
      v-model="imageUrlInput"
      placeholder="Enter image URL"
      class="mt-2 block w-28rem border-1 border-gray-300 rounded-lg border-solid p-2 focus:outline-none focus:ring-2 focus:ring-blue-600"
      @change="onUrlChange"
    />
    <textarea
      placeholder="Paste an image here"
      class="mt-2 block w-28rem border-1 border-gray-300 rounded-lg border-solid p-2 focus:outline-none focus:ring-2 focus:ring-blue-600"
    />
  </div>
</template>
