# Testing Guide for Zene Frontend

This project uses **Vitest** for unit and integration testing with Vue 3 + TypeScript.

## Overview

- **Test Runner**: Vitest (fast Vite-native test runner)
- **Test Environment**: happy-dom (lightweight DOM implementation)
- **Component Testing**: @vue/test-utils
- **API Mocking**: MSW (Mock Service Worker)
- **Assertions**: Built-in Vitest + @testing-library/jest-dom

## Quick Start

```bash
# Run all tests
npm run test

# Run tests once (CI mode)
npm run test:run

# Run tests with coverage report
npm run test:coverage

# Run tests with UI
npm run test:ui
```

## Project Structure

```
frontend/
├── test/
│   ├── setup.ts           # Global test setup (MSW, mocks)
│   └── utils.ts           # Test utilities and helpers
├── src/
│   ├── components/
│   │   └── __tests__/     # Component tests
│   ├── composables/
│   │   └── __tests__/     # Composable tests
│   └── routes/
│       └── __tests__/     # Route tests
└── vite.config.ts         # Vitest configuration
```

## Writing Tests

### Component Tests

```typescript
import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import MyComponent from '../MyComponent.vue'

describe('MyComponent', () => {
  it('should render correctly', () => {
    const wrapper = mount(MyComponent, {
      props: { 
        // Required props here
      },
      global: {
        mocks: {
          $router: { push: vi.fn() },
          $route: { path: '/', params: {}, query: {} },
        },
        stubs: {
          'RouterLink': true,
          'RouterView': true,
        },
      },
    })
    
    expect(wrapper.exists()).toBe(true)
    expect(wrapper.text()).toContain('Expected text')
  })
})
```

### Composable Tests

```typescript
import { describe, it, expect, vi } from 'vitest'
import { useMyComposable } from '../useMyComposable'

describe('useMyComposable', () => {
  it('should return expected functionality', () => {
    const { someFunction, someState } = useMyComposable()
    
    expect(someFunction).toBeDefined()
    expect(someState.value).toBe(expectedValue)
  })
})
```

### Testing with API Calls

API calls are automatically mocked using MSW. The setup handles common endpoints:

```typescript
// This will be automatically mocked
const albums = await fetch('/api/albums')
// Returns mock data defined in test/setup.ts
```

To add custom API mocks for specific tests:

```typescript
import { server } from '../../../test/setup'
import { http, HttpResponse } from 'msw'

it('should handle custom API response', async () => {
  server.use(
    http.get('/api/custom-endpoint', () => {
      return HttpResponse.json({ custom: 'data' })
    })
  )
  
  // Your test here
})
```

## API Mocking with MSW

The test setup includes comprehensive API mocking for:

- **Authentication**: `/api/check-session`, `/api/logout`
- **Music Data**: `/api/albums`, `/api/artists`, `/api/tracks`, `/api/genres`
- **Search**: `/api/search`
- **Settings**: `/api/settings`
- **Users**: `/api/users`, `/api/current-user`
- **Media**: Album/artist art endpoints

### Example: Mocking GET/POST Requests

```typescript
// GET request example (already handled in setup)
http.get('/api/albums', () => {
  return HttpResponse.json([
    { id: 1, title: 'Test Album', artist: 'Test Artist' }
  ])
})

// POST request example
http.post('/api/albums', async ({ request }) => {
  const body = await request.json()
  return HttpResponse.json({ id: 1, ...body })
})

// Error response example
http.get('/api/error-endpoint', () => {
  return new HttpResponse(null, { status: 500 })
})
```

## Testing Different Component Types

### Components with Props

```typescript
const mockAlbum = {
  id: 1,
  title: 'Test Album',
  artist: 'Test Artist',
  musicbrainz_album_id: 'test-id'
}

const wrapper = mount(Album, {
  props: { album: mockAlbum }
})
```

### Components with Router

```typescript
const wrapper = mount(MyComponent, {
  global: {
    mocks: {
      $router: { push: vi.fn() },
      $route: { path: '/test', params: { id: '1' } }
    }
  }
})
```

### Components with User Interactions

```typescript
it('should handle click events', async () => {
  const wrapper = mount(MyComponent)
  const button = wrapper.find('button')
  
  await button.trigger('click')
  
  expect(wrapper.emitted('click')).toBeTruthy()
})
```

## Coverage Reporting

Run coverage with:

```bash
npm run test:coverage
```

Coverage reports are generated in:
- Console output (text format)
- `coverage/index.html` (HTML format)
- `coverage/coverage.json` (JSON format)

## Best Practices

1. **Focus on behavior, not implementation**
2. **Test user interactions and edge cases**
3. **Mock external dependencies (APIs, complex composables)**
4. **Use descriptive test names**
5. **Keep tests simple and focused**
6. **Test error states and loading states**

## Troubleshooting

### Common Issues

1. **Router injection errors**: Add router mocks to global mocks
2. **Missing props**: Provide required props in mount options
3. **API calls not mocked**: Check MSW handlers in `test/setup.ts`
4. **Composable circular dependencies**: Mock complex composables

### Debug Tips

```typescript
// Debug component HTML
console.log(wrapper.html())

// Debug component data
console.log(wrapper.vm.$data)

// Debug MSW requests
server.events.on('request:start', ({ request }) => {
  console.log('MSW intercepted:', request.method, request.url)
})
```

## Generated Tests

The project includes auto-generated tests for all components and composables. These provide:

- Basic rendering tests
- Vue instance validation
- Props validation (where applicable)

To regenerate tests:

```bash
node generate-tests.mjs
```

## Adding New Tests

1. Create test file following naming convention: `ComponentName.test.ts`
2. Place in appropriate `__tests__` directory
3. Follow the examples in existing tests
4. Run tests to ensure they pass
5. Add to CI pipeline if needed

## CI/CD Integration

Add to your CI pipeline:

```yaml
- name: Run tests
  run: npm run test:run

- name: Generate coverage
  run: npm run test:coverage
```

For more detailed information, refer to the [Vitest documentation](https://vitest.dev/) and [Vue Test Utils documentation](https://test-utils.vuejs.org/).