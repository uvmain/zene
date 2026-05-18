<script setup lang="ts">
import type { SubsonicUser } from '~/types/subsonicUser'
import { postAvatarImage } from '~/logic/backendFetch'
import { debugLog } from '~/logic/logger'
import { createUser, defaultNewUser, deleteUser, fetchCurrentUser, fetchUsers, updateUser } from '~/logic/users'

const users = ref<SubsonicUser[]>([])
const currentUser = ref<SubsonicUser>({} as SubsonicUser)
const showCreateUserDialog = ref(false)
const showEditUserDialog = ref(false)
const showDeleteUserDialog = ref(false)
const newUser = ref<SubsonicUser>({ ...defaultNewUser })
const editingUser = ref<SubsonicUser>({} as SubsonicUser)
const editingUserAdminRole = ref(false)
const userToDelete = ref<SubsonicUser>({} as SubsonicUser)
const avatar = ref<string | null>(null)

const adminCount = computed(() => users.value.filter(user => user.adminRole).length)

async function getCurrentUser() {
  currentUser.value = await fetchCurrentUser()
}

async function getUsers() {
  // sort by adminRole first
  users.value = (await fetchUsers()).sort((a, b) => Number(b.adminRole) - Number(a.adminRole))
}

async function handleCreateUser() {
  if (!currentUser.value?.adminRole || !newUser.value) {
    return
  }
  await createUser(newUser.value)
  if (avatar.value) {
    const imageBlob = await (await fetch(avatar.value)).blob()
    await postAvatarImage({ userId: editingUser.value.id, file: imageBlob })
    avatar.value = null
  }
  await getUsers()
  showCreateUserDialog.value = false
  newUser.value = { ...defaultNewUser }
}

async function handleUpdateUser() {
  if (!currentUser.value?.adminRole || !editingUser.value)
    return
  await updateUser(editingUser.value)
  if (avatar.value) {
    const imageBlob = await (await fetch(avatar.value)).blob()
    await postAvatarImage({ userId: editingUser.value.id, file: imageBlob })
    avatar.value = null
  }
  showEditUserDialog.value = false
  editingUser.value = {} as SubsonicUser
  await getUsers()
}

async function handleDeleteUser() {
  if (!currentUser.value?.adminRole || !userToDelete.value)
    return
  const userDeletionResponse = await deleteUser(userToDelete.value)
  if (userDeletionResponse.status !== 'ok') {
    debugLog(`Failed to delete user: ${userDeletionResponse.error?.message}`)
  }
  await getUsers()
  showDeleteUserDialog.value = false
  userToDelete.value = {} as SubsonicUser
}

function openCreateUserDialog() {
  if (!currentUser.value?.adminRole)
    return
  newUser.value = { ...defaultNewUser }
  showCreateUserDialog.value = true
}

function closeCreateUserDialog() {
  showCreateUserDialog.value = false
  newUser.value = { ...defaultNewUser }
}

function openEditUserDialog(user: SubsonicUser) {
  if (!currentUser.value?.adminRole)
    return
  editingUser.value = { ...user }
  editingUserAdminRole.value = user.adminRole
  showEditUserDialog.value = true
}

function closeEditUserDialog() {
  showEditUserDialog.value = false
  editingUser.value = {} as SubsonicUser
  editingUserAdminRole.value = false
}

function openDeleteUserDialog(user: SubsonicUser) {
  if (!currentUser.value?.adminRole)
    return
  userToDelete.value = user
  showDeleteUserDialog.value = true
}

function closeDeleteUserDialog() {
  showDeleteUserDialog.value = false
  userToDelete.value = {} as SubsonicUser
}

onMounted(async () => {
  await getCurrentUser()
  if (currentUser.value?.adminRole) {
    await getUsers()
  }
})
</script>

