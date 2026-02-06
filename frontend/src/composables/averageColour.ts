import type { ColorInstance } from 'color'
import Color from 'color'

export function getImageColour(imageObject: HTMLImageElement): ColorInstance {
  // draw the image to one pixel and let the browser find the dominant color
  const ctx = document.createElement('canvas').getContext('2d')!
  ctx.canvas.width = 1
  ctx.canvas.height = 1

  ctx.drawImage(imageObject, 0, 0, 1, 1)

  // get pixel color
  const i = ctx.getImageData(0, 0, 1, 1).data

  return Color.rgb(i[0], i[1], i[2])
}
