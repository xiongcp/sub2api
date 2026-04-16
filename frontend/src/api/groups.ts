/**
 * User Groups API endpoints (non-admin)
 * Handles group-related operations for regular users
 */

import { apiClient } from './client'
import type { Group, UserGroupSummary } from '@/types'

/**
 * Get available groups that the current user can bind to API keys
 * This returns groups based on user's permissions:
 * - Standard groups: public (non-exclusive) or explicitly allowed
 * - Subscription groups: user has active subscription
 * @returns List of available groups
 */
export async function getAvailable(): Promise<Group[]> {
  const { data } = await apiClient.get<Group[]>('/groups/available')
  return data
}

/**
 * Get user-visible summaries for groups the current user can use
 * @returns List of lightweight group summaries
 */
export async function getAvailableSummary(): Promise<UserGroupSummary[]> {
  const { data } = await apiClient.get<UserGroupSummary[]>('/groups/available/summary')
  return data
}

/**
 * Get current user's custom group rate multipliers
 * @returns Map of group_id to custom rate_multiplier
 */
export async function getUserGroupRates(): Promise<Record<number, number>> {
  const { data } = await apiClient.get<Record<number, number> | null>('/groups/rates')
  return data || {}
}

export const userGroupsAPI = {
  getAvailable,
  getAvailableSummary,
  getUserGroupRates
}

export default userGroupsAPI
