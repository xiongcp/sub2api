const CUSTOM_CSS_STYLE_ID = 'app-custom-css'

export function upsertCustomCss(css: string, doc: Document = document): void {
  if (!doc) return

  const normalized = css.trim()
  const existing = doc.getElementById(CUSTOM_CSS_STYLE_ID)
  if (!normalized) {
    existing?.remove()
    return
  }

  const styleEl =
    existing instanceof HTMLStyleElement
      ? existing
      : Object.assign(doc.createElement('style'), { id: CUSTOM_CSS_STYLE_ID })

  styleEl.textContent = normalized
  if (!styleEl.parentNode) {
    doc.head.appendChild(styleEl)
  }
}

export { CUSTOM_CSS_STYLE_ID }
