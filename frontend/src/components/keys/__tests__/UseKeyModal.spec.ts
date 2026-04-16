import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key
  })
}))

vi.mock('@/composables/useClipboard', () => ({
  useClipboard: () => ({
    copyToClipboard: vi.fn().mockResolvedValue(true)
  })
}))

import UseKeyModal from '../UseKeyModal.vue'

describe('UseKeyModal', () => {
  it('renders updated GPT-5.4 mini/nano names in OpenCode config', async () => {
    const wrapper = mount(UseKeyModal, {
      props: {
        show: true,
        apiKey: 'sk-test',
        baseUrl: 'https://example.com/v1',
        dynamicBaseUrl: 'https://dynamic.example.com/v1',
        platform: 'openai',
        usageGuideContent: {
          description: '',
          note: '',
          no_group_title: '',
          no_group_description: '',
          openai: {
            description: '',
            config_toml_hint: '',
            note: '',
            note_windows: '',
            model_comment: '',
            claude_note: '',
            gemini_note: ''
          },
          gemini: {
            description: '',
            config_toml_hint: '',
            note: '',
            note_windows: '',
            model_comment: '',
            claude_note: '',
            gemini_note: ''
          },
          antigravity: {
            description: '',
            config_toml_hint: '',
            note: '',
            note_windows: '',
            model_comment: '',
            claude_note: '',
            gemini_note: ''
          },
          opencode: {
            hint: 'server hint'
          }
        }
      },
      global: {
        stubs: {
          BaseDialog: {
            template: '<div><slot /><slot name="footer" /></div>'
          },
          Icon: {
            template: '<span />'
          }
        }
      }
    })

    const opencodeTab = wrapper.findAll('button').find((button) =>
      button.text().includes('keys.useKeyModal.cliTabs.opencode')
    )

    expect(opencodeTab).toBeDefined()
    await opencodeTab!.trigger('click')
    await nextTick()

    const codeBlock = wrapper.find('pre code')
    expect(codeBlock.exists()).toBe(true)
    expect(codeBlock.text()).toContain('"name": "GPT-5.4 Mini"')
    expect(codeBlock.text()).toContain('"name": "GPT-5.4 Nano"')
    expect(codeBlock.text()).toContain('"baseURL": "https://dynamic.example.com/v1"')
    expect(wrapper.text()).toContain('server hint')
  })

  it('prefers server-provided no-group copy over i18n defaults', () => {
    const wrapper = mount(UseKeyModal, {
      props: {
        show: true,
        apiKey: 'sk-test',
        baseUrl: 'https://example.com/v1',
        platform: null,
        usageGuideContent: {
          description: '',
          note: '',
          no_group_title: 'Assign a group first',
          no_group_description: 'The latest copy comes from the server.',
          openai: {
            description: '',
            config_toml_hint: '',
            note: '',
            note_windows: '',
            model_comment: '',
            claude_note: '',
            gemini_note: ''
          },
          gemini: {
            description: '',
            config_toml_hint: '',
            note: '',
            note_windows: '',
            model_comment: '',
            claude_note: '',
            gemini_note: ''
          },
          antigravity: {
            description: '',
            config_toml_hint: '',
            note: '',
            note_windows: '',
            model_comment: '',
            claude_note: '',
            gemini_note: ''
          },
          opencode: {
            hint: ''
          }
        }
      },
      global: {
        stubs: {
          BaseDialog: {
            template: '<div><slot /><slot name="footer" /></div>'
          },
          Icon: {
            template: '<span />'
          }
        }
      }
    })

    expect(wrapper.text()).toContain('Assign a group first')
    expect(wrapper.text()).toContain('The latest copy comes from the server.')
  })
})
