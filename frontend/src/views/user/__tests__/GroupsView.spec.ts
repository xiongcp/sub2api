import { describe, expect, it, vi, beforeEach } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'

import GroupsView from '../GroupsView.vue'

const { getAvailableSummary, showError } = vi.hoisted(() => ({
  getAvailableSummary: vi.fn(),
  showError: vi.fn(),
}))

const messages: Record<string, string> = {
  'userGroups.title': '可用分组',
  'userGroups.description': '查看当前账号可以使用的分组，以及每个分组的基础信息。',
  'userGroups.emptyTitle': '暂无可用分组',
  'userGroups.emptyDescription': '当前账号还没有可用分组，请联系管理员授权或开通订阅。',
  'userGroups.noDescription': '暂无分组说明',
  'userGroups.failedToLoad': '加载可用分组失败',
  'userGroups.access.public': '公开可用',
  'userGroups.access.exclusive': '已授权',
  'userGroups.access.subscription': '订阅可用',
  'userGroups.type.standard': '标准分组',
  'userGroups.type.subscription': '订阅分组',
  'userGroups.fields.platform': '平台',
  'userGroups.fields.type': '类型',
  'userGroups.fields.rate': '倍率',
  'userGroups.platforms.anthropic': 'Anthropic',
  'userGroups.platforms.openai': 'OpenAI',
  'userGroups.platforms.gemini': 'Gemini',
  'userGroups.platforms.antigravity': 'Antigravity',
}

vi.mock('@/api', () => ({
  userGroupsAPI: {
    getAvailableSummary,
  },
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({ showError }),
}))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => messages[key] ?? key,
    }),
  }
})

const AppLayoutStub = { template: '<div><slot /></div>' }
const GroupBadgeStub = {
  props: ['name'],
  template: '<div class="group-badge">{{ name }}</div>',
}
const EmptyStateStub = {
  props: ['title', 'description'],
  template: '<div class="empty-state">{{ title }}|{{ description }}</div>',
}

describe('user GroupsView', () => {
  beforeEach(() => {
    getAvailableSummary.mockReset()
    showError.mockReset()
  })

  it('renders sorted available group summaries', async () => {
    getAvailableSummary.mockResolvedValue([
      {
        id: 3,
        name: 'Subscription Beta',
        description: 'subscription-desc',
        platform: 'gemini',
        rate_multiplier: 0.8,
        subscription_type: 'subscription',
        access_scope: 'subscription',
      },
      {
        id: 2,
        name: 'Exclusive Alpha',
        description: 'exclusive-desc',
        platform: 'openai',
        rate_multiplier: 2,
        subscription_type: 'standard',
        access_scope: 'exclusive',
      },
      {
        id: 1,
        name: 'Public Gamma',
        description: null,
        platform: 'anthropic',
        rate_multiplier: 1.5,
        subscription_type: 'standard',
        access_scope: 'public',
      },
    ])

    const wrapper = mount(GroupsView, {
      global: {
        stubs: {
          AppLayout: AppLayoutStub,
          GroupBadge: GroupBadgeStub,
          EmptyState: EmptyStateStub,
          LoadingSpinner: true,
        },
      },
    })

    await flushPromises()

    const names = wrapper.findAll('.group-badge').map((node) => node.text())
    expect(names).toEqual(['Exclusive Alpha', 'Public Gamma', 'Subscription Beta'])
    expect(wrapper.text()).toContain('已授权')
    expect(wrapper.text()).toContain('公开可用')
    expect(wrapper.text()).toContain('订阅可用')
    expect(wrapper.text()).toContain('暂无分组说明')
    expect(showError).not.toHaveBeenCalled()
  })

  it('shows empty state when no groups are available', async () => {
    getAvailableSummary.mockResolvedValue([])

    const wrapper = mount(GroupsView, {
      global: {
        stubs: {
          AppLayout: AppLayoutStub,
          GroupBadge: GroupBadgeStub,
          EmptyState: EmptyStateStub,
          LoadingSpinner: true,
        },
      },
    })

    await flushPromises()

    expect(wrapper.find('.empty-state').text()).toContain('暂无可用分组')
    expect(showError).not.toHaveBeenCalled()
  })
})
