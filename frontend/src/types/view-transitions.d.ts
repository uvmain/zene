// Type declarations for View Transitions API
interface Document {
  startViewTransition?: (callback?: () => void | Promise<void>) => ViewTransition
}

interface ViewTransition {
  finished: Promise<void>
  ready: Promise<void>
  updateCallbackDone: Promise<void>
  skipTransition: () => void
}
