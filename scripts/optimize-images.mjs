import path from 'node:path'
import { fileURLToPath } from 'node:url'
import sharp from 'sharp'

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..')

for (const name of ['Before', 'After']) {
  const input = path.join(root, `front/public/images/${name}.png`)
  for (const width of [960, 1672]) {
    await sharp(input)
      .resize({ width, withoutEnlargement: true })
      .webp({ quality: 82, smartSubsample: true })
      .toFile(path.join(root, `front/public/images/${name}-${width}.webp`))
  }
}

await sharp(path.join(root, 'front/public/favicon.png')).resize(32, 32).png({ compressionLevel: 9 }).toFile(path.join(root, 'front/public/favicon-32.png'))
await sharp(path.join(root, 'front/public/favicon.png')).resize(180, 180).png({ compressionLevel: 9, palette: true, quality: 88 }).toFile(path.join(root, 'front/public/apple-touch-icon.png'))
