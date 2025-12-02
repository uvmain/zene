<script setup>
import { openSubsonicFetchRequest } from '~/composables/backendFetch'
import { useDebug } from '~/composables/useDebug'
import { useSettings } from '~/composables/useSettings'

const { streamQuality, StreamQualities } = useSettings()
const { toggleDebug, debugLog, useDebugBool } = useDebug()

const open = ref(false)
const dropdownRef = useTemplateRef('dropdown')

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
  const response = await openSubsonicFetchRequest('startScan.view')
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
  <div ref="dropdown" class="relative z-20 inline-block text-left">
    <icon-nrk-settings class="cursor-pointer text-2xl" @click="toggle" />

    <transition name="fade">
      <div
        v-if="open"
        class="-md absolute right-0 mt-2 w-64 origin-top-right bg-white shadow-lg ring-1 ring-black ring-opacity-5 md:w-56"
      >
        <div class="py-1">
          <div
            class="block cursor-pointer px-4 py-3 text-base text-gray-700 hover:bg-gray-100 md:py-2 md:text-sm"
            @click="runScan"
          >
            Run a Scan
          </div>
          <div
            class="block cursor-pointer px-4 py-3 text-base text-gray-700 hover:bg-gray-100 md:py-2 md:text-sm"
            @click="toggleDebug()"
          >
            Debug: {{ useDebugBool ? 'On' : 'Off' }}
          </div>
          <div class="px-4 py-3 md:py-2">
            <label class="mb-2 block text-base text-gray-500 md:mb-1 md:text-sm">Stream Quality</label>
            <select
              v-model="streamQuality"
              class="-md w-full border border-gray-300 px-3 py-2 text-base text-gray-700 md:px-2 md:py-1 md:text-sm focus:outline-none focus:ring focus:ring-blue-200"
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
            class="block px-4 py-3 text-base text-gray-700 hover:bg-gray-100 md:py-2 md:text-sm"
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
