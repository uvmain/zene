<script setup lang="ts">
import { useNavbar } from '~/composables/useNavbar'
import { useSearch } from '~/composables/useSearch'

defineProps({
  routeName: { type: String, required: true },
  routeProp: { type: String, required: true },
})

const route = useRoute()
const { closeSearch } = useSearch()
const { closeMobileNav } = useNavbar()

const currentRoute = computed(() => {
  return route.path
})

function handleLinkClick() {
  closeSearch()
  closeMobileNav()
}
</script>

<template>
  <RouterLink
    :to="routeProp"
    class="navlink"
    :class="{
      'translate-x--8': currentRoute !== routeProp,
      'hover:translate-x--6': currentRoute !== routeProp,
      'translate-x-0': currentRoute === routeProp,
    }"
    @click="handleLinkClick"
  >
    <icon-nrk-arrow-right
      class="size-8 translate-y--1 text-cyan opacity-0 transition-all duration-300"
      :class="{
        'opacity-100': currentRoute === routeProp,
      }"
    />
    {{ routeName }}
  </RouterLink>
</template>

<style scoped>
a {
  @apply uppercase;
}

.navlink {
  @apply block flex gap-x-1 py-2 text-white font-semibold no-underline transition-all duration-100;
}
</style>
