<script setup lang="ts">
import { useDark, useSessionStorage, useToggle } from '@vueuse/core'
import { getRandomSlug } from '../composables/slugs'

defineProps({
  showAdd: { type: Boolean, default: false },
  showEdit: { type: Boolean, default: false },
})

const emits = defineEmits(['add', 'edit'])

const isModalOpened = ref(false)
const isDark = useDark()
const toggleDark = useToggle(isDark)

function openModal() {
  isModalOpened.value = true
}

function closeModal() {
  isModalOpened.value = false
}

const router = useRouter()

const userLoginState = useSessionStorage('login-state', false)

function enableEdit() {
  if (userLoginState.value) {
    emits('edit')
  }
  else {
    console.warn('Unable to enter edit mode, please log in')
  }
}

async function navigateHome() {
  router.push('/')
}

function navigateAlbums() {
  router.push('/albums')
}

async function navigateRandom() {
  const slug = await getRandomSlug()
  router.push(`/${slug}`)
}

async function navigateUpload() {
  if (userLoginState.value)
    router.push('/upload')
}

async function navigateTags() {
  router.push('/tags')
}
</script>

<template>
  <div class="px-2 pt-2 standard lg:px-6">
    <header class="mx-0 mx-auto flex justify-around pt-2 lg:max-w-8/10 lg:justify-between lg:p-6">
      <div class="flex items-center gap-0 gap-4">
        <div class="hidden p-2 text-2xl lg:block hover:cursor-pointer" @click="navigateHome">
          home
        </div>
        <icon-tabler-home class="p-2 text-2xl lg:hidden" @click="navigateHome" />

        <div class="hidden p-2 text-2xl lg:block hover:cursor-pointer" @click="navigateAlbums">
          albums
        </div>
        <icon-tabler-vinyl class="p-2 text-2xl lg:hidden" @click="navigateAlbums" />

        <div class="hidden p-2 text-2xl lg:block hover:cursor-pointer" @click="navigateTags">
          tags
        </div>
        <icon-tabler-tags class="p-2 text-2xl lg:hidden" @click="navigateTags" />

        <div class="hidden p-2 text-2xl lg:block hover:cursor-pointer" @click="navigateRandom">
          random
        </div>
        <icon-tabler-arrows-shuffle class="p-2 text-2xl lg:hidden" @click="navigateRandom" />
      </div>
      <div class="flex gap-4">
        <TooltipIcon v-if="userLoginState" tooltip-text="Upload" class="hover:cursor-pointer" @click="navigateUpload">
          <icon-tabler-upload class="text-2xl" />
        </TooltipIcon>
        <slot name="1" />
        <TooltipIcon v-if="showAdd && userLoginState" tooltip-text="Add" class="hover:cursor-pointer" @click="emits('add')">
          <icon-tabler-library-plus class="text-2xl" />
        </TooltipIcon>
        <slot name="2" />
        <TooltipIcon v-if="showEdit && userLoginState" tooltip-text="Edit Mode" class="hover:cursor-pointer" @click="enableEdit">
          <icon-tabler-edit class="text-2xl" />
        </TooltipIcon>
        <slot name="3" />
        <TooltipIcon :tooltip-text="isDark ? 'Light Mode' : 'Dark Mode'" class="hover:cursor-pointer" @click="toggleDark()">
          <icon-tabler-sun v-if="isDark" class="text-2xl" />
          <icon-tabler-moon-stars v-else class="text-2xl" />
        </TooltipIcon>
        <slot name="4" />
        <TooltipIcon :tooltip-text="userLoginState ? 'Log Out' : 'Log In'" class="hover:cursor-pointer" @click="openModal">
          <icon-tabler-user class="text-2xl" />
        </TooltipIcon>
      </div>
      <LoginModal :is-open="isModalOpened" @modal-close="closeModal" @submit="navigateHome" />
    </header>
  </div>
</template>
