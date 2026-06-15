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

export async function setHeroColourFromImage(imageElement: HTMLImageElement): Promise<void> {
  const options: ExtractionOptions = {
    colorCount: 10,
    quality: 10,
    worker: true,
    ignoreWhite: true,
    minSaturation: 0.1,
  }
  const palette = await getPalette(imageElement, options) ?? []
  const colours = palette.filter((colour) => {
    const hsl = colour.hsl()
    return hsl.l > 10 && hsl.l < 90 && hsl.s > 10
  }).sort((a, b) => {
    const aScore = (a.hsl().s * 1.5 + a.proportion * 1.0) * (a.isLight === true ? 1.5 : 1.0)
    const bScore = (b.hsl().s * 1.5 + b.proportion * 1.0) * (b.isLight === true ? 1.5 : 1.0)
    return bScore - aScore
  })

  const newColour = colours.length > 0 ? colours[0]?.toString() : accentColour.value

  document.documentElement.style.setProperty('--hero-colour', newColour)
}

export async function setAccentFromImage(imageElement: HTMLImageElement): Promise<void> {
  const options: ExtractionOptions = {
    colorCount: 10,
    quality: 10,
    worker: true,
    ignoreWhite: true,
    minSaturation: 0.1,
  }
  const palette = await getPalette(imageElement, options) ?? []
  const colours = palette.filter((colour) => {
    const hsl = colour.hsl()
    return hsl.l > 10 && hsl.l < 90 && hsl.s > 10
  }).sort((a, b) => {
    const aScore = (a.hsl().s * 1.5 + a.proportion * 1.0) * (a.isLight === true ? 1.5 : 1.0)
    const bScore = (b.hsl().s * 1.5 + b.proportion * 1.0) * (b.isLight === true ? 1.5 : 1.0)
    return bScore - aScore
  })

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
