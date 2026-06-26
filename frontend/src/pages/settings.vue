<script setup lang="ts">
import type { StreamQuality } from '~/stores/main'
import type { FfVersionsResponse } from '~/types'
import { useDark, useToggle } from '@vueuse/core'
import { deleteAudioCache, downloadNewFfBinaries, fetchFfVersions, openSubsonicFetchRequest } from '~/logic/backendFetch'
import { initializeAccentColour, resetAccentColour } from '~/logic/colours'
import { clearApiKey, toggleWakeLock } from '~/logic/common'
import { toggleDebug } from '~/logic/logger'
import * as Store from '~/stores/main'

const isDark = useDark()
const toggleDark = useToggle(isDark)
const router = useRouter()
const forceTags = ref(false)
const forceArt = ref(false)
const showLogoutModal = ref(false)
const ffVersions = ref<FfVersionsResponse | null>(null)

const streamQualitiesArray = computed<(string | number)[]>(() => {
  return Object.values(Store.streamQualities)
})

const currentStreamQuality = computed(() => {
  return Store.streamQuality.value
})

function colourToHex(colour: string): string {
  const parser = new Option().style
  parser.color = ''
  parser.color = colour

  const rgbMatch = parser.color.match(/^rgb\((\d+),\s*(\d+),\s*(\d+)\)$/)
  if (!rgbMatch) {
    return '#fa742f'
  }

  const [r, g, b] = rgbMatch.slice(1).map(value => Number.parseInt(value, 10))
  const toHex = (value: number) => value.toString(16).padStart(2, '0')
  return `#${toHex(r)}${toHex(g)}${toHex(b)}`
}

const accentColourInputValue = computed({
  get: () => colourToHex(Store.accentColour.value),
  set: (value: string) => {
    Store.accentColour.value = value
    initializeAccentColour()
  },
})

watch(Store.autoSwitchColours, (newValue) => {
  if (!newValue) {
    initializeAccentColour()
  }
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
  if (Store.streamQuality.value === quality) {
    return
  }
  Store.streamQuality.value = quality
}

async function getFfVersions() {
  ffVersions.value = await fetchFfVersions()
}

async function logOut() {
  clearApiKey()
  router.push('/login')
}

async function downloadFfBinaries() {
  await downloadNewFfBinaries()
  await getFfVersions()
}

function resetBackendOverrides() {
  Store.backendUrl.value = new URL(window.location.href).origin
}

onMounted(() => {
  getFfVersions()
})
</script>

<template>
  <div class="p-4 flex flex-col gap-y-6">
    <div class="flex gap-4">
      <ZButton @click="runScan()">
        <span class="text-nowrap">Run a scan</span>
      </ZButton>
      <label for="force" class="flex gap-x-2 items-center">
        <input id="force" v-model="forceTags" type="checkbox" class="accent-main-400 size-4">
        <span class="text-nowrap">Force</span>
      </label>
      <label v-if="forceTags" for="include-art" class="flex gap-x-2 items-center">
        <input id="include-art" v-model="forceArt" type="checkbox" class="accent-main-400 size-4">
        <span class="text-nowrap">Include Art</span>
      </label>
    </div>
    <ZButton @click="toggleDebug()">
      <span class="text-nowrap">Debug: {{ Store.debugEnabled.value ? 'On' : 'Off' }}</span>
    </ZButton>
    <ZButton @click="toggleDark()">
      <span class="text-nowrap">Dark Mode: {{ isDark ? 'On' : 'Off' }}</span>
    </ZButton>
    <ZButton @click="showLogoutModal = true">
      <span class="text-nowrap">Logout</span>
    </ZButton>
    <ZButton @click="deleteAudioCache()">
      <span class="text-nowrap">Delete Audio Cache</span>
    </ZButton>
    <DropdownMenu
      title="Stream Quality"
      :options="streamQualitiesArray"
      align="right"
      :current-option="currentStreamQuality"
      @select="setStreamQuality"
    />

    <div class="flex flex-row gap-2 items-center">
      <label for="auto-switch-colours" class="flex gap-x-2 items-center">
        <input id="auto-switch-colours" v-model="Store.autoSwitchColours.value" type="checkbox" class="accent-main-400 size-4">
        <span class="text-nowrap">Auto Switch Colours</span>
      </label>
      <input
        id="accent"
        v-model="accentColourInputValue"
        type="color"
        name="accent"
        colorspace="display-p3"
      />
      <label for="accent">Accent color</label>
      <ZButton @click="resetAccentColour()">
        Reset
      </ZButton>
    </div>

    <div class="flex flex-row gap-2 items-center">
      <label for="wake-lock" class="flex gap-x-2 items-center">
        <input id="wake-lock" v-model="Store.wakeLockEnabled.value" type="checkbox" class="accent-main-400 size-4" @change="toggleWakeLock">
        <span class="text-nowrap">Wake Lock - this is required to prevent the screen from turning off while using chromecast</span>
      </label>
    </div>

    <div v-if="ffVersions" class="mr-auto p-2 border-muted corner-cut flex flex-col gap-y-2">
      <div>
        <span class="text-lg font-semibold">FFmpeg Version: </span>
        <span>{{ ffVersions.ffmpeg_version }}</span>
      </div>
      <div>
        <span class="text-lg font-semibold">FFprobe Version: </span>
        <span>{{ ffVersions.ffprobe_version }}</span>
      </div>
      <ZButton @click="downloadFfBinaries()">
        Update FFmpeg and FFprobe
      </ZButton>
    </div>

    <div class="mr-auto flex flex-row gap-4 items-center">
      <label for="override-backend-url" class="text-lg font-semibold flex gap-x-2 items-center">
        <span>Override Backend URL: </span>
        <input
          id="override-backend-url"
          v-model="Store.backendUrl.value"
          type="text"
          class="p-2 w-64"
          placeholder="using default backend URL"
        />
      </label>
      <ZButton @click="resetBackendOverrides">
        Reset Backend URL
      </ZButton>
    </div>

    <UserManagement />

    <!-- Logout Modal -->
    <Modal :show-modal="showLogoutModal" modal-title="Logout" @close="showLogoutModal = false">
      <template #content>
        <div class="p-6">
          <h1 class="text-xl text-primary font-semibold mb-6">
            Confirm Logout
          </h1>
          <div class="flex justify-center space-x-3">
            <ZButton @click="showLogoutModal = false">
              Cancel
            </ZButton>
            <ZButton @click="logOut">
              Logout
            </ZButton>
          </div>
          <div>
          </div>
        </div>
      </template>
    </Modal>
  </div>
</template>
