<script setup lang="ts">
import { openMobileNav } from '~/logic/navbar'
</script>

<template>
  <header class="p-2 flex items-center justify-center lg:hidden">
    <icon-nrk-list class="lg:text-2xl" @click="openMobileNav()" />
  </header>
</template>
