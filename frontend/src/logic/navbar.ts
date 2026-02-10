let mobileNavOpen = false

export function isMobileNavOpen(): boolean {
  return mobileNavOpen
}

export function toggleMobileNav() {
  mobileNavOpen = !mobileNavOpen
}

export function closeMobileNav() {
  mobileNavOpen = false
}

export function openMobileNav() {
  mobileNavOpen = true
}
