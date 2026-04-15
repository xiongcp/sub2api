import { describe, expect, it } from 'vitest'
import {
  DEFAULT_BRANDING_CUSTOM_CSS,
  DEFAULT_BRANDING_GLOBAL_FOOTER_HTML,
  DEFAULT_BRANDING_HOME_CONTENT,
  DEFAULT_BRANDING_LOGIN_EXTRA_HTML,
  DEFAULT_BRANDING_REGISTER_EXTRA_HTML,
} from '@/constants/defaultBranding'

describe('defaultBranding', () => {
  it('provides a non-empty home template and styling', () => {
    expect(DEFAULT_BRANDING_HOME_CONTENT).toContain('brand-google-home')
    expect(DEFAULT_BRANDING_CUSTOM_CSS).toContain('.brand-google-home')
    expect(DEFAULT_BRANDING_CUSTOM_CSS).toContain('.auth-layout-shell')
  })

  it('provides auth and footer snippets', () => {
    expect(DEFAULT_BRANDING_LOGIN_EXTRA_HTML).toContain('brand-google-auth-note')
    expect(DEFAULT_BRANDING_REGISTER_EXTRA_HTML).toContain('brand-google-auth-note')
    expect(DEFAULT_BRANDING_GLOBAL_FOOTER_HTML).toContain('brand-google-footer')
  })
})
