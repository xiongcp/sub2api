import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { nextTick } from 'vue'

import AppHeader from '../AppHeader.vue'
import { useAppStore, useAuthStore } from '@/stores'

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn(),
  }),
  useRoute: () => ({
    name: 'Dashboard',
    params: {},
    meta: {
      title: 'Dashboard',
      description: '',
    },
  }),
}))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key,
    }),
  }
})

describe('AppHeader top banner', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  function mountHeader() {
    return mount(AppHeader, {
      global: {
        stubs: {
          AnnouncementBell: true,
          LocaleSwitcher: true,
          SubscriptionProgressMini: true,
          Icon: {
            template: '<span />',
          },
          RouterLink: {
            template: '<a><slot /></a>',
          },
        },
      },
    })
  }

  function seedStores(text = '充值联系 support@example.com') {
    const authStore = useAuthStore()
    authStore.token = 'token'
    authStore.user = {
      id: 1,
      email: 'user@example.com',
      username: 'tester',
      role: 'user',
      balance: 8.5,
    } as any

    const appStore = useAppStore()
    appStore.cachedPublicSettings = {
      top_banner_enabled: true,
      top_banner_text: text,
      custom_menu_items: [],
    } as any
  }

  it('已登录且配置开启时显示顶部横幅', () => {
    seedStores()

    const wrapper = mountHeader()

    expect(wrapper.text()).toContain('充值联系 support@example.com')
  })

  it('关闭后写入本地缓存并隐藏横幅', async () => {
    seedStores()

    const wrapper = mountHeader()
    await wrapper.get('[data-testid="top-banner-close"]').trigger('click')
    await nextTick()

    expect(localStorage.getItem('sub2api.top_banner_dismissed_signature')).toBe('充值联系 support@example.com')
    expect(wrapper.text()).not.toContain('充值联系 support@example.com')
  })

  it('文案变化后即使之前关闭过也会重新显示', async () => {
    seedStores()

    const wrapper = mountHeader()
    await wrapper.get('[data-testid="top-banner-close"]').trigger('click')

    const appStore = useAppStore()
    appStore.cachedPublicSettings = {
      top_banner_enabled: true,
      top_banner_text: '续费请联系 new-support@example.com',
      custom_menu_items: [],
    } as any
    await nextTick()

    expect(wrapper.text()).toContain('续费请联系 new-support@example.com')
  })
})
