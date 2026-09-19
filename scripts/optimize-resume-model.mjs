import { spawnSync } from 'node:child_process'
import { fileURLToPath } from 'node:url'
import path from 'node:path'

const directory = path.dirname(fileURLToPath(import.meta.url))
const input = path.resolve(process.argv[2] || path.join(directory, '../front/public/resume/resume-ready.glb'))
const output = path.resolve(process.argv[3] || path.join(directory, '../front/public/resume/resume-ready.optimized.glb'))
const cli = path.join(directory, 'node_modules/@gltf-transform/cli/bin/cli.js')
const result = spawnSync(process.execPath, [
  cli, 'optimize', input, output,
  '--compress', 'meshopt',
  '--texture-compress', 'webp',
  '--texture-size', '2048',
  '--flatten', 'false',
  '--join', 'false',
  '--palette', 'false',
  '--simplify', 'false',
], { stdio: 'inherit' })

if (result.error)
  throw result.error
process.exitCode = result.status ?? 1
