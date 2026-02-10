<script setup lang="ts">
import type { SubsonicUser } from '~/types/subsonicUser'
import { fetchCurrentUser } from '~/logic/users'

const currentUser = ref<SubsonicUser | null>(null)

async function getCurrentUser() {
  currentUser.value = await fetchCurrentUser()
  console.log('Current user:', currentUser.value)
}

onBeforeMount(async () => {
  await getCurrentUser()
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
