<script setup lang="ts">
import type { StreamQuality } from '~/logic/store'
import { useDark, useToggle } from '@vueuse/core'
import { openSubsonicFetchRequest } from '~/logic/backendFetch'
import { clearApiKey } from '~/logic/common'
import { toggleDebug } from '~/logic/logger'
import { debugEnabled, streamQualities, streamQuality } from '~/logic/store'

const isDark = useDark()
const toggleDark = useToggle(isDark)
const router = useRouter()
const forceTags = ref(false)
const forceArt = ref(false)

// convert streamQualities to an array of string | number
const streamQualitiesArray = computed<(string | number)[]>(() => {
  return Object.values(streamQualities)
})

async function runScan() {
  const formData = new FormData()
  formData.append('force', forceTags.value ? 'true' : 'false')
  formData.append('include-art', forceArt.value ? 'true' : 'false')
  await openSubsonicFetchRequest('startScan.view', {
    body: formData,
  })
}

function setStreamQuality(quality: StreamQuality) {
  if (streamQuality.value === quality) {
    return
  }
  streamQuality.value = quality
}

async function logOut() {
  clearApiKey()
  router.push('/login')
}
</script>

<template>
  <div class="p-4 flex flex-col gap-y-6">
    <div class="flex gap-4">
      <ZButton @click="runScan()">
        <span class="text-nowrap">Run a scan</span>
      </ZButton>
      <label for="force" class="flex gap-x-2 items-center">
        <input id="force" v-model="forceTags" type="checkbox" class="size-6">
        <span class="text-nowrap">Force</span>
      </label>
      <label v-if="forceTags" for="include-art" class="flex gap-x-2 items-center">
        <input id="include-art" v-model="forceArt" type="checkbox" class="size-6">
        <span class="text-nowrap">Include Art</span>
      </label>
    </div>
    <ZButton :primary="debugEnabled" @click="toggleDebug()">
      <span class="text-nowrap">Debug: {{ debugEnabled ? 'On' : 'Off' }}</span>
    </ZButton>
    <ZButton :primary="isDark" @click="toggleDark()">
      <span class="text-nowrap">Dark Mode: {{ isDark ? 'On' : 'Off' }}</span>
    </ZButton>
    <ZButton @click="logOut()">
      <span class="text-nowrap">Logout</span>
    </ZButton>
    <DropdownMenu
      title="Stream Quality"
      :options="streamQualitiesArray"
      align="right"
      @select="setStreamQuality"
    />
    <UserManagement />
  </div>
</template>
