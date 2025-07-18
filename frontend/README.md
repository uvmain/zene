# uvite

This template should help get you started developing with Vue 3 in Vite.

## Recommended IDE Setup

[VSCode](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Type Support for `.vue` Imports in TS

TypeScript cannot handle type information for `.vue` imports by default, so we replace the `tsc` CLI with `vue-tsc` for type checking. In editors, we need [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) to make the TypeScript language service aware of `.vue` types.

## Customize configuration

See [Vite Configuration Reference](https://vitejs.dev/config/).

## Project Setup

```sh
npm install
```

### Compile and Hot-Reload for Development

```sh
npm run dev
```

### Type-Check, Compile and Minify for Production

```sh
npm run build
```

### Lint with [ESLint](https://eslint.org/)

```sh
npm run lint
```

### Testing with [Vitest](https://vitest.dev/)

The project includes a comprehensive testing framework using Vitest with Vue Test Utils and MSW for API mocking.

#### Test Scripts

```sh
# Run tests in watch mode (development)
npm run test
```
Starts Vitest in watch mode, automatically re-running tests when files change. Perfect for development and TDD workflows.

```sh
# Run tests once (CI/CD)
npm run test:run
```
Runs all tests once and exits. Ideal for CI/CD pipelines and pre-commit hooks.

```sh
# Generate coverage report
npm run test:coverage
```
Runs all tests and generates a detailed coverage report in HTML and JSON formats. Reports are saved to the `coverage/` directory.

```sh
# Run with UI interface
npm run test:ui
```
Opens Vitest's web-based UI for an interactive testing experience with real-time results, test filtering, and debugging capabilities.

#### Additional Testing Resources

- **Testing Guide**: See `TESTING.md` for comprehensive testing patterns and examples
- **Setup Details**: See `VITEST_SETUP_COMPLETE.md` for framework configuration details
- **Test Files**: Component tests in `src/components/__tests__/`, composable tests in `src/composables/__tests__/`
