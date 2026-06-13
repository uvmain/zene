<script setup lang="ts">
import type { StreamQuality } from '~/stores/main'
import { useDark, useToggle } from '@vueuse/core'
import { deleteAudioCache, downloadNewFfBinaries, fetchFfVersions, openSubsonicFetchRequest } from '~/logic/backendFetch'
import { initializeAccentColour, resetAccentColour, updateAccentColour } from '~/logic/colours'
import { clearApiKey } from '~/logic/common'
import { toggleDebug } from '~/logic/logger'
import { accentColour, autoSwitchColours, debugEnabled, streamQualities, streamQuality } from '~/stores/main'

interface FfVersions {
  ffmpeg_version: string
  ffprobe_version: string
}

const isDark = useDark()
const toggleDark = useToggle(isDark)
const router = useRouter()
const forceTags = ref(false)
const forceArt = ref(false)
const showLogoutModal = ref(false)
const ffVersions = ref<FfVersions | null>(null)

const streamQualitiesArray = computed<(string | number)[]>(() => {
  return Object.values(streamQualities)
})

const currentStreamQuality = computed(() => {
  return streamQuality.value
})

watch(autoSwitchColours, (newValue) => {
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
  if (streamQuality.value === quality) {
    return
  }
  streamQuality.value = quality
}

async function getFfVersions() {
  const versions = await fetchFfVersions()
  ffVersions.value = versions
}

async function logOut() {
  clearApiKey()
  router.push('/login')
}

async function downloadFfBinaries() {
  await downloadNewFfBinaries()
  await getFfVersions()
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
      <span class="text-nowrap">Debug: {{ debugEnabled ? 'On' : 'Off' }}</span>
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
        <input id="auto-switch-colours" v-model="autoSwitchColours" type="checkbox" class="accent-main-400 size-4">
        <span class="text-nowrap">Auto Switch Colours</span>
      </label>
      <input
        id="accent"
        v-model="accentColour"
        type="color"
        name="accent"
        colorspace="display-p3"
        @input="updateAccentColour"
      />
      <label for="accent">Accent color</label>
      <ZButton @click="resetAccentColour()">
        Reset
      </ZButton>
    </div>

    <div v-if="ffVersions" class="mr-auto p-2 border-muted corner-cut flex flex-col gap-y-2">
      <div class="text-lg font-semibold">
        FFmpeg Version: {{ ffVersions.ffmpeg_version }}
      </div>
      <div class="text-lg font-semibold">
        FFprobe Version: {{ ffVersions.ffprobe_version }}
      </div>
      <ZButton @click="downloadFfBinaries()">
        Update FFmpeg and FFprobe
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
