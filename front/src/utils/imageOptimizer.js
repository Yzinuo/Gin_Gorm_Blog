/**
 * Selects a pre-generated delivery variant from an asset manifest. R2 does not
 * provide Aliyun's x-oss-process query API, so string URLs are never mutated.
 */
export function optimizeImage(asset, options = {}) {
  if (!asset)
    return asset
  if (typeof asset === 'string')
    return asset
  const variants = Array.isArray(asset.srcset) ? [...asset.srcset].sort((a, b) => a.width - b.width) : []
  if (!variants.length)
    return asset.src || ''
  const target = Number(options.width) || variants.at(-1).width
  return (variants.find(item => item.width >= target) || variants.at(-1)).src
}

export const ImagePresets = {
  thumbnail: { width: 200 },
  card: { width: 700 },
  banner: { width: 1920 },
  full: { width: 1200 },
}

export function optimizeImageWithPreset(asset, preset = 'card') {
  return optimizeImage(asset, ImagePresets[preset] || ImagePresets.card)
}

export function generateSrcset(asset) {
  if (!asset || typeof asset === 'string')
    return ''
  return (asset.srcset || []).map(item => `${item.src} ${item.width}w`).join(', ')
}

export function preloadImage(url) {
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.onload = () => resolve()
    img.onerror = reject
    img.src = url
  })
}

export function preloadImages(urls) {
  return Promise.all(urls.map(url => preloadImage(url)))
}
