import { accentColour } from './store'

export function initializeAccentColour() {
  document.documentElement.style.setProperty('--main-colour', accentColour.value)
}

export function updateAccentColour(event: Event) {
  const color = (event.target as HTMLInputElement).value
  document.documentElement.style.setProperty('--main-colour', color)
  accentColour.value = color
}

export function resetAccentColour() {
  const defaultColour = 'hsl(22 95% 60%)'
  document.documentElement.style.setProperty('--main-colour', defaultColour)
  accentColour.value = defaultColour
}
