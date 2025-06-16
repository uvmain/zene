<script setup lang="ts">
import type { User } from '../types/auth'
import { useBackendFetch } from '../composables/useBackendFetch'

const { getCurrentUser } = useBackendFetch()

const currentUser = ref<User | null>(null)

async function fetchCurrentUser() {
  currentUser.value = await getCurrentUser()
}

onBeforeMount(async () => {
  await fetchCurrentUser()
})
</script>

<template>
  <div>
    <div v-if="!currentUser?.is_admin" class="text-red-500">
      You do not have permission to manage users.
    </div>
    <UserManagement v-else />
  </div>
</template>
