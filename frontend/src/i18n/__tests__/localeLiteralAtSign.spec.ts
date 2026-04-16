import { describe, expect, it } from 'vitest'

import en from '../locales/en'
import zh from '../locales/zh'

describe('locale literal @ rendering', () => {
  it('stores english top banner placeholder using literal interpolation syntax', () => {
    expect(en.admin.settings.site.topBannerTextPlaceholder).toBe(
      "e.g., For recharge support, contact Telegram: {'@'}support"
    )
  })

  it('stores chinese top banner placeholder using literal interpolation syntax', () => {
    expect(zh.admin.settings.site.topBannerTextPlaceholder).toBe(
      "例如：充值请联系 Telegram: {'@'}support"
    )
  })
})
