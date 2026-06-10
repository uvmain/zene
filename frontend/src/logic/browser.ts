export function isBrowserChrome(): boolean {
  return /Chrome/.test(navigator.userAgent)
    && /Google Inc/.test(navigator.vendor)
}
