import fs from 'node:fs'
import process from 'node:process'

const args = process.argv.slice(2)

if (args.length < 2) {
  console.error('Usage: node copy.js <source> <destination>')
  process.exit(1)
}

const [src, dest] = args

fs.cp(src, dest, { recursive: true }, (err) => {
  if (err) {
    console.error(`Error copying frontend dist: ${err}`)
  }
  console.log('Frontend dist copied to webapi dist')
})
