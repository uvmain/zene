const isMobileNavOpen = ref(false)

export function useNavbar() {
  const toggleMobileNav = () => {
    isMobileNavOpen.value = !isMobileNavOpen.value
  }

  const closeMobileNav = () => {
    isMobileNavOpen.value = false
  }

  const openMobileNav = () => {
    isMobileNavOpen.value = true
  }

  return {
    isMobileNavOpen,
    toggleMobileNav,
    closeMobileNav,
    openMobileNav,
  }
}
