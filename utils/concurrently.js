import { spawn } from 'node:child_process'
import path from 'node:path'
import process from 'node:process'
import { fileURLToPath } from 'node:url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

const firstArg = process.argv[2]
let caddyScript = ''
let backendScript = ''
let frontendScript = ''

switch (firstArg) {
  case 'preview':
    caddyScript = 'caddy:preview'
    backendScript = 'preview:go'
    frontendScript = ''
    break
  default:
    caddyScript = 'caddy'
    backendScript = 'dev:backend'
    frontendScript = 'dev:frontend'
    break
}

const components = [
  {
    name: 'caddy',
    path: path.resolve(__dirname, '../'),
    script: `pnpm run ${caddyScript}`,
    colour: 'green',
  },
  {
    name: 'backend',
    path: path.resolve(__dirname, '../'),
    script: `pnpm run ${backendScript}`,
    colour: 'cyan',
  },
  firstArg !== 'preview' && {
    name: 'frontend',
    path: path.resolve(__dirname, '../'),
    script: `pnpm run ${frontendScript}`,
    colour: 'yellow',
  },
].filter(Boolean)

const activeChildren = new Map()
let isShuttingDown = false
let shutdownExitCode = 0

function killChildProcess(childProcess) {
  if (!childProcess?.pid) {
    return
  }

  if (process.platform === 'win32') {
    const killer = spawn('taskkill', [
      '/pid',
      String(childProcess.pid),
      '/T',
      '/F',
    ], {
      stdio: 'ignore',
      windowsHide: true,
    })
    killer.on('error', () => {
      childProcess.kill('SIGTERM')
    })
    return
  }

  childProcess.kill('SIGTERM')
}

function shutdownAll(exitedName, code, signal) {
  if (isShuttingDown) {
    return
  }

  isShuttingDown = true
  shutdownExitCode = typeof code === 'number'
    ? code
    : 1

  const reason = typeof code === 'number'
    ? `code ${code}`
    : `signal ${signal || 'unknown'}`

  colourConsole('red', 'manager', `${exitedName} exited (${reason}); stopping remaining processes.`)

  for (const [name, childProcess] of activeChildren.entries()) {
    if (name === exitedName) {
      continue
    }

    killChildProcess(childProcess)
  }
}

components.forEach((component) => {
  const { name, path, script, colour } = component

  const childProcess = spawn(script, {
    cwd: path,
    shell: true,
    stdio: 'overlapped',
  })

  activeChildren.set(name, childProcess)

  childProcess.stdout.on('data', (data) => {
    colourConsole(colour, name, data.toString())
  })

  childProcess.stderr.on('data', (data) => {
    colourConsole(colour, name, data.toString())
  })

  childProcess.on('error', (error) => {
    colourConsole(colour, name, `Error running npm script: ${error}`)
    shutdownAll(name)
  })

  childProcess.on('close', (code, signal) => {
    activeChildren.delete(name)

    if (code === 0) {
      colourConsole(colour, name, 'Script completed successfully.')
    }
    else {
      const reason = typeof code === 'number'
        ? `code ${code}`
        : `signal ${signal || 'unknown'}`
      colourConsole(colour, name, `Script exited with ${reason}.`)
    }

    if (!isShuttingDown) {
      shutdownAll(name, code, signal)
    }

    if (isShuttingDown && activeChildren.size === 0) {
      process.exit(shutdownExitCode)
    }
  })
})

process.on('SIGINT', () => {
  shutdownAll('manager', 130, 'SIGINT')
})

process.on('SIGTERM', () => {
  shutdownAll('manager', 143, 'SIGTERM')
})

export function colourConsole(colour, name, message) {
  if (message.trim() === '') {
    return
  }

  const paddedName = `[${name}]`

  // don't display empty lines, add padding to newlines
  message = message.split('\n').filter(line => line.trim() !== '').join('\n')
  message = message.replace(/\n/g, `\n${' '.repeat(name.length + 3)} `)

  if (colour === 'red')
    console.log('\u001B[31m', `${paddedName}\u001B[0m`, message)
  else if (colour === 'green')
    console.log('\u001B[32m', `${paddedName}\u001B[0m`, message)
  else if (colour === 'yellow')
    console.log('\u001B[33m', `${paddedName}\u001B[0m`, message)
  else if (colour === 'blue')
    console.log('\u001B[34m', `${paddedName}\u001B[0m`, message)
  else if (colour === 'purple')
    console.log('\u001B[35m', `${paddedName}\u001B[0m`, message)
  else if (colour === 'cyan')
    console.log('\u001B[36m', `${paddedName}\u001B[0m`, message)
  else console.log(colour, `${paddedName}`, message)
}
