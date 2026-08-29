<script setup lang="ts">
import { fetchShares } from '~/logic/backendFetch'
import type { Share, SubsonicSharesResponse } from '~/types/share'
import type { SubsonicUser } from '~/types/subsonicUser'
import { fetchCurrentUser } from '~/logic/users'
import { SubsonicResponse } from '~/types/subsonic'
import { openSubsonicFetchRequest } from '~/logic/backendFetch'

const currentUser = ref<SubsonicUser>({} as SubsonicUser)
const shares = ref<SubsonicSharesResponse | null>(null)
const shareToDelete = ref<Share>({} as Share)
const showDeleteShareDialog = ref(false)

function openDeleteShareDialog(share: Share) {
  if (!currentUser.value?.adminRole)
    return
  shareToDelete.value = share
  showDeleteShareDialog.value = true
}

function closeDeleteShareDialog() {
  showDeleteShareDialog.value = false
  shareToDelete.value = {} as Share
}

async function handleDeleteShare() {
  if (!shareToDelete.value.id) {
    throw new Error('No share selected for deletion')
  }
  if (!currentUser.value?.adminRole && shareToDelete.value.username !== currentUser.value.username) {
    return
  }
  const formData = new FormData()
  const shareId = shareToDelete.value.id.toString()
  formData.append('id', shareId)
  const response = await openSubsonicFetchRequest<SubsonicResponse>('deleteshare', {
    body: formData,
  })
  if (response.status !== 'ok') {
    throw new Error(response.error?.message ?? 'Unknown error')
  }
  closeDeleteShareDialog()
  shares.value = await fetchShares()
}

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
        <ZButton
          :title="`Delete ${share.description}`"
          :red="true"
          @click="openDeleteShareDialog(share)"
        >
        Delete
      </ZButton>
      </div>
    </div>
    <div v-else>
      Loading shares...
    </div>

    <!-- Delete Share Modal -->
    <Modal :show-modal="showDeleteShareDialog" modal-title="Delete Share" @close="closeDeleteShareDialog">
      <template #content>
        <p class="text-muted">
          Are you sure you want to delete the share "{{ shareToDelete.description }}"? This action cannot be undone.
        </p>
        <div class="mt-6 flex justify-end space-x-3">
          <ZButton :green="true" @click="closeDeleteShareDialog">
            Cancel
          </ZButton>
          <ZButton :red="true" @click="handleDeleteShare">
            Delete
          </ZButton>
        </div>
      </template>
    </Modal>
  </div>
</template>

<style lang="css" scoped>
.share-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 1rem;
}
</style>