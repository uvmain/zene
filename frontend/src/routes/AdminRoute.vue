<script setup lang="ts">
import type { SubsonicUser } from '~/types/subsonicUser'
import { useBackendFetch } from '~/composables/useBackendFetch'

const { getCurrentUser } = useBackendFetch()

const currentUser = ref<SubsonicUser | null>(null)

async function fetchCurrentUser() {
  currentUser.value = await getCurrentUser()
  console.log('Current user:', currentUser.value)
}

onBeforeMount(async () => {
  await fetchCurrentUser()
})
</script>

<template>
  <div>
    <div v-if="!currentUser?.adminRole" class="text-red-500">
      You do not have permission to manage users.
    </div>
    <UserManagement v-else />
  </div>
</template>
