// Type definitions for butterchurn
// Project: https://github.com/jberg/butterchurn

declare module 'butterchurn' {
  export interface ButterchurnVisualizerOptions {
    width?: number
    height?: number
    pixelRatio?: number
    textureRatio?: number
    meshWidth?: number
    meshHeight?: number
    canvas?: HTMLCanvasElement
    offscreenCanvas?: boolean
    hideFps?: boolean
    hideControls?: boolean
    [key: string]: any
  }

  export interface ButterchurnVisualizer {
    connectAudio: (audioNode: AudioNode) => void
    disconnectAudio: (audioNode: AudioNode) => void
    loadPreset: (preset: object, blendTime: number) => void
    setRendererSize: (width: number, height: number) => void
    render: () => void
    destroy: () => void
    /** Other methods and properties may exist */
  }

  export interface Butterchurn {
    createVisualizer: (
      context: AudioContext,
      canvas: HTMLCanvasElement,
      opts?: ButterchurnVisualizerOptions,
    ) => ButterchurnVisualizer
  }

  const butterchurn: Butterchurn
  export default butterchurn
}

declare module 'butterchurn-presets' {
  export interface ButterchurnPreset {
    [key: string]: any
  }

  export interface ButterchurnPresets {
    getPresets: () => { [presetName: string]: ButterchurnPreset }
  }

  const butterchurnPresets: ButterchurnPresets
  export default butterchurnPresets
}
