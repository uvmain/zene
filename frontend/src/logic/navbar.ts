export const isMobileNavOpen = ref(false)

export function toggleMobileNav() {
  isMobileNavOpen.value = !isMobileNavOpen.value
}

export function closeMobileNav() {
  isMobileNavOpen.value = false
}
