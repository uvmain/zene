export const isMobileNavOpen = ref(false)

export function toggleMobileNav() {
  isMobileNavOpen.value = !isMobileNavOpen.value
}

export function closeMobileNav() {
  isMobileNavOpen.value = false
}

export function openMobileNav() {
  isMobileNavOpen.value = true
}
