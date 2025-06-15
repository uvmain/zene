<script setup lang="ts">
import type { User } from '../types/auth'
import { backendFetchRequest, getCurrentUser, getUsers } from '../composables/fetchFromBackend'

const users = ref<User[]>([])
const currentUser = ref<User | null>(null)
const showCreateUserDialog = ref(false)
const showEditUserDialog = ref(false)
const showDeleteUserDialog = ref(false)

const newUser = ref({ username: '', password: '', isAdmin: false })
const editingUser = ref<User | null>(null)
const userToDelete = ref<User | null>(null)

async function fetchCurrentUser() {
  currentUser.value = await getCurrentUser()
}

async function fetchUsers() {
  users.value = await getUsers()
}

async function handleCreateUser() {
  if (!currentUser.value?.is_admin)
    return
  try {
    const response = await backendFetchRequest('users', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newUser.value),
    })
    if (!response.ok) {
      const errData = await response.json()
      throw new Error(errData.message || `Failed to create user: ${response.statusText}`)
    }
    await fetchUsers()
    showCreateUserDialog.value = false
    newUser.value = { username: '', password: '', isAdmin: false }
  }
  catch (error) {
    console.error(`Error creating user: ${error}`)
  }
}

async function handleUpdateUser() {
  if (!currentUser.value?.is_admin || !editingUser.value)
    return
  try {
    const response = await backendFetchRequest(`users/${editingUser.value.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: editingUser.value.username,
        // password: sort this out later
        is_admin: editingUser.value.is_admin,
      }),
    })
    if (!response.ok) {
      const errData = await response.json()
      throw new Error(errData.message || `Failed to update user: ${response.statusText}`)
    }
    await fetchUsers()
    showEditUserDialog.value = false
    editingUser.value = null
  }
  catch (error) {
    console.error(`Error updating user: ${error}`)
  }
}

async function handleDeleteUser() {
  if (!currentUser.value?.is_admin || !userToDelete.value)
    return
  try {
    const response = await backendFetchRequest(`users/${userToDelete.value.id}`, {
      method: 'DELETE',
    })
    if (!response.ok) {
      const errData = await response.json()
      throw new Error(errData.message || `Failed to delete user: ${response.statusText}`)
    }
    await fetchUsers()
    showDeleteUserDialog.value = false
    userToDelete.value = null
  }
  catch (error) {
    console.error(`Error deleting user: ${error}`)
  }
}

function openCreateUserDialog() {
  if (!currentUser.value?.is_admin)
    return
  newUser.value = { username: '', password: '', isAdmin: false }
  showCreateUserDialog.value = true
}

function openEditUserDialog(user: User) {
  if (!currentUser.value?.is_admin)
    return
  editingUser.value = { ...user }
  showEditUserDialog.value = true
}

function openDeleteUserDialog(user: User) {
  if (!currentUser.value?.is_admin)
    return
  userToDelete.value = user
  showDeleteUserDialog.value = true
}

onMounted(async () => {
  await fetchCurrentUser()
  if (currentUser.value?.is_admin) {
    await fetchUsers()
  }
})
</script>

<template>
  <div class="p-4 md:p-6">
    <h1 class="mb-6 text-2xl font-semibold">
      Manage Users
    </h1>
    <div v-if="currentUser?.is_admin">
      <div>
        <div class="mb-4">
          <button
            class="rounded bg-blue-500 px-4 py-2 text-white font-bold hover:bg-blue-700"
            @click="openCreateUserDialog"
          >
            Create User
          </button>
        </div>

        <div v-if="!users.length" class="py-4 text-center">
          Loading users...
        </div>
        <div v-if="users.length > 0" class="overflow-x-auto">
          <table class="min-w-full rounded-lg bg-white shadow-md">
            <thead class="bg-gray-200">
              <tr>
                <th class="px-4 py-3 text-left text-xs text-gray-500 font-medium tracking-wider uppercase">
                  Username
                </th>
                <th class="px-4 py-3 text-left text-xs text-gray-500 font-medium tracking-wider uppercase">
                  Admin
                </th>
                <th class="px-4 py-3 text-left text-xs text-gray-500 font-medium tracking-wider uppercase">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-for="user in users" :key="user.id">
                <td class="whitespace-nowrap px-4 py-4 text-gray-600">
                  {{ user.username }}
                </td>
                <td class="whitespace-nowrap px-4 py-4 text-gray-600">
                  {{ user.is_admin ? 'Yes' : 'No' }}
                </td>
                <td class="whitespace-nowrap px-4 py-4 space-x-2">
                  <button
                    class="rounded bg-yellow-500 px-3 py-1 text-sm text-white font-semibold hover:bg-yellow-600"
                    @click="openEditUserDialog(user)"
                  >
                    Edit
                  </button>
                  <button
                    class="rounded bg-red-500 px-3 py-1 text-sm text-white font-semibold hover:bg-red-600"
                    @click="openDeleteUserDialog(user)"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div v-if="showCreateUserDialog" class="fixed inset-0 z-30 h-full w-full flex items-center justify-center overflow-y-auto bg-gray-600 bg-opacity-50">
        <div class="max-w-md w-full rounded-lg bg-white p-6 shadow-xl">
          <h3 class="mb-4 text-lg text-gray-900 font-medium leading-6">
            Create New User
          </h3>
          <form @submit.prevent="handleCreateUser">
            <div class="mb-4">
              <label for="new-username" class="block text-sm text-gray-700 font-medium">Username</label>
              <input id="new-username" v-model="newUser.username" type="text" required class="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 shadow-sm focus:border-indigo-500 sm:text-sm focus:outline-none focus:ring-indigo-500">
            </div>
            <div class="mb-4">
              <label for="new-password" class="block text-sm text-gray-700 font-medium">Password</label>
              <input id="new-password" v-model="newUser.password" type="password" required class="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 shadow-sm focus:border-indigo-500 sm:text-sm focus:outline-none focus:ring-indigo-500">
            </div>
            <div class="mb-4">
              <label class="flex items-center">
                <input v-model="newUser.isAdmin" type="checkbox" class="border-gray-300 rounded text-indigo-600 shadow-sm focus:border-indigo-300 focus:ring focus:ring-offset-0 focus:ring-indigo-200 focus:ring-opacity-50">
                <span class="ml-2 text-sm text-gray-600">Is Admin</span>
              </label>
            </div>
            <div class="mt-6 flex justify-end space-x-3">
              <button type="button" class="rounded-md bg-gray-100 px-4 py-2 text-sm text-gray-700 font-medium hover:bg-gray-200" @click="showCreateUserDialog = false">
                Cancel
              </button>
              <button type="submit" class="rounded-md bg-blue-600 px-4 py-2 text-sm text-white font-medium hover:bg-blue-700">
                Create
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Edit User Dialog -->
      <div v-if="showEditUserDialog && editingUser" class="fixed inset-0 z-30 h-full w-full flex items-center justify-center overflow-y-auto bg-gray-600 bg-opacity-50">
        <div class="max-w-md w-full rounded-lg bg-white p-6 shadow-xl">
          <h3 class="mb-4 text-lg text-gray-900 font-medium leading-6">
            Edit User: {{ editingUser.username }}
          </h3>
          <form @submit.prevent="handleUpdateUser">
            <div class="mb-4">
              <label for="edit-username" class="block text-sm text-gray-700 font-medium">Username</label>
              <input id="edit-username" v-model="editingUser.username" type="text" required class="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 shadow-sm focus:border-indigo-500 sm:text-sm focus:outline-none focus:ring-indigo-500">
            </div>
            <div class="mb-4">
              <label class="flex items-center">
                <input v-model="editingUser.is_admin" type="checkbox" class="border-gray-300 rounded text-indigo-600 shadow-sm focus:border-indigo-300 focus:ring focus:ring-offset-0 focus:ring-indigo-200 focus:ring-opacity-50">
                <span class="ml-2 text-sm text-gray-600">Is Admin</span>
              </label>
            </div>
            <div class="mt-6 flex justify-end space-x-3">
              <button type="button" class="rounded-md bg-gray-100 px-4 py-2 text-sm text-gray-700 font-medium hover:bg-gray-200" @click="showEditUserDialog = false">
                Cancel
              </button>
              <button type="submit" class="rounded-md bg-yellow-600 px-4 py-2 text-sm text-white font-medium hover:bg-yellow-700">
                Update
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Delete User Confirmation Dialog -->
      <div v-if="showDeleteUserDialog && userToDelete" class="fixed inset-0 z-30 h-full w-full flex items-center justify-center overflow-y-auto bg-gray-600 bg-opacity-50">
        <div class="max-w-md w-full rounded-lg bg-white p-6 shadow-xl">
          <h3 class="mb-2 text-lg text-gray-900 font-medium leading-6">
            Confirm Deletion
          </h3>
          <p class="mb-4 text-sm text-gray-500">
            Are you sure you want to delete the user "{{ userToDelete.username }}"? This action cannot be undone.
          </p>
          <div class="mt-6 flex justify-end space-x-3">
            <button type="button" class="rounded-md bg-gray-100 px-4 py-2 text-sm text-gray-700 font-medium hover:bg-gray-200" @click="showDeleteUserDialog = false">
              Cancel
            </button>
            <button class="rounded-md bg-red-600 px-4 py-2 text-sm text-white font-medium hover:bg-red-700" @click="handleDeleteUser">
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

<style scoped>
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
