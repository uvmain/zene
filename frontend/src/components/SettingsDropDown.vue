<script setup>
import { openSubsonicFetchRequest } from '~/logic/backendFetch'
import { debugLog, toggleDebug } from '~/logic/logger'
import { debugEnabled, streamQualities, streamQuality } from '~/logic/store'

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
  <abbr title="Settings">
    <div ref="dropdown" class="text-left flex inline-block items-center relative z-20">
      <icon-nrk-settings class="text-2xl cursor-pointer" @click="toggle" />

      <transition name="fade">
        <div
          v-if="open"
          class="mt-2 bg-white w-64 ring-1 ring-black ring-opacity-5 shadow-lg origin-top-right right-0 absolute lg:w-56"
        >
          <div class="py-1">
            <div
              class="text-base text-gray-700 px-4 py-3 block cursor-pointer lg:text-sm lg:py-2 hover:bg-gray-100"
              @click="runScan"
            >
              Run a Scan
            </div>
            <div
              class="text-base text-gray-700 px-4 py-3 block cursor-pointer lg:text-sm lg:py-2 hover:bg-gray-100"
              @click="toggleDebug()"
            >
              Debug: {{ debugEnabled ? 'On' : 'Off' }}
            </div>
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
            <router-link
              to="/admin"
              class="text-base text-gray-700 px-4 py-3 block lg:text-sm lg:py-2 hover:bg-gray-100"
              @click="close"
            >
              Manage Users
            </router-link>
          </div>
        </div>
      </transition>
    </div>
  </abbr>
</template>

<style scoped lang="css">
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
