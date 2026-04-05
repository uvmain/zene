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
  <div class="text p-6 flex flex-col gap-2 max-w-100vw items-center" @paste="onPaste">
    <input
      id="fileInput"
      type="file"
      accept="image/*"
      class="hidden"
      @change="onFileChange"
    />
    <label
      for="fileInput"
      class="text-sm text-muted p-2 text-center border-1 border-gray-300 corner-cut border-solid bg-zshade-200 w-28rem cursor-pointer focus:outline-none dark-border-zshade-700 dark:bg-zshade-600"
    >
      Browse
    </label>
    <input
      v-model="imageUrlInput"
      placeholder="Enter image URL"
      class="mt-2 p-2 border-1 border-gray-300 border-solid w-28rem focus:outline-none focus:ring-2 focus:ring-primary1"
      @change="onUrlChange"
    />
    <textarea
      placeholder="Paste an image here"
      class="mt-2 p-2 border-1 border-gray-300 border-solid w-28rem focus:outline-none focus:ring-2 focus:ring-primary1"
    />
  </div>
</template>
