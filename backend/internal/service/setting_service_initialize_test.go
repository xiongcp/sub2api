//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type settingInitializeRepoStub struct {
	values   map[string]string
	writes   map[string]string
	writeCnt int
}

func (s *settingInitializeRepoStub) Get(ctx context.Context, key string) (*Setting, error) {
	panic("unexpected Get call")
}

func (s *settingInitializeRepoStub) GetValue(ctx context.Context, key string) (string, error) {
	panic("unexpected GetValue call")
}

func (s *settingInitializeRepoStub) Set(ctx context.Context, key, value string) error {
	panic("unexpected Set call")
}

func (s *settingInitializeRepoStub) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	panic("unexpected GetMultiple call")
}

func (s *settingInitializeRepoStub) SetMultiple(ctx context.Context, settings map[string]string) error {
	s.writeCnt++
	s.writes = make(map[string]string, len(settings))
	for key, value := range settings {
		s.writes[key] = value
		s.values[key] = value
	}
	return nil
}

func (s *settingInitializeRepoStub) GetAll(ctx context.Context) (map[string]string, error) {
	out := make(map[string]string, len(s.values))
	for key, value := range s.values {
		out[key] = value
	}
	return out, nil
}

func (s *settingInitializeRepoStub) Delete(ctx context.Context, key string) error {
	panic("unexpected Delete call")
}

func TestSettingService_InitializeDefaultSettings_FillsMissingBrandingDefaults(t *testing.T) {
	repo := &settingInitializeRepoStub{
		values: map[string]string{
			SettingKeyRegistrationEnabled: "true",
			SettingKeySiteName:            "Existing Site",
		},
	}
	svc := NewSettingService(repo, &config.Config{
		Default: config.DefaultConfig{
			UserConcurrency: 3,
			UserBalance:     1.5,
		},
	})

	err := svc.InitializeDefaultSettings(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, repo.writeCnt)
	require.Equal(t, defaultBrandingHomeContent, repo.writes[SettingKeyHomeContent])
	require.Equal(t, defaultBrandingCustomCSS, repo.writes[SettingKeyCustomCSS])
	require.Equal(t, defaultBrandingLoginExtraHTML, repo.writes[SettingKeyLoginExtraHTML])
	require.Equal(t, defaultBrandingRegisterExtraHTML, repo.writes[SettingKeyRegisterExtraHTML])
	require.Equal(t, "", repo.writes[SettingKeyPaymentFooterHTML])
	require.Equal(t, defaultBrandingGlobalFooterHTML, repo.writes[SettingKeyGlobalFooterHTML])
	require.NotContains(t, repo.writes, SettingKeyRegistrationEnabled)
	require.Equal(t, "Existing Site", repo.values[SettingKeySiteName])
}

func TestSettingService_InitializeDefaultSettings_PreservesExplicitEmptyBrandingValues(t *testing.T) {
	repo := &settingInitializeRepoStub{
		values: map[string]string{
			SettingKeyHomeContent:       "",
			SettingKeyCustomCSS:         "",
			SettingKeyLoginExtraHTML:    "",
			SettingKeyRegisterExtraHTML: "",
			SettingKeyPaymentFooterHTML: "",
			SettingKeyGlobalFooterHTML:  "",
		},
	}
	svc := NewSettingService(repo, &config.Config{
		Default: config.DefaultConfig{
			UserConcurrency: 2,
			UserBalance:     0.5,
		},
	})

	err := svc.InitializeDefaultSettings(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, repo.writeCnt)
	require.NotContains(t, repo.writes, SettingKeyHomeContent)
	require.NotContains(t, repo.writes, SettingKeyCustomCSS)
	require.NotContains(t, repo.writes, SettingKeyLoginExtraHTML)
	require.NotContains(t, repo.writes, SettingKeyRegisterExtraHTML)
	require.NotContains(t, repo.writes, SettingKeyPaymentFooterHTML)
	require.NotContains(t, repo.writes, SettingKeyGlobalFooterHTML)
}

func TestSettingService_InitializeDefaultSettings_SetsSMTPSecurityModeDefault(t *testing.T) {
	repo := &settingInitializeRepoStub{
		values: map[string]string{},
	}
	svc := NewSettingService(repo, &config.Config{})

	err := svc.InitializeDefaultSettings(context.Background())
	require.NoError(t, err)
	require.Equal(t, string(SMTPSecurityModeStartTLS), repo.writes[SettingKeySMTPSecurityMode])
	require.Equal(t, "false", repo.writes[SettingKeySMTPUseTLS])
}
