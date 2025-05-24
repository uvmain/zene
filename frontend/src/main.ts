import { ViteSSG } from 'vite-ssg/single-page'

import App from './App.vue'
import 'virtual:uno.css'

export const createApp = ViteSSG(App as Component)
