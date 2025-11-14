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
    class="navlink group/navlink"
    :class="{
      'translate-x--8': currentRoute !== routeProp,
      'hover:translate-x-0': currentRoute !== routeProp && !currentRoute.startsWith(`${routeProp}/`),
      'translate-x-0': currentRoute === routeProp || currentRoute.startsWith(`${routeProp}/`),
    }"
    @click="handleLinkClick"
  >
    <icon-nrk-media-ffw
      class="size-8 text-primary1 opacity-0 transition-all duration-300"
      :class="{
        'opacity-100': currentRoute === routeProp || currentRoute.startsWith(`${routeProp}/`),
        'group-hover/navlink:opacity-50': currentRoute !== routeProp && !currentRoute.startsWith(`${routeProp}/`),
      }"
    />
    {{ routeName }}
  </RouterLink>
</template>

<style scoped>
.navlink {
  @apply block flex  text-muted font-semibold no-underline transition-all duration-100 uppercase flex items-center;
}
</style>
