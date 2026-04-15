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
    expect(DEFAULT_BRANDING_HOME_CONTENT).toContain('brand-google-home__brand-mark')
    expect(DEFAULT_BRANDING_HOME_CONTENT).toContain('brand-google-home__badge-dot')
    expect(DEFAULT_BRANDING_CUSTOM_CSS).toContain('.brand-google-home')
    expect(DEFAULT_BRANDING_CUSTOM_CSS).toContain('.auth-layout-shell')
    expect(DEFAULT_BRANDING_CUSTOM_CSS).toContain('--brand-google-primary')
    expect(DEFAULT_BRANDING_CUSTOM_CSS).toContain('.brand-google-auth-note__marker')
    expect(DEFAULT_BRANDING_CUSTOM_CSS).toContain('.brand-google-footer__divider')
  })

  it('provides auth and footer snippets', () => {
    expect(DEFAULT_BRANDING_LOGIN_EXTRA_HTML).toContain('brand-google-auth-note')
    expect(DEFAULT_BRANDING_LOGIN_EXTRA_HTML).toContain('brand-google-auth-note__heading')
    expect(DEFAULT_BRANDING_LOGIN_EXTRA_HTML).toContain('Sign in')
    expect(DEFAULT_BRANDING_REGISTER_EXTRA_HTML).toContain('brand-google-auth-note')
    expect(DEFAULT_BRANDING_REGISTER_EXTRA_HTML).toContain('Create account')
    expect(DEFAULT_BRANDING_GLOBAL_FOOTER_HTML).toContain('brand-google-footer')
    expect(DEFAULT_BRANDING_GLOBAL_FOOTER_HTML).toContain('brand-google-footer__divider')
  })
})
