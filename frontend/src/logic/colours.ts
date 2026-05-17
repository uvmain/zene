import type { ExtractionOptions } from 'colorthief'
import { getSwatches } from 'colorthief'
import { accentColour } from './store'

const DEFAULT_COLOUR = 'hsl(22 95% 60%)' as const

export function initializeAccentColour() {
  document.documentElement.style.setProperty('--main-colour', accentColour.value)
}

export function updateAccentColour(event: Event) {
  const color = (event.target as HTMLInputElement).value
  document.documentElement.style.setProperty('--main-colour', color)
  accentColour.value = color
}

export function resetAccentColour() {
  document.documentElement.style.setProperty('--main-colour', DEFAULT_COLOUR)
  accentColour.value = DEFAULT_COLOUR
}

export async function setAccentFromImage(imageElement: HTMLImageElement): Promise<void> {
  const options: ExtractionOptions = {
    colorCount: 10,
    quality: 5,
    colorSpace: 'oklch',
  }
  const swatches = await getSwatches(imageElement, options)
  let newColour: string
  if (swatches.Vibrant?.color) {
    newColour = swatches.Vibrant.color.toString()
  }
  else if (swatches.DarkVibrant?.color) {
    newColour = swatches.DarkVibrant.color.toString()
  }
  else if (swatches.LightVibrant?.color) {
    newColour = swatches.LightVibrant.color.toString()
  }
  else if (swatches.Muted?.color) {
    newColour = swatches.Muted.color.toString()
  }
  else {
    newColour = DEFAULT_COLOUR
  }
  document.documentElement.style.setProperty('--main-colour', newColour)
  accentColour.value = newColour
}
