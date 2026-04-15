import { describe, expect, it } from 'vitest'
import { CUSTOM_CSS_STYLE_ID, upsertCustomCss } from '@/utils/customCss'

describe('upsertCustomCss', () => {
  it('creates a style element when css is provided', () => {
    document.head.innerHTML = ''

    upsertCustomCss('body { color: red; }')

    const styleEl = document.getElementById(CUSTOM_CSS_STYLE_ID)
    expect(styleEl).not.toBeNull()
    expect(styleEl?.textContent).toBe('body { color: red; }')
  })

  it('updates an existing style element', () => {
    document.head.innerHTML = ''

    upsertCustomCss('body { color: red; }')
    upsertCustomCss('body { color: blue; }')

    const styles = document.querySelectorAll(`#${CUSTOM_CSS_STYLE_ID}`)
    expect(styles).toHaveLength(1)
    expect(styles[0]?.textContent).toBe('body { color: blue; }')
  })

  it('removes the style element when css is empty', () => {
    document.head.innerHTML = ''

    upsertCustomCss('body { color: red; }')
    upsertCustomCss('')

    expect(document.getElementById(CUSTOM_CSS_STYLE_ID)).toBeNull()
  })
})
