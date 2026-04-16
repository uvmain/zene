<script setup lang="ts">
import { isMobileNavOpen, toggleMobileNav } from '~/logic/navbar'
</script>

<template>
  <icon-nrk-media-ffw
    class="size-12 rotate-y-0 transition-all duration-300 fixed z-100"
    :class="{
      'rotate-y-180 text-accent-400 translate-x-2': isMobileNavOpen,
      'text-primary': !isMobileNavOpen,
    }"
    @click="toggleMobileNav()"
  />
</template>
