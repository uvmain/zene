<script setup lang="ts">
import { closeMobileNav } from '~/logic/navbar'

defineProps({
  routeName: { type: String, required: true },
  routeProp: { type: String, required: true },
  disableIndicator: { type: Boolean, required: false, default: false },
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
      'translate-x--8': currentRoute !== routeProp && !disableIndicator,
      'hover:translate-x-0': disableIndicator || (currentRoute !== routeProp && !currentRoute.startsWith(`${routeProp}/`)),
      'translate-x-0': (currentRoute === routeProp || currentRoute.startsWith(`${routeProp}/`)) && !disableIndicator,
    }"
    @click="handleLinkClick"
  >
    <icon-nrk-media-ffw
      v-if="!disableIndicator"
      class="opacity-0 size-8 transition-all duration-300"
      :class="{
        'opacity-100 text-main-400': currentRoute === routeProp || currentRoute.startsWith(`${routeProp}/`),
        'group-hover/navlink:opacity-50 text-main-400': currentRoute !== routeProp && !currentRoute.startsWith(`${routeProp}/`),
      }"
    />
    <div
      class="flex transition-all duration-300 items-center justify-start"
      :class="{
        'hover:translate-x-4': disableIndicator,
      }"
    >
      <slot name="icon" />
      <span>{{ routeName }}</span>
    </div>
  </RouterLink>
</template>
