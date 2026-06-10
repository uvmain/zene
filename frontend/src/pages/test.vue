<script setup lang="ts">
import { isBrowserChrome } from '~/logic/browser'

const isChrome = isBrowserChrome()
</script>

<template>
  is Chrome?: {{ isChrome }}
  <ChromecastButton />
</template>
