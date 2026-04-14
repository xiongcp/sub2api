package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccountIsBelowQuotaSchedulingThreshold_Absolute(t *testing.T) {
	account := &Account{
		Status:      StatusActive,
		Schedulable: true,
		Extra: map[string]any{
			"quota_limit":         10.0,
			"quota_used":          9.6,
			"quota_min_remaining": 0.5,
		},
	}

	require.True(t, account.IsBelowQuotaSchedulingThreshold())
	require.False(t, account.IsSchedulable())
}

func TestAccountIsBelowQuotaSchedulingThreshold_Ratio(t *testing.T) {
	account := &Account{
		Status:      StatusActive,
		Schedulable: true,
		Extra: map[string]any{
			"quota_limit":               100.0,
			"quota_used":                96.0,
			"quota_min_remaining_ratio": 0.05,
		},
	}

	require.True(t, account.IsBelowQuotaSchedulingThreshold())
	require.False(t, account.IsSchedulable())
}

func TestAccountIsBelowQuotaSchedulingThreshold_NoThreshold(t *testing.T) {
	account := &Account{
		Status:      StatusActive,
		Schedulable: true,
		Extra: map[string]any{
			"quota_limit": 10.0,
			"quota_used":  9.6,
		},
	}

	require.False(t, account.IsBelowQuotaSchedulingThreshold())
	require.True(t, account.IsSchedulable())
}

func TestAccountIsBelowQuotaSchedulingThreshold_NoTotalQuotaLimit(t *testing.T) {
	account := &Account{
		Status:      StatusActive,
		Schedulable: true,
		Extra: map[string]any{
			"quota_daily_limit":         10.0,
			"quota_daily_used":          9.9,
			"quota_min_remaining":       0.5,
			"quota_min_remaining_ratio": 0.1,
		},
	}

	require.False(t, account.IsBelowQuotaSchedulingThreshold())
	require.True(t, account.IsSchedulable())
}
