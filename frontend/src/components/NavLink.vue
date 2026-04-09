<script setup lang="ts">
import { closeMobileNav } from '~/logic/navbar'

defineProps({
  routeName: { type: String, required: true },
  routeProp: { type: String, required: true },
})

const route = useRoute()

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
    class="group/navlink text-muted font-semibold no-underline flex uppercase transition-all duration-100 items-center"
    :class="{
      'translate-x--8': currentRoute !== routeProp,
      'hover:translate-x-0': currentRoute !== routeProp && !currentRoute.startsWith(`${routeProp}/`),
      'translate-x-0': currentRoute === routeProp || currentRoute.startsWith(`${routeProp}/`),
    }"
    @click="handleLinkClick"
  >
    <icon-nrk-media-ffw
      class="opacity-0 size-8 transition-all duration-300"
      :class="{
        'opacity-100 text-primary-500': currentRoute === routeProp || currentRoute.startsWith(`${routeProp}/`),
        'group-hover/navlink:opacity-50 text-secondary-500': currentRoute !== routeProp && !currentRoute.startsWith(`${routeProp}/`),
      }"
    />
    <span class="text-2xl lg:text-lg">{{ routeName }}</span>
  </RouterLink>
</template>
