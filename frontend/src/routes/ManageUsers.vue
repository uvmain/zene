<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { backendFetchRequest } from '@/composables/fetchFromBackend' // Adjusted path

// Define types for User and CurrentUser
interface User {
  id: number
  username: string
  is_admin: boolean
}

interface CurrentUser {
  id: number
  username: string
  is_admin: boolean
}

const users = ref<User[]>([])
const currentUser = ref<CurrentUser | null>(null)
const isLoading = ref(true)
const error = ref<string | null>(null)

// Dialog states
const showCreateUserDialog = ref(false)
const showEditUserDialog = ref(false)
const showDeleteUserDialog = ref(false)

// Form data
const newUser = ref({ username: '', password: '', isAdmin: false })
const editingUser = ref<User | null>(null)
const userToDelete = ref<User | null>(null)

// Fetch current user data (assuming an endpoint like /api/me)
async function fetchCurrentUser() {
  try {
    // TODO: Replace with actual endpoint if different
    const response = await backendFetchRequest('/api/me', { method: 'GET' })
    if (!response.ok) {
      throw new Error(`Failed to fetch current user: ${response.statusText}`)
    }
    currentUser.value = await response.json()
  } catch (e: any) {
    error.value = e.message
    console.error("Error fetching current user:", e)
  }
}

// Fetch all users
async function fetchUsers() {
  if (!currentUser.value?.is_admin) {
    // Non-admins should not be able to fetch all users
    // This check should ideally be on the backend, but good for UX too
    users.value = []
    isLoading.value = false
    // error.value = "You do not have permission to view users." // Optional: display error
    return
  }
  isLoading.value = true
  error.value = null
  try {
    // TODO: Replace with actual endpoint /api/users
    const response = await backendFetchRequest('/api/users', { method: 'GET' })
    if (!response.ok) {
      throw new Error(`Failed to fetch users: ${response.statusText}`)
    }
    users.value = await response.json()
  } catch (e: any) {
    error.value = e.message
    console.error("Error fetching users:", e)
  } finally {
    isLoading.value = false
  }
}

// CRUD operations
async function handleCreateUser() {
  if (!currentUser.value?.is_admin) return
  try {
    // TODO: Replace with actual endpoint POST /api/users
    const response = await backendFetchRequest('/api/users', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newUser.value),
    })
    if (!response.ok) {
      const errData = await response.json()
      throw new Error(errData.message || `Failed to create user: ${response.statusText}`)
    }
    await fetchUsers() // Refresh users list
    showCreateUserDialog.value = false
    newUser.value = { username: '', password: '', isAdmin: false } // Reset form
  } catch (e: any) {
    error.value = e.message
    console.error("Error creating user:", e)
  }
}

