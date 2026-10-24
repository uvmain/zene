/* eslint-disable unused-imports/no-unused-vars */
import { rmSync } from 'node:fs'
import process from 'node:process'

const args = process.argv.slice(2)

for (const path of args) {
  try {
    rmSync(path, { force: true, recursive: true })
  }
  catch (error) {
    console.error(`Failed to remove ${path}:`, error)
  }
}
