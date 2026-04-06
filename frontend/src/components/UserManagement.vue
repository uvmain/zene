<script setup lang="ts">
import type { SubsonicUser } from '~/types/subsonicUser'
import { createUser, deleteUser, fetchCurrentUser, fetchUsers, updateUser } from '~/logic/users'

const users = ref<SubsonicUser[]>([])
const currentUser = ref<SubsonicUser>({} as SubsonicUser)
const showCreateUserDialog = ref(false)
const showEditUserDialog = ref(false)
const showDeleteUserDialog = ref(false)

const newUser = ref<SubsonicUser>({} as SubsonicUser)
const editingUser = ref<SubsonicUser>({} as SubsonicUser)
const userToDelete = ref<SubsonicUser>({} as SubsonicUser)

async function getCurrentUser() {
  currentUser.value = await fetchCurrentUser()
}

async function getUsers() {
  users.value = await fetchUsers()
}

async function handleCreateUser() {
  if (!currentUser.value?.adminRole || !newUser.value)
    return

  await createUser(newUser.value)

  await fetchUsers()
  showCreateUserDialog.value = false
  newUser.value = {} as SubsonicUser
}

async function handleUpdateUser() {
  if (!currentUser.value?.adminRole || !editingUser.value)
    return
  await updateUser(editingUser.value)
  await fetchUsers()
  showEditUserDialog.value = false
  editingUser.value = {} as SubsonicUser
}

async function handleDeleteUser() {
  if (!currentUser.value?.adminRole || !userToDelete.value)
    return
  await deleteUser(userToDelete.value)
  await fetchUsers()
  showDeleteUserDialog.value = false
  userToDelete.value = {} as SubsonicUser
}

function openCreateUserDialog() {
  if (!currentUser.value?.adminRole)
    return
  newUser.value = {} as SubsonicUser
  showCreateUserDialog.value = true
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
  <div v-if="!currentUser?.adminRole" class="text-red-500">
    You do not have permission to manage users.
  </div>
  <div v-else class="p-4 lg:p-6">
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
            Create User
          </button>
        </div>

        <div v-if="!users.length" class="py-4 text-center">
          Loading users...
        </div>
        <div v-if="users.length > 0" class="overflow-x-auto">
          <table class="bg-white min-w-full shadow-md">
            <thead class="bg-gray-200">
              <tr>
                <th class="text-xs text-gray-500 tracking-wider font-medium px-4 py-3 text-left uppercase">
                  Username
                </th>
                <th class="text-xs text-gray-500 tracking-wider font-medium px-4 py-3 text-left uppercase">
                  Admin
                </th>
                <th class="text-xs text-gray-500 tracking-wider font-medium px-4 py-3 text-left uppercase">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody class="divide-gray-200 divide-y">
              <tr v-for="user in users" :key="user.username">
                <td class="text-gray-600 px-4 py-4 whitespace-nowrap">
                  {{ user.username }}
                </td>
                <td class="text-gray-600 px-4 py-4 whitespace-nowrap">
                  {{ user.adminRole ? 'Yes' : 'No' }}
                </td>
                <td class="px-4 py-4 whitespace-nowrap space-x-2">
                  <button
                    class="text-sm font-semibold px-3 py-1 bg-yellow-500 hover:bg-yellow-600"
                    @click="openEditUserDialog(user)"
                  >
                    Edit
                  </button>
                  <button
                    class="text-sm font-semibold px-3 py-1 bg-red-500 hover:bg-red-600"
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

      <div v-if="showCreateUserDialog" class="bg-gray-600 bg-opacity-50 flex h-full w-full items-center inset-0 justify-center fixed z-30 overflow-y-auto">
        <div class="p-6 bg-white max-w-md w-full shadow-xl">
          <h3 class="text-lg text-gray-900 leading-6 font-medium mb-4">
            Create New User
          </h3>
          <form @submit.prevent="handleCreateUser">
            <div class="mb-4">
              <label for="new-username" class="text-sm text-gray-700 font-medium">Username</label>
              <input id="new-username" v-model="newUser.username" type="text" required class="mt-1 px-3 py-2 border border-gray-300 w-full shadow-sm sm:text-sm focus:outline-none focus:border-indigo-500 focus:ring-indigo-500">
            </div>
            <div class="mb-4">
              <label for="new-password" class="text-sm text-gray-700 font-medium">Password</label>
              <input id="new-password" v-model="newUser.password" type="password" required class="mt-1 px-3 py-2 border border-gray-300 w-full shadow-sm sm:text-sm focus:outline-none focus:border-indigo-500 focus:ring-indigo-500">
            </div>
            <div class="mb-4">
              <label for="new-email" class="text-sm text-gray-700 font-medium">Email</label>
              <input id="new-email" v-model="newUser.email" type="email" required class="mt-1 px-3 py-2 border border-gray-300 w-full shadow-sm sm:text-sm focus:outline-none focus:border-indigo-500 focus:ring-indigo-500">
            </div>
            <div class="mb-4">
              <label class="flex items-center">
                <input v-model="newUser.adminRole" type="checkbox" class="text-indigo-600 border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50 focus:ring-offset-0">
                <span class="text-sm text-gray-600 ml-2">Is Admin</span>
              </label>
            </div>
            <div class="mt-6 flex justify-end space-x-3">
              <button type="button" class="text-sm text-gray-700 font-medium px-4 py-2 bg-gray-100 hover:bg-gray-200" @click="showCreateUserDialog = false">
                Cancel
              </button>
              <button type="submit" class="text-sm font-medium px-4 py-2 bg-blue-600 hover:bg-blue-700">
                Create
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Edit User Dialog -->
      <div v-if="showEditUserDialog && editingUser" class="bg-gray-600 bg-opacity-50 flex h-full w-full items-center inset-0 justify-center fixed z-30 overflow-y-auto">
        <div class="p-6 bg-white max-w-md w-full shadow-xl">
          <h3 class="text-lg text-gray-900 leading-6 font-medium mb-4">
            Edit User: {{ editingUser.username }}
          </h3>
          <form @submit.prevent="handleUpdateUser">
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
            </div>
            <div class="mt-6 flex justify-end space-x-3">
              <button type="button" class="text-sm text-gray-700 font-medium px-4 py-2 bg-gray-100 hover:bg-gray-200" @click="showEditUserDialog = false">
                Cancel
              </button>
              <button type="submit" class="text-sm font-medium px-4 py-2 bg-yellow-600 hover:bg-yellow-700">
                Update
              </button>
            </div>
          </form>
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