async function handleUpdateUser() {
  if (!currentUser.value?.is_admin || !editingUser.value) return
  try {
    // TODO: Replace with actual endpoint PUT /api/users/:id
    const response = await backendFetchRequest(`/api/users/${editingUser.value.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: editingUser.value.username,
        // Only send password if it's being changed. Assume backend handles empty password as no change.
        // password: editingUser.value.password, // This needs a dedicated field in the edit form
        is_admin: editingUser.value.is_admin,
      }),
    })
    if (!response.ok) {
      const errData = await response.json()
      throw new Error(errData.message || `Failed to update user: ${response.statusText}`)
    }
    await fetchUsers() // Refresh users list
    showEditUserDialog.value = false
    editingUser.value = null
  } catch (e: any) {
    error.value = e.message
    console.error("Error updating user:", e)
  }
}

async function handleDeleteUser() {
  if (!currentUser.value?.is_admin || !userToDelete.value) return
  try {
    // TODO: Replace with actual endpoint DELETE /api/users/:id
    const response = await backendFetchRequest(`/api/users/${userToDelete.value.id}`, {
      method: 'DELETE',
    })
    if (!response.ok) {
      const errData = await response.json()
      throw new Error(errData.message || `Failed to delete user: ${response.statusText}`)
    }
    await fetchUsers() // Refresh users list
    showDeleteUserDialog.value = false
    userToDelete.value = null
  } catch (e: any) {
    error.value = e.message
    console.error("Error deleting user:", e)
  }
}

// Open dialog functions
function openCreateUserDialog() {
  if (!currentUser.value?.is_admin) return
  newUser.value = { username: '', password: '', isAdmin: false }
  error.value = null; // Clear previous errors
  showCreateUserDialog.value = true
}

function openEditUserDialog(user: User) {
  if (!currentUser.value?.is_admin) return
  editingUser.value = { ...user } // Clone user to avoid modifying original data directly
  error.value = null; // Clear previous errors
  // Add a field for new password if needed, separate from editingUser.password_hash
  showEditUserDialog.value = true
}

function openDeleteUserDialog(user: User) {
  if (!currentUser.value?.is_admin) return
  userToDelete.value = user
  error.value = null; // Clear previous errors
  showDeleteUserDialog.value = true
}

onMounted(async () => {
  await fetchCurrentUser()
  if (currentUser.value?.is_admin) {
    await fetchUsers()
  } else {
    isLoading.value = false
    // Optionally, set an error message or redirect if not admin
    // error.value = "You do not have permission to access this page.";
  }
})
</script>

<template>
  <div class="p-4 md:p-6">
    <h1 class="text-2xl font-semibold mb-6">Manage Users</h1>

    <div v-if="!currentUser?.is_admin && !isLoading" class="text-red-500">
      You do not have permission to manage users.
    </div>

    <div v-if="currentUser?.is_admin">
      <div class="mb-4">
        <button
          @click="openCreateUserDialog"
          class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        >
          Create User
        </button>
      </div>

      <div v-if="isLoading && !users.length" class="text-center py-4">Loading users...</div>
      <div v-if="error && !isLoading" class="text-red-500 bg-red-100 p-3 rounded mb-4">
        Error: {{ error }}
      </div>

      <div v-if="!isLoading && users.length > 0" class="overflow-x-auto">
        <table class="min-w-full bg-white shadow-md rounded-lg">
          <thead class="bg-gray-200">
            <tr>
              <th class="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Username</th>
              <th class="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Admin</th>
              <th class="py-3 px-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="user in users" :key="user.id">
              <td class="py-4 px-4 whitespace-nowrap">{{ user.username }}</td>
              <td class="py-4 px-4 whitespace-nowrap">{{ user.is_admin ? 'Yes' : 'No' }}</td>
              <td class="py-4 px-4 whitespace-nowrap space-x-2">
                <button
                  @click="openEditUserDialog(user)"
                  class="text-sm bg-yellow-500 hover:bg-yellow-600 text-white font-semibold py-1 px-3 rounded"
                >
                  Edit
                </button>
                <button
                  @click="openDeleteUserDialog(user)"
                  class="text-sm bg-red-500 hover:bg-red-600 text-white font-semibold py-1 px-3 rounded"
                >
                  Delete
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
       <div v-if="!isLoading && users.length === 0 && !error" class="text-center py-4 text-gray-500">
        No users found.
      </div>
    </div>

    <!-- Create User Dialog -->
    <div v-if="showCreateUserDialog" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full flex justify-center items-center z-30">
      <div class="bg-white p-6 rounded-lg shadow-xl w-full max-w-md">
        <h3 class="text-lg font-medium leading-6 text-gray-900 mb-4">Create New User</h3>
        <form @submit.prevent="handleCreateUser">
          <div v-if="error" class="text-red-500 bg-red-100 p-3 rounded mb-4">Error: {{ error }}</div>
          <div class="mb-4">
            <label for="new-username" class="block text-sm font-medium text-gray-700">Username</label>
            <input type="text" id="new-username" v-model="newUser.username" required class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
          </div>
          <div class="mb-4">
            <label for="new-password" class="block text-sm font-medium text-gray-700">Password</label>
            <input type="password" id="new-password" v-model="newUser.password" required class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
          </div>
          <div class="mb-4">
            <label class="flex items-center">
              <input type="checkbox" v-model="newUser.isAdmin" class="rounded border-gray-300 text-indigo-600 shadow-sm focus:border-indigo-300 focus:ring focus:ring-offset-0 focus:ring-indigo-200 focus:ring-opacity-50">
              <span class="ml-2 text-sm text-gray-600">Is Admin</span>
            </label>
          </div>
          <div class="mt-6 flex justify-end space-x-3">
            <button type="button" @click="showCreateUserDialog = false" class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-md">Cancel</button>
            <button type="submit" class="px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-md">Create</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Edit User Dialog -->
    <div v-if="showEditUserDialog && editingUser" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full flex justify-center items-center z-30">
      <div class="bg-white p-6 rounded-lg shadow-xl w-full max-w-md">
        <h3 class="text-lg font-medium leading-6 text-gray-900 mb-4">Edit User: {{ editingUser.username }}</h3>
        <form @submit.prevent="handleUpdateUser">
          <div v-if="error" class="text-red-500 bg-red-100 p-3 rounded mb-4">Error: {{ error }}</div>
          <div class="mb-4">
            <label for="edit-username" class="block text-sm font-medium text-gray-700">Username</label>
            <input type="text" id="edit-username" v-model="editingUser.username" required class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
          </div>
          <!-- Add password change field if required -->
          <!-- <div class="mb-4">
            <label for="edit-password" class="block text-sm font-medium text-gray-700">New Password (optional)</label>
            <input type="password" id="edit-password" v-model="newPasswordForEdit" class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
          </div> -->
          <div class="mb-4">
            <label class="flex items-center">
              <input type="checkbox" v-model="editingUser.is_admin" class="rounded border-gray-300 text-indigo-600 shadow-sm focus:border-indigo-300 focus:ring focus:ring-offset-0 focus:ring-indigo-200 focus:ring-opacity-50">
              <span class="ml-2 text-sm text-gray-600">Is Admin</span>
            </label>
          </div>
          <div class="mt-6 flex justify-end space-x-3">
            <button type="button" @click="showEditUserDialog = false" class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-md">Cancel</button>
            <button type="submit" class="px-4 py-2 text-sm font-medium text-white bg-yellow-600 hover:bg-yellow-700 rounded-md">Update</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Delete User Confirmation Dialog -->
    <div v-if="showDeleteUserDialog && userToDelete" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full flex justify-center items-center z-30">
      <div class="bg-white p-6 rounded-lg shadow-xl w-full max-w-md">
        <h3 class="text-lg font-medium leading-6 text-gray-900 mb-2">Confirm Deletion</h3>
        <div v-if="error" class="text-red-500 bg-red-100 p-3 rounded mb-4">Error: {{ error }}</div>
        <p class="text-sm text-gray-500 mb-4">
          Are you sure you want to delete the user "{{ userToDelete.username }}"? This action cannot be undone.
        </p>
        <div class="mt-6 flex justify-end space-x-3">
          <button type="button" @click="showDeleteUserDialog = false" class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-md">Cancel</button>
          <button @click="handleDeleteUser" class="px-4 py-2 text-sm font-medium text-white bg-red-600 hover:bg-red-700 rounded-md">Delete</button>
        </div>
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
