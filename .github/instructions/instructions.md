# Zene Project Overview
## Backend
- This project uses Golang for the backend, with the ncruces database/sql driver for SQLite.
- The backend is main.go and router.go in the root directory, and uses modularised services in the /core directory.
- The backend uses stdlib net/http for HTTP handling and routing.
- The backend presents http handlers compatible with the OpenSubsonic api spec: https://opensubsonic.netlify.app/docs/api-reference/
## Frontend
- The frontend is built with Vue 3, using the Composition API and TypeScript.
- The frontend uses @antfu/eslint-config for linting - this means no semi colons, and single quotes for strings.
- The frontend uses a custom backend fetch composable for making API requests: frontend/src/composables/backendFetch.ts
- The frontend uses vite-ssg for static site generation, with a custom vite.config.ts for configuration.
## Project Structure
- The project has a package.json in the root directory, which is used for managing local development dependencies and scripts.
- The project uses npm as the package manager.
- the caddy-baron package downloads Caddy and sets it up to serve the frontend and backend together at https://zene.localhost.
- the ffmpeg-baron package downloads FFmpeg for use in the backend in local development.
- the ffprobe-baron package downloads FFprobe for use in the backend in local development.
- During local development, the frontend is served by Vite at http://localhost:5173 and the backend at http://localhost:8080.
- When the application is built for production/preview, the built frontend dist is embedded in the backend binary.
## Pull requests
- Pull requests should be made against the main branch.
- A PR is not considered ready for review until the code is linted;
    - "npm run lint --workspace frontend" should be run and any linting errors resolved before opening a PR.
    - "npm run cspell" should be run before opening a PR to check for spelling errors - the project uses the en-GB dictionary.