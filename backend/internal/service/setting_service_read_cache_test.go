//go:build unit

package service

import (
	"context"
	"encoding/json"
	"slices"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type settingReadRepoStub struct {
	values           map[string]string
	getValueCalls    int
	getMultipleCalls int
	setMultipleCalls int
}

func (s *settingReadRepoStub) Get(ctx context.Context, key string) (*Setting, error) {
	panic("unexpected Get call")
}

func (s *settingReadRepoStub) GetValue(ctx context.Context, key string) (string, error) {
	s.getValueCalls++
	return s.values[key], nil
}

func (s *settingReadRepoStub) Set(ctx context.Context, key, value string) error {
	panic("unexpected Set call")
}

func (s *settingReadRepoStub) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	s.getMultipleCalls++
	out := make(map[string]string, len(keys))
	for _, key := range keys {
		if value, ok := s.values[key]; ok {
			out[key] = value
		}
	}
	return out, nil
}

func (s *settingReadRepoStub) SetMultiple(ctx context.Context, settings map[string]string) error {
	s.setMultipleCalls++
	if s.values == nil {
		s.values = make(map[string]string, len(settings))
	}
	for k, v := range settings {
		s.values[k] = v
	}
	return nil
}

func (s *settingReadRepoStub) GetAll(ctx context.Context) (map[string]string, error) {
	panic("unexpected GetAll call")
}

func (s *settingReadRepoStub) Delete(ctx context.Context, key string) error {
	panic("unexpected Delete call")
}

type settingReadCacheStub struct {
	store     map[string][]byte
	deleted   []string
	published []string
	handler   func(string)
}

func newSettingReadCacheStub() *settingReadCacheStub {
	return &settingReadCacheStub{
		store: make(map[string][]byte),
	}
}

func (s *settingReadCacheStub) Get(ctx context.Context, key string, dest any) (bool, error) {
	data, ok := s.store[key]
	if !ok {
		return false, nil
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return false, err
	}
	return true, nil
}

func (s *settingReadCacheStub) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	s.store[key] = data
	return nil
}

func (s *settingReadCacheStub) Delete(ctx context.Context, key string) error {
	delete(s.store, key)
	s.deleted = append(s.deleted, key)
	return nil
}

func (s *settingReadCacheStub) PublishInvalidation(ctx context.Context, key string) error {
	s.published = append(s.published, key)
	return nil
}

func (s *settingReadCacheStub) SubscribeInvalidation(ctx context.Context, handler func(key string)) error {
	s.handler = handler
	return nil
}

func TestSettingService_GetPublicSettings_UsesReadCacheAfterFirstLoad(t *testing.T) {
	repo := &settingReadRepoStub{
		values: map[string]string{
			SettingKeySiteName: "Cached Site",
		},
	}
	cache := newSettingReadCacheStub()
	svc := NewSettingService(repo, &config.Config{})
	svc.SetReadCache(cache)

	first, err := svc.GetPublicSettings(context.Background())
	require.NoError(t, err)
	second, err := svc.GetPublicSettings(context.Background())
	require.NoError(t, err)

	require.Equal(t, "Cached Site", first.SiteName)
	require.Equal(t, "Cached Site", second.SiteName)
	require.Equal(t, 1, repo.getMultipleCalls)
	_, ok := cache.store[settingReadCacheKeyPublic]
	require.True(t, ok)
}

func TestSettingService_GetPublicSettings_ReadsRedisCacheBeforeRepo(t *testing.T) {
	repo := &settingReadRepoStub{
		values: map[string]string{
			SettingKeySiteName: "DB Site",
		},
	}
	cache := newSettingReadCacheStub()
	require.NoError(t, cache.Set(context.Background(), settingReadCacheKeyPublic, &PublicSettings{
		SiteName: "Redis Site",
	}, time.Minute))

	svc := NewSettingService(repo, &config.Config{})
	svc.SetReadCache(cache)

	settings, err := svc.GetPublicSettings(context.Background())
	require.NoError(t, err)
	require.Equal(t, "Redis Site", settings.SiteName)
	require.Equal(t, 0, repo.getMultipleCalls)
}

func TestSettingService_GetFrontendURL_ReadsRedisCacheBeforeRepo(t *testing.T) {
	repo := &settingReadRepoStub{
		values: map[string]string{
			SettingKeyFrontendURL: "https://db.example.com",
		},
	}
	cache := newSettingReadCacheStub()
	require.NoError(t, cache.Set(context.Background(), settingReadCacheKeyFrontend, "https://redis.example.com", time.Minute))

	svc := NewSettingService(repo, &config.Config{})
	svc.SetReadCache(cache)

	require.Equal(t, "https://redis.example.com", svc.GetFrontendURL(context.Background()))
	require.Equal(t, 0, repo.getValueCalls)
}

func TestSettingService_UpdateSettings_InvalidatesReadCaches(t *testing.T) {
	repo := &settingReadRepoStub{
		values: map[string]string{
			SettingKeySiteName:                  "Old Site",
			SettingKeyAPIBaseURL:                "https://old.example.com/v1",
			SettingKeyAPIKeyUsageGuideContent:   `{"description":"old"}`,
			SettingKeyFrontendURL:               "https://old-frontend.example.com",
			SettingKeyRegistrationEnabled:       "true",
			SettingKeyEmailVerifyEnabled:        "true",
			SettingKeyPasswordResetEnabled:      "true",
			SettingKeyTableDefaultPageSize:      "20",
			SettingKeyTablePageSizeOptions:      "[20,50,100]",
			SettingKeyBalanceLowNotifyThreshold: "0",
		},
	}
	cache := newSettingReadCacheStub()
	svc := NewSettingService(repo, &config.Config{})
	svc.SetReadCache(cache)

	_, err := svc.GetPublicSettings(context.Background())
	require.NoError(t, err)
	_, err = svc.GetAPIKeyUsageGuide(context.Background())
	require.NoError(t, err)
	require.Equal(t, "https://old-frontend.example.com", svc.GetFrontendURL(context.Background()))

	err = svc.UpdateSettings(context.Background(), &SystemSettings{
		SiteName:   "New Site",
		APIBaseURL: "https://new.example.com/v1",
		APIKeyUsageGuideContent: APIKeyUsageGuideContent{
			Description: "new",
		},
		FrontendURL: "https://new-frontend.example.com",
	})
	require.NoError(t, err)

	require.True(t, slices.Contains(cache.deleted, settingReadCacheKeyPublic))
	require.True(t, slices.Contains(cache.deleted, settingReadCacheKeyGuide))
	require.True(t, slices.Contains(cache.deleted, settingReadCacheKeyFrontend))
	require.True(t, slices.Contains(cache.published, settingReadCacheKeyPublic))
	require.True(t, slices.Contains(cache.published, settingReadCacheKeyGuide))
	require.True(t, slices.Contains(cache.published, settingReadCacheKeyFrontend))

	settings, err := svc.GetPublicSettings(context.Background())
	require.NoError(t, err)
	guide, err := svc.GetAPIKeyUsageGuide(context.Background())
	require.NoError(t, err)
	frontendURL := svc.GetFrontendURL(context.Background())

	require.Equal(t, "New Site", settings.SiteName)
	require.Equal(t, "https://new.example.com/v1", guide.APIBaseURL)
	require.Equal(t, "new", guide.Content.Description)
	require.Equal(t, "https://new-frontend.example.com", frontendURL)
	require.Equal(t, 4, repo.getMultipleCalls)
	require.Equal(t, 2, repo.getValueCalls)
	require.Equal(t, 1, repo.setMultipleCalls)
}
