<script setup>
import { useBackendFetch } from '../composables/useBackendFetch'
import { useDebug } from '../composables/useDebug'
import { useSettings } from '../composables/useSettings'

const { streamQuality, StreamQualities } = useSettings()
const { backendFetchRequest } = useBackendFetch()
const { toggleDebug, debugLog, useDebugBool } = useDebug()

const open = ref(false)
const dropdownRef = ref(null)

function toggle() {
  open.value = !open.value
}

function close() {
  open.value = false
}

function handleClickOutside(event) {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target)) {
    close()
  }
}

async function runScan() {
  const response = await backendFetchRequest('scan', {
    method: 'POST',
  })
  const json = await response.json()
  debugLog(JSON.stringify(json))
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div ref="dropdownRef" class="relative z-20 inline-block text-left">
    <icon-tabler-settings class="text-2xl" @click="toggle" />

    <transition name="fade">
      <div
        v-if="open"
        class="absolute right-0 mt-2 w-56 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5"
      >
        <div class="py-1">
          <div
            class="block cursor-pointer px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
            @click="runScan"
          >
            Run a Scan
          </div>
          <div
            class="block cursor-pointer px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
            @click="toggleDebug()"
          >
            Debug: {{ useDebugBool ? 'On' : 'Off' }}
          </div>
          <div class="px-4 py-2">
            <label class="mb-1 block text-sm text-gray-500">Stream Quality</label>
            <select
              v-model="streamQuality"
              class="w-full border border-gray-300 rounded-md px-2 py-1 text-sm text-gray-700 focus:outline-none focus:ring focus:ring-blue-200"
            >
              <option
                v-for="quality in StreamQualities"
                :key="quality"
                :value="quality"
              >
                {{ quality === 'native' ? 'Original Quality' : `${quality} kbps` }}
              </option>
            </select>
          </div>
          <router-link
            to="/admin"
            class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
            @click="close"
          >
            Manage Users
          </router-link>
        </div>
      </div>
    </transition>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
