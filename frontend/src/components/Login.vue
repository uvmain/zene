<script setup lang="ts">
import { checkIfLoggedIn, login } from '../composables/auth'

const username = ref('')
const password = ref('')

onBeforeMount(() => {
  checkIfLoggedIn()
})
</script>

<template>
  <div class="fixed left-0 top-0 z-999 size-full animate-fade-in backdrop-blur-xl">
    <div class="mx-auto mb-auto mt-150px w-300px border rounded px-30px pb-30px pt-20px">
      <div class="w-300 flex flex-col gap-4 p-6">
        <form class="flex flex-col gap-2">
          <div class="flex flex-row items-center gap-2">
            <label for="username">Username:</label>
            <input id="username" v-model="username" type="text" name="username" autocomplete="username" class="recipeCardBackground text border rounded px-2 py-1">
          </div>
          <div class="flex flex-row items-center gap-2">
            <label for="password">Password:</label>
            <input id="password" v-model="password" type="password" name="password" autocomplete="current-password" class="recipeCardBackground text border rounded px-2 py-1" @keydown.enter="login(username, password)">
          </div>
        </form>
      </div>
      <div class="flex justify-center gap-4">
        <button aria-label="login" class="textButton px-4 py-2" @click="login(username, password)">
          Login
        </button>
      </div>
    </div>
  </div>
</template>
