import type { ExtractionOptions } from 'colorthief'
import { getPalette } from 'colorthief'
import { accentColour } from '~/stores/main'

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
    return hsl.l > 5 && hsl.l < 95 && hsl.s > 10 && colour.population > 5
  }).sort((a, b) => b.hsl().s - a.hsl().s)

  const newColour = colours.length > 0 ? colours[0]?.toString() : accentColour.value

  document.documentElement.style.setProperty('--main-colour', newColour)
}

export async function setAccentFromImageUrl(imageUrl: string): Promise<void> {
  const img = new Image()
  img.crossOrigin = 'anonymous'
  img.src = imageUrl
  img.onload = async () => {
    await setAccentFromImage(img)
  }
}
