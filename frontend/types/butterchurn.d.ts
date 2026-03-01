// Type definitions for butterchurn
// Project: https://github.com/jberg/butterchurn

declare module 'butterchurn' {
  export interface VisualizerOptions {
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
    deterministic?: boolean
    testMode?: boolean
    onlyUseWASM?: boolean
    [key: string]: any
  }

  export class Visualizer {
    constructor(audioContext: AudioContext, canvas: HTMLCanvasElement, opts: VisualizerOptions)

    // Properties
    opts: VisualizerOptions
    rng: any
    deterministicMode: boolean
    audio: any
    internalCanvas: HTMLCanvasElement | OffscreenCanvas
    gl: WebGLRenderingContext | WebGL2RenderingContext | null
    outputGl: CanvasRenderingContext2D | null
    renderer: any
    audioNode?: AudioNode

    // Default values
    baseValsDefaults: Record<string, any>
    shapeBaseValsDefaults: Record<string, any>
    waveBaseValsDefaults: Record<string, any>
    qs: string[]
    ts: string[]
    globalPerFrameVars: string[]
    globalPerPixelVars: string[]
    globalShapeVars: string[]
    shapeBaseVars: string[]
    globalWaveVars: string[]

    // Methods
    loseGLContext(): void
    connectAudio(audioNode: AudioNode): void
    disconnectAudio(audioNode: AudioNode): void
    createQVars(): Record<string, any>
    createTVars(): Record<string, any>
    createPerFramePool(baseVals: Record<string, any>): Record<string, any>
    createPerPixelPool(baseVals: Record<string, any>): Record<string, any>
    createCustomShapePerFramePool(baseVals: Record<string, any>): Record<string, any>
    createCustomWavePerFramePool(baseVals: Record<string, any>): Record<string, any>
    loadPreset(presetMap: object, blendTime?: number): Promise<void>
    loadJSPreset(preset: object, blendTime: number): void
    loadExtraImages(imageData: any): void
    setRendererSize(width: number, height: number, opts?: object): void
    setInternalMeshSize(width: number, height: number): void
    setOutputAA(useAA: boolean): void
    setCanvas(canvas: HTMLCanvasElement): void
    render(opts?: object): any
    launchSongTitleAnim(text: string): void
    toDataURL(): string
    warpBufferToDataURL(): string

    // Static methods
    static overrideDefaultVars(baseValsDefaults: Record<string, any>, baseVals: Record<string, any>): Record<string, any>
    static makeShapeResetPool(pool: Record<string, any>, variables: string[], idx: number): Record<string, any>
    static base64ToArrayBuffer(base64: string): ArrayBuffer
  }

  export interface Butterchurn {
    createVisualizer: (
      context: AudioContext,
      canvas: HTMLCanvasElement,
      opts?: VisualizerOptions,
    ) => Visualizer
  }

  const butterchurn: Butterchurn
  export default butterchurn
}