<template>
  <div v-if="currentUser?.adminRole" class="flex flex-col gap-4">
    <h1 class="text-2xl font-semibold">
      Manage Users
    </h1>
    <ZButton
      hover-text="Create new user"
      @click="openCreateUserDialog"
    >
      Create new user
    </ZButton>
    <div>
      <div class="lg:w-2xl">
        <div class="text-primary px-4 py-3 corner-cut background-3 gap-4 grid grid-cols-3 uppercase">
          <span>User</span>
          <span>Admin</span>
          <span>Actions</span>
        </div>
        <div v-for="user in users" :key="user.username" class="text-muted p-4 border-primary border-t-0 background-2 gap-4 grid grid-cols-3 items-center">
          <Avatar :user="user" />
          <span>{{ user.adminRole ? 'Yes' : 'No' }}</span>
          <div class="flex flex-col gap-3 lg:flex-row">
            <ZButton
              :hover-text="`Edit ${user.username} user`"
              @click="openEditUserDialog(user)"
            >
              Edit
            </ZButton>
            <ZButton
              :hover-text="`Delete ${user.username} user`"
              :red="true"
              @click="openDeleteUserDialog(user)"
            >
              Delete
            </ZButton>
          </div>
        </div>
      </div>

      <!-- Create User Modal -->
      <Modal :show-modal="showCreateUserDialog" modal-title="Create New User" @close="closeCreateUserDialog">
        <template #content>
          <div class="text-muted flex flex-col gap-4">
            <div class="flex flex-row gap-2 items-center justify-between">
              <label for="new-username" class="text-sm">Username</label>
              <input id="new-username" v-model="newUser.username" type="text" required class="input-text">
            </div>
            <div class="flex flex-row gap-2 items-center justify-between">
              <label for="new-password" class="text-sm">Password</label>
              <input id="new-password" v-model="newUser.password" type="password" required class="input-text">
            </div>
            <div class="flex flex-row gap-2 items-center justify-between">
              <label for="new-email" class="text-sm">Email</label>
              <input id="new-email" v-model="newUser.email" type="email" required class="input-text">
            </div>
            <div class="flex flex-row gap-2 items-center justify-center">
              <label for="is-admin" class="flex items-center" />
              <input id="is-admin" v-model="newUser.adminRole" type="checkbox" class="text-primary border-gray-300 focus:border-indigo-300 focus:ring focus:ring-main-200 focus:ring-opacity-50 focus:ring-offset-0">
              <span class="text-sm">Admin user</span>
            </div>
          </div>
          <ImageSelector v-model="avatar">
            <template #title>
              <span class="text-sm text-muted">
                Optional user avatar
              </span>
            </template>
          </ImageSelector>
          <div class="mt-6 flex justify-center space-x-3">
            <ZButton :red="true" @click="closeCreateUserDialog">
              Cancel
            </ZButton>
            <ZButton :green="true" @click="handleCreateUser">
              Create
            </ZButton>
          </div>
        </template>
      </Modal>

      <!-- Edit User Modal -->
      <Modal :show-modal="showEditUserDialog" modal-title="Edit User" @close="closeEditUserDialog">
        <template #content>
          <div class="text-muted flex flex-col gap-4">
            <div class="flex flex-row gap-2 items-center justify-between">
              <label for="new-username" class="text-sm">Username</label>
              <input id="new-username" v-model="editingUser.username" type="text" required class="input-text">
            </div>
            <div class="flex flex-row gap-2 items-center justify-between">
              <label for="new-password" class="text-sm">Password</label>
              <input id="new-password" v-model="editingUser.password" type="password" required class="input-text">
            </div>
            <div class="flex flex-row gap-2 items-center justify-between">
              <label for="new-email" class="text-sm">Email</label>
              <input id="new-email" v-model="editingUser.email" type="email" required class="input-text">
            </div>
            <div class="flex flex-row gap-2 items-center justify-center">
              <label for="is-admin" class="flex items-center" />
              <input id="is-admin" v-model="editingUser.adminRole" type="checkbox" class="text-primary border-gray-300 focus:border-indigo-300 focus:ring focus:ring-main-200 focus:ring-opacity-50 focus:ring-offset-0">
              <span class="text-sm">Admin user</span>
            </div>
            <p v-if="adminCount === 1 && editingUserAdminRole" class="text-red-400 max-w-lg text-wrap">
              Warning: This is the only admin user. Removing admin role will leave the system without any admin accounts.
            </p>
          </div>
          <ImageSelector v-model="avatar">
            <template #title>
              <span class="text-sm text-muted">
                Optional user avatar
              </span>
            </template>
          </ImageSelector>
          <div class="mt-6 flex justify-center space-x-3">
            <ZButton :red="true" @click="closeEditUserDialog">
              Cancel
            </ZButton>
            <ZButton :green="true" @click="handleUpdateUser">
              Update
            </ZButton>
          </div>
        </template>
      </Modal>

      <!-- Delete User Modal -->
      <Modal :show-modal="showDeleteUserDialog" modal-title="Delete User" @close="closeDeleteUserDialog">
        <template #content>
          <p class="text-muted">
            Are you sure you want to delete the user "{{ userToDelete.username }}"? This action cannot be undone.
          </p>
          <p v-if="adminCount === 1 && userToDelete.adminRole" class="text-red-400 max-w-lg text-wrap">
            Warning: This is the only admin user. Deleting this user will leave the system without any admin accounts.
          </p>
          <div class="mt-6 flex justify-end space-x-3">
            <ZButton :green="true" @click="closeDeleteUserDialog">
              Cancel
            </ZButton>
            <ZButton :red="true" @click="handleDeleteUser">
              Delete
            </ZButton>
          </div>
        </template>
      </Modal>
      <div>
      </div>
    </div>
  </div>
</template>
