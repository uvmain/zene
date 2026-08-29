<script setup lang="ts">
import { fetchShares } from '~/logic/backendFetch'
import type { SubsonicSharesResponse } from '~/types/share'
import type { SubsonicUser } from '~/types/subsonicUser'
import { fetchCurrentUser } from '~/logic/users'

const currentUser = ref<SubsonicUser>({} as SubsonicUser)
const shares = ref<SubsonicSharesResponse | null>(null)

onMounted(async () => {
  const response = await fetchCurrentUser()
  if (response) {
    currentUser.value = response
  }
  shares.value = await fetchShares()
})
</script>

<template>
  <div v-if="currentUser?.adminRole" class="flex flex-col gap-4 lg:max-w-7xl">
    <h1 class="text-2xl font-semibold">
      Manage Shares
    </h1>
    <div v-if="shares">
      <div class="text-primary px-4 py-3 corner-cut background-3 share-grid uppercase">
        <span>User</span>
        <span class="col-span-2">Share</span>
        <span>Visits</span>
        <span>Created</span>
      </div>
      <div v-for="share in shares.shares.share" :key="share.id" class="text-muted p-4 border-primary border-t-0 background-2 share-grid">
        <div> {{ share.username }} </div>
        <a class="col-span-2 underline" :href="share.url" target="_blank">{{ share.description }}</a>
        <div> {{ share.visitCount }} </div>
        <div> {{ new Date(share.created).toLocaleString() }} </div>
      </div>
    </div>
    <div v-else>
      Loading shares...
    </div>
  </div>
</template>

<style lang="css" scoped>
.share-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 1rem;
}
</style>