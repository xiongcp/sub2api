package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/redis/go-redis/v9"
)

const (
	settingReadCacheKeyPrefix      = "setting:read:"
	settingReadCacheInvalidationCh = "setting:read:invalidate"
)

type settingReadCache struct {
	rdb *redis.Client
}

func NewSettingReadCache(rdb *redis.Client) service.SettingReadCache {
	return &settingReadCache{rdb: rdb}
}

func (c *settingReadCache) Get(ctx context.Context, key string, dest any) (bool, error) {
	data, err := c.rdb.Get(ctx, c.redisKey(key)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return false, err
	}
	return true, nil
}

func (c *settingReadCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, c.redisKey(key), data, ttl).Err()
}

func (c *settingReadCache) Delete(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, c.redisKey(key)).Err()
}

func (c *settingReadCache) PublishInvalidation(ctx context.Context, key string) error {
	return c.rdb.Publish(ctx, settingReadCacheInvalidationCh, key).Err()
}

func (c *settingReadCache) SubscribeInvalidation(ctx context.Context, handler func(key string)) error {
	pubsub := c.rdb.Subscribe(ctx, settingReadCacheInvalidationCh)
	if _, err := pubsub.Receive(ctx); err != nil {
		_ = pubsub.Close()
		return fmt.Errorf("subscribe setting read cache invalidation: %w", err)
	}

	go func() {
		defer func() {
			if err := pubsub.Close(); err != nil {
				log.Printf("Warning: failed to close setting read cache pubsub: %v", err)
			}
		}()

		ch := pubsub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				if msg != nil {
					handler(msg.Payload)
				}
			}
		}
	}()

	return nil
}

func (c *settingReadCache) redisKey(key string) string {
	return settingReadCacheKeyPrefix + key
}
