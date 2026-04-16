//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func TestSettingService_GetAPIKeyUsageGuide(t *testing.T) {
	repo := &settingPublicRepoStub{
		values: map[string]string{
			SettingKeyAPIBaseURL: "https://api.example.com/v1",
			SettingKeyAPIKeyUsageGuideContent: `{
				"description":"dynamic description",
				"openai":{"config_toml_hint":"top of file"},
				"opencode":{"hint":"dynamic hint"}
			}`,
		},
	}
	svc := NewSettingService(repo, &config.Config{})

	guide, err := svc.GetAPIKeyUsageGuide(context.Background())
	require.NoError(t, err)
	require.Equal(t, "https://api.example.com/v1", guide.APIBaseURL)
	require.Equal(t, "dynamic description", guide.Content.Description)
	require.Equal(t, "top of file", guide.Content.OpenAI.ConfigTomlHint)
	require.Equal(t, "dynamic hint", guide.Content.OpenCode.Hint)
}

func TestParseAPIKeyUsageGuideContent_InvalidJSONReturnsEmpty(t *testing.T) {
	content := ParseAPIKeyUsageGuideContent("{not-json")
	require.Equal(t, APIKeyUsageGuideContent{}, content)
}
