const allowedTags = new Set([
  'a', 'blockquote', 'br', 'code', 'del', 'em', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6',
  'hr', 'img', 'li', 'ol', 'p', 'pre', 'strong', 'table', 'tbody', 'td', 'th', 'thead',
  'tr', 'ul',
])

const allowedAttributes = {
  a: new Set(['href', 'title']),
  code: new Set(['class']),
  h1: new Set(['id']),
  h2: new Set(['id']),
  h3: new Set(['id']),
  h4: new Set(['id']),
  h5: new Set(['id']),
  h6: new Set(['id']),
  img: new Set(['alt', 'height', 'src', 'title', 'width']),
  td: new Set(['colspan', 'rowspan']),
  th: new Set(['colspan', 'rowspan']),
}

const removeWithContent = new Set(['base', 'embed', 'iframe', 'link', 'meta', 'object', 'script', 'style'])

function isSafeUrl(value) {
  try {
    const url = new URL(value, window.location.origin)
    return url.protocol === 'http:' || url.protocol === 'https:'
  }
  catch {
    return false
  }
}

export function safeExternalUrl(value) {
  return typeof value === 'string' && isSafeUrl(value) ? value : ''
}

// Markdown rendering still needs semantic HTML, but never raw event handlers,
// executable elements, inline styles, or unsafe URLs from stored content.
export function sanitizeHtml(value) {
  if (typeof value !== 'string' || !value)
    return ''

  const template = document.createElement('template')
  template.innerHTML = value

  for (const element of [...template.content.querySelectorAll('*')]) {
    const tag = element.tagName.toLowerCase()
    if (!allowedTags.has(tag)) {
      if (removeWithContent.has(tag))
        element.remove()
      else
        element.replaceWith(...element.childNodes)
      continue
    }

    const permitted = allowedAttributes[tag] || new Set()
    for (const attribute of [...element.attributes]) {
      const name = attribute.name.toLowerCase()
      if (!permitted.has(name)) {
        element.removeAttribute(attribute.name)
        continue
      }
      if ((name === 'href' || name === 'src') && !isSafeUrl(attribute.value))
        element.removeAttribute(attribute.name)
    }

    if (tag === 'a' && element.hasAttribute('href'))
      element.setAttribute('rel', 'noopener noreferrer')
  }

  return template.innerHTML
}
