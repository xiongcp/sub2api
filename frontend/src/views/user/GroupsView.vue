<template>
  <AppLayout>
    <div class="space-y-6">
      <section class="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-800">
        <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">
          {{ t('userGroups.title') }}
        </h1>
        <p class="mt-2 max-w-3xl text-sm text-gray-600 dark:text-gray-300">
          {{ t('userGroups.description') }}
        </p>
      </section>

      <div v-if="loading" class="flex items-center justify-center py-12">
        <LoadingSpinner />
      </div>

      <EmptyState
        v-else-if="sortedGroups.length === 0"
        :title="t('userGroups.emptyTitle')"
        :description="t('userGroups.emptyDescription')"
      />

      <section v-else class="grid grid-cols-1 gap-4 xl:grid-cols-2">
        <article
          v-for="group in sortedGroups"
          :key="group.id"
          class="rounded-3xl border border-gray-200 bg-white p-5 shadow-sm transition-shadow hover:shadow-md dark:border-dark-700 dark:bg-dark-800"
        >
          <div class="flex items-start justify-between gap-3">
            <div class="min-w-0 space-y-3">
              <GroupBadge
                :name="group.name"
                :platform="group.platform"
                :subscription-type="group.subscription_type"
                :rate-multiplier="group.rate_multiplier"
              />
              <p class="text-sm leading-6 text-gray-600 dark:text-gray-300">
                {{ group.description || t('userGroups.noDescription') }}
              </p>
            </div>
            <span
              class="inline-flex shrink-0 rounded-full bg-blue-50 px-3 py-1 text-xs font-semibold text-blue-700 dark:bg-blue-500/10 dark:text-blue-300"
            >
              {{ accessScopeLabel(group.access_scope) }}
            </span>
          </div>

          <div class="mt-5 grid grid-cols-1 gap-3 sm:grid-cols-3">
            <div class="rounded-2xl bg-gray-50 px-4 py-3 dark:bg-dark-700/70">
              <p class="text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400">
                {{ t('userGroups.fields.platform') }}
              </p>
              <p class="mt-1 text-sm font-medium text-gray-900 dark:text-white">
                {{ platformLabel(group.platform) }}
              </p>
            </div>
            <div class="rounded-2xl bg-gray-50 px-4 py-3 dark:bg-dark-700/70">
              <p class="text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400">
                {{ t('userGroups.fields.type') }}
              </p>
              <p class="mt-1 text-sm font-medium text-gray-900 dark:text-white">
                {{ subscriptionTypeLabel(group.subscription_type) }}
              </p>
            </div>
            <div class="rounded-2xl bg-gray-50 px-4 py-3 dark:bg-dark-700/70">
              <p class="text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400">
                {{ t('userGroups.fields.rate') }}
              </p>
              <p class="mt-1 text-sm font-medium text-gray-900 dark:text-white">
                ×{{ formatMultiplier(group.rate_multiplier) }}
              </p>
            </div>
          </div>
        </article>
      </section>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { userGroupsAPI } from '@/api'
import { useAppStore } from '@/stores/app'
import AppLayout from '@/components/layout/AppLayout.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import GroupBadge from '@/components/common/GroupBadge.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import type { GroupPlatform, SubscriptionType, UserGroupAccessScope, UserGroupSummary } from '@/types'

const { t } = useI18n()
const appStore = useAppStore()

const loading = ref(false)
const groups = ref<UserGroupSummary[]>([])

const sortedGroups = computed(() =>
  [...groups.value].sort((a, b) => {
    if (a.subscription_type !== b.subscription_type) {
      return a.subscription_type === 'standard' ? -1 : 1
    }
    return a.name.localeCompare(b.name)
  })
)

const accessScopeLabel = (accessScope: UserGroupAccessScope) => t(`userGroups.access.${accessScope}`)
const subscriptionTypeLabel = (type: SubscriptionType) => t(`userGroups.type.${type}`)
const platformLabel = (platform: GroupPlatform) => t(`userGroups.platforms.${platform}`)
const formatMultiplier = (value: number) => value.toFixed(Number.isInteger(value) ? 0 : 2)

const loadGroups = async () => {
  loading.value = true
  try {
    groups.value = await userGroupsAPI.getAvailableSummary()
  } catch (error) {
    console.error('Failed to load available group summaries:', error)
    appStore.showError(t('userGroups.failedToLoad'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadGroups()
})
</script>
