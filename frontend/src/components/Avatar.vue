<script setup lang="ts">
import type { SubsonicUser } from '~/types/subsonicUser'
import { getAuthenticatedAvatarUrl } from '~/logic/backendFetch'

const props = defineProps({
  user: { type: Object as PropType<SubsonicUser>, required: true },
})

const avatarUrl = ref<string>('')
const useDefaultAvatar = ref<boolean>(false)

function onImageError() {
  useDefaultAvatar.value = true
}

onMounted(() => {
  avatarUrl.value = getAuthenticatedAvatarUrl(props.user.id)
})
</script>

<template>
  <img
    v-if="!useDefaultAvatar && avatarUrl && avatarUrl.length > 0"
    :src="avatarUrl"
    alt="User Avatar"
    class="rounded-full size-10 object-cover"
    @error="onImageError"
  />
  <div v-else class="rounded-full background-3 flex size-10 items-center justify-center">
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 64 64"
      fill="none"
      width="64"
      height="64"
      stroke-width="4"
      stroke-linecap="round"
      stroke-linejoin="round"
      aria-hidden="true"
      role="img"
      class="p-1 stroke-text-900 dark:stroke-text-100"
    >
      <circle id="head" cx="32" cy="22" r="10" />
      <path id="shoulders" d="M14 52c0-10 8-18 18-18s18 8 18 18" />
    </svg>
  </div>
  <span>{{ user.username }}</span>
</template>
