import type { ExtractionOptions } from 'colorthief'
import { getPalette } from 'colorthief'
import { accentColour } from './store'

const DEFAULT_COLOUR: string = 'hsl(22 95% 60%)' as const

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
    colorCount: 6,
  }
  const palette = await getPalette(imageElement, options) ?? []
  const colours = palette.filter((colour) => {
    const hsl = colour.hsl()
    return hsl.s > 20 && hsl.l > 10 && hsl.l < 90
  }).sort((a, b) => b.population - a.population)

  const newColour = colours[0]?.toString() ?? accentColour.value

  document.documentElement.style.setProperty('--main-colour', newColour)
}
