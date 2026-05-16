<script setup lang="ts">
import type { StreamQuality } from '~/logic/store'
import { useDark, useToggle } from '@vueuse/core'
import { deleteAudioCache, openSubsonicFetchRequest } from '~/logic/backendFetch'
import { resetAccentColour, updateAccentColour } from '~/logic/colours'
import { clearApiKey } from '~/logic/common'
import { toggleDebug } from '~/logic/logger'
import { accentColour, debugEnabled, streamQualities, streamQuality } from '~/logic/store'

const isDark = useDark()
const toggleDark = useToggle(isDark)
const router = useRouter()
const forceTags = ref(false)
const forceArt = ref(false)
const showLogoutModal = ref(false)

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
        <input id="force" v-model="forceTags" type="checkbox" class="accent-primary-400 size-4">
        <span class="text-nowrap">Force</span>
      </label>
      <label v-if="forceTags" for="include-art" class="flex gap-x-2 items-center">
        <input id="include-art" v-model="forceArt" type="checkbox" class="accent-primary-400 size-4">
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
      :current-option="streamQuality"
      @select="setStreamQuality"
    />

    <div class="flex flex-row gap-2 items-center">
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
