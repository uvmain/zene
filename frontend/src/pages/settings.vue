<script setup>
import { useDark, useToggle } from '@vueuse/core'
import { openSubsonicFetchRequest } from '~/logic/backendFetch'
import { clearApiKey } from '~/logic/common'
import { debugLog, toggleDebug } from '~/logic/logger'
import { debugEnabled, streamQualities, streamQuality } from '~/logic/store'

const isDark = useDark()
const toggleDark = useToggle(isDark)
const router = useRouter()

async function runScan() {
  const response = await openSubsonicFetchRequest('startScan.view')
  const json = await response.json()
  debugLog(JSON.stringify(json))
}

async function logOut() {
  clearApiKey()
  router.push('/login')
}
</script>

<template>
  <div class="p-4 flex flex-col gap-y-6">
    <ZButton @click="runScan()">
      <span class="text-nowrap">Run a scan</span>
    </ZButton>
    <ZButton :primary="debugEnabled" @click="toggleDebug()">
      <span class="text-nowrap">Debug: {{ debugEnabled ? 'On' : 'Off' }}</span>
    </ZButton>
    <ZButton :primary="isDark" @click="toggleDark()">
      <span class="text-nowrap">Dark Mode: {{ isDark ? 'On' : 'Off' }}</span>
    </ZButton>
    <ZButton @click="logOut()">
      <span class="text-nowrap">Logout</span>
    </ZButton>
    <div class="px-4 py-3 lg:py-2">
      <label class="text-base text-gray-500 mb-2 block lg:text-sm lg:mb-1">Stream Quality</label>
      <select
        v-model="streamQuality"
        class="-md text-base text-gray-700 px-3 py-2 border border-gray-300 w-full lg:text-sm lg:px-2 lg:py-1 focus:outline-none focus:ring focus:ring-blue-200"
      >
        <option
          v-for="quality in streamQualities"
          :key="quality"
          :value="quality"
        >
          {{ quality === 'native' ? 'Original Quality' : `${quality} kbps` }}
        </option>
      </select>
    </div>
    <UserManagement />
  </div>
</template>
