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
  users.value = await fetchUsers()
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
  <div v-if="currentUser?.adminRole" class="mt-4">
    <h1 class="text-2xl font-semibold mb-6">
      Manage Users
    </h1>
    <div>
      <div>
        <div class="mb-4">
          <button
            class="z-button"
            @click="openCreateUserDialog"
          >
            Create new user
          </button>
        </div>
        <div>
          <table class="text-left background-2 w-full">
            <thead class="text-primary background-3">
              <tr>
                <th class="text-xs px-4 py-3 uppercase">
                  User
                </th>
                <th class="text-xs px-4 py-3 uppercase">
                  Admin
                </th>
                <th class="text-xs px-4 py-3 uppercase">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody class="text-muted divide-background-300 divide-y dark:divide-background-700">
              <tr v-for="user in users" :key="user.username">
                <td class="px-4 py-2 flex flex-row whitespace-nowrap items-center space-x-3">
                  <Avatar :user="user" />
                </td>
                <td class="px-4 whitespace-nowrap">
                  {{ user.adminRole ? 'Yes' : 'No' }}
                </td>
                <td class="px-4 py-4 flex flex-row whitespace-nowrap space-x-2">
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
                </td>
              </tr>
            </tbody>
          </table>
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
              <input id="is-admin" v-model="newUser.adminRole" type="checkbox" class="text-primary border-gray-300 focus:border-indigo-300 focus:ring focus:ring-primary-200 focus:ring-opacity-50 focus:ring-offset-0">
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
              <input id="is-admin" v-model="editingUser.adminRole" type="checkbox" class="text-primary border-gray-300 focus:border-indigo-300 focus:ring focus:ring-primary-200 focus:ring-opacity-50 focus:ring-offset-0">
              <span class="text-sm">Admin user</span>
            </div>
            <p v-if="adminCount === 1 && editingUserAdminRole" class="text-red-600 max-w-lg text-wrap">
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
          <p v-if="adminCount === 1 && userToDelete.adminRole" class="text-red-600">
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
