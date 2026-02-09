<script setup lang="ts">
import { useNavbar } from '~/composables/useNavbar'

defineProps({
  routeName: { type: String, required: true },
  routeProp: { type: String, required: true },
})

const route = useRoute()
const { closeMobileNav } = useNavbar()

const currentRoute = computed(() => {
  return route.path
})

function handleLinkClick() {
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
      class="size-8 opacity-0 transition-all duration-300"
      :class="{
        'opacity-100 text-primary1': currentRoute === routeProp || currentRoute.startsWith(`${routeProp}/`),
        'group-hover/navlink:opacity-50 text-secondary1': currentRoute !== routeProp && !currentRoute.startsWith(`${routeProp}/`),
      }"
    />
    <span class="text-2xl lg:text-lg">{{ routeName }}</span>
  </RouterLink>
</template>

<style scoped>
.navlink {
  @apply block flex  text-muted font-semibold no-underline transition-all duration-100 uppercase flex items-center;
}
</style>
