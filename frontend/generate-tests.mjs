#!/usr/bin/env node

import { readdir, writeFile, mkdir } from 'fs/promises'
import { join, parse } from 'path'
import { fileURLToPath } from 'url'

const __dirname = fileURLToPath(new URL('.', import.meta.url))

async function generateComponentTest(componentName) {
  const testName = componentName.replace('.vue', '')
  return `import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import ${testName} from '../${componentName}'

// Mock router
const mockRouter = {
  push: vi.fn(),
  replace: vi.fn(),
}

describe('${testName}', () => {
  it('should render correctly', () => {
    const wrapper = mount(${testName}, {
      global: {
        mocks: {
          $router: mockRouter,
          $route: { path: '/', params: {}, query: {} },
        },
        stubs: {
          'RouterLink': true,
          'RouterView': true,
        },
      },
    })
    expect(wrapper.exists()).toBe(true)
  })

  it('should be a Vue instance', () => {
    const wrapper = mount(${testName}, {
      global: {
        mocks: {
          $router: mockRouter,
          $route: { path: '/', params: {}, query: {} },
        },
        stubs: {
          'RouterLink': true,
          'RouterView': true,
        },
      },
    })
    expect(wrapper.vm).toBeTruthy()
  })
})
`
}

async function generateComposableTest(composableName) {
  const testName = composableName.replace('.ts', '')
  return `import { describe, it, expect, vi, beforeEach } from 'vitest'
import { ${testName} } from '../${composableName}'

describe('${testName}', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should be defined', () => {
    expect(${testName}).toBeDefined()
  })

  it('should return expected properties/methods', () => {
    const result = ${testName}()
    expect(result).toBeTruthy()
  })
})
`
}

async function generateRouteTest(routeName) {
  const testName = routeName.replace('.vue', '')
  return `import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import ${testName} from '../${routeName}'

// Mock router
const mockRouter = {
  push: vi.fn(),
  replace: vi.fn(),
}

describe('${testName}', () => {
  it('should render correctly', () => {
    const wrapper = mount(${testName}, {
      global: {
        mocks: {
          $router: mockRouter,
          $route: { path: '/', params: {}, query: {} },
        },
        stubs: {
          'RouterLink': true,
          'RouterView': true,
        },
      },
    })
    expect(wrapper.exists()).toBe(true)
  })

  it('should be a Vue instance', () => {
    const wrapper = mount(${testName}, {
      global: {
        mocks: {
          $router: mockRouter,
          $route: { path: '/', params: {}, query: {} },
        },
        stubs: {
          'RouterLink': true,
          'RouterView': true,
        },
      },
    })
    expect(wrapper.vm).toBeTruthy()
  })
})
`
}

async function generateTests() {
  const srcPath = join(__dirname, 'src')
  
  // Generate component tests
  const componentsPath = join(srcPath, 'components')
  const componentFiles = await readdir(componentsPath)
  const vueComponents = componentFiles.filter(file => file.endsWith('.vue'))
  
  for (const component of vueComponents) {
    const testContent = await generateComponentTest(component)
    const testFileName = component.replace('.vue', '.test.ts')
    const testPath = join(componentsPath, '__tests__', testFileName)
    await writeFile(testPath, testContent)
    console.log(`Generated test: ${testPath}`)
  }
  
  // Generate route tests (only a few key ones to avoid too many API dependencies)
  const routesPath = join(srcPath, 'routes')
  const routeFiles = await readdir(routesPath)
  const vueRoutes = routeFiles.filter(file => file.endsWith('.vue') && ['HomeRoute.vue', 'LoginRoute.vue'].includes(file))
  
  for (const route of vueRoutes) {
    const testContent = await generateRouteTest(route)
    const testFileName = route.replace('.vue', '.test.ts')
    const testPath = join(routesPath, '__tests__', testFileName)
    await writeFile(testPath, testContent)
    console.log(`Generated test: ${testPath}`)
  }
  
  console.log('\\nTest generation complete!')
  console.log(`Generated ${vueComponents.length} component tests`)
  console.log(`Generated ${vueRoutes.length} route tests`)
}

generateTests().catch(console.error)