<script setup lang="ts">
import type { SubsonicUser } from '~/types/subsonicUser'
import { getAuthenticatedAvatarUrl, postAvatarImage } from '~/logic/backendFetch'
import { onImageError } from '~/logic/common'
import { createUser, defaultNewUser, deleteUser, fetchCurrentUser, fetchUsers, updateUser } from '~/logic/users'

const users = ref<SubsonicUser[]>([])
const currentUser = ref<SubsonicUser>({} as SubsonicUser)
const showCreateUserDialog = ref(false)
const showEditUserDialog = ref(false)
const showDeleteUserDialog = ref(false)
const newUser = ref<SubsonicUser>({ ...defaultNewUser })
const editingUser = ref<SubsonicUser>({} as SubsonicUser)
const userToDelete = ref<SubsonicUser>({} as SubsonicUser)
const avatar = ref<string | null>(null)

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
  await deleteUser(userToDelete.value)
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
  editingUser.value = user
  showEditUserDialog.value = true
}

function openDeleteUserDialog(user: SubsonicUser) {
  if (!currentUser.value?.adminRole)
    return
  userToDelete.value = user
  showDeleteUserDialog.value = true
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
                  <img
                    :src="getAuthenticatedAvatarUrl(user.id)"
                    alt="User Avatar"
                    class="rounded-full size-10 object-cover"
                    @error="onImageError"
                  />
                  <span>{{ user.username }}</span>
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

      <Modal :show-modal="showCreateUserDialog" modal-title="Create New User" @close="closeCreateUserDialog">
        <template #content>
          <form class="text-muted flex flex-col gap-4">
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
          </form>
        </template>
      </Modal>

      <!-- Edit User Dialog -->
      <div v-if="showEditUserDialog && editingUser" class="bg-gray-600 bg-opacity-50 flex h-full w-full items-center inset-0 justify-center fixed z-30 overflow-y-auto">
        <div class="p-6 bg-white max-w-md w-full shadow-xl">
          <h3 class="text-lg text-gray-900 leading-6 font-medium mb-4">
            Edit User: {{ editingUser.username }}
          </h3>
          <div>
            <div class="mb-4">
              <label class="flex items-center">
                <input v-model="editingUser.adminRole" type="checkbox" class="text-indigo-600 border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50 focus:ring-offset-0">
                <span class="text-sm text-gray-600 ml-2">Is Admin</span>
              </label>
              <div class="mb-4">
                <label for="new-password" class="text-sm text-gray-700 font-medium">Password</label>
                <input id="new-password" v-model="editingUser.password" type="password" class="mt-1 px-3 py-2 border border-gray-300 w-full shadow-sm sm:text-sm focus:outline-none focus:border-indigo-500 focus:ring-indigo-500">
              </div>
              <div class="mb-4">
                <label for="new-email" class="text-sm text-gray-700 font-medium">Email</label>
                <input id="new-email" v-model="editingUser.email" type="email" class="mt-1 px-3 py-2 border border-gray-300 w-full shadow-sm sm:text-sm focus:outline-none focus:border-indigo-500 focus:ring-indigo-500">
              </div>
              <div v-if="avatar" class="text-black">
                avatar: {{ avatar }}
              </div>
              <ImageSelector v-model="avatar" />
            </div>
            <div class="mt-6 flex justify-end space-x-3">
              <button type="button" class="text-sm text-gray-700 font-medium px-4 py-2 bg-gray-100 hover:bg-gray-200" @click="showEditUserDialog = false">
                Cancel
              </button>
              <button type="button" class="text-sm font-medium px-4 py-2 bg-yellow-600 hover:bg-yellow-700" @click="handleUpdateUser()">
                Update
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Delete User Confirmation Dialog -->
      <div v-if="showDeleteUserDialog && userToDelete" class="bg-gray-600 bg-opacity-50 flex h-full w-full items-center inset-0 justify-center fixed z-30 overflow-y-auto">
        <div class="p-6 bg-white max-w-md w-full shadow-xl">
          <h3 class="text-lg text-gray-900 leading-6 font-medium mb-2">
            Confirm Deletion
          </h3>
          <p class="text-sm text-gray-500 mb-4">
            Are you sure you want to delete the user "{{ userToDelete.username }}"? This action cannot be undone.
          </p>
          <div class="mt-6 flex justify-end space-x-3">
            <button type="button" class="text-sm text-gray-700 font-medium px-4 py-2 bg-gray-100 hover:bg-gray-200" @click="showDeleteUserDialog = false">
              Cancel
            </button>
            <button class="text-sm font-medium px-4 py-2 bg-red-600 hover:bg-red-700" @click="handleDeleteUser">
              Delete
            </button>
          </div>
        </div>
      </div>
      <div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="css">
/* Basic styling for the table and modals, can be expanded or use a UI library like TailwindCSS more extensively */
button {
  transition: background-color 0.2s ease-in-out;
}
/* Basic modal transition */
.modal-enter-active, .modal-leave-active {
  transition: opacity 0.3s ease;
}
.modal-enter-from, .modal-leave-to {
  opacity: 0;
}
</style>
