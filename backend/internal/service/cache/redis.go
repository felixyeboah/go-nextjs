package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nanayaw/fullstack/internal/config"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService(cfg *config.RedisConfig) (*RedisService, error) {
	opt, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisService{client: client}, nil
}

// Rate limiting
func (s *RedisService) CheckRateLimit(ctx context.Context, key string, limit int, duration int) (bool, error) {
	pipe := s.client.Pipeline()

	// Increment counter
	incr := pipe.Incr(ctx, key)
	// Set expiration if key is new
	pipe.Expire(ctx, key, time.Duration(duration)*time.Second)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check rate limit: %w", err)
	}

	count := incr.Val()
	return count <= int64(limit), nil
}

func (s *RedisService) ResetRateLimit(ctx context.Context, key string) error {
	if err := s.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to reset rate limit: %w", err)
	}
	return nil
}

// Caching
func (s *RedisService) CacheData(ctx context.Context, key string, data interface{}, ttl int) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := s.client.Set(ctx, key, jsonData, time.Duration(ttl)*time.Second).Err(); err != nil {
		return fmt.Errorf("failed to cache data: %w", err)
	}

	return nil
}

func (s *RedisService) GetCachedData(ctx context.Context, key string, dest interface{}) error {
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil // Key doesn't exist
		}
		return fmt.Errorf("failed to get cached data: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal cached data: %w", err)
	}

	return nil
}

func (s *RedisService) InvalidateCache(ctx context.Context, key string) error {
	if err := s.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to invalidate cache: %w", err)
	}
	return nil
}

func (s *RedisService) InvalidateCachePattern(ctx context.Context, pattern string) error {
	iter := s.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := s.client.Del(ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("failed to invalidate cache pattern: %w", err)
		}
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan cache keys: %w", err)
	}
	return nil
}

// Session management
func (s *RedisService) StoreSession(ctx context.Context, sessionID string, userID string, expiration time.Duration) error {
	if err := s.client.Set(ctx, fmt.Sprintf("session:%s", sessionID), userID, expiration).Err(); err != nil {
		return fmt.Errorf("failed to store session: %w", err)
	}
	return nil
}

func (s *RedisService) GetSession(ctx context.Context, sessionID string) (string, error) {
	userID, err := s.client.Get(ctx, fmt.Sprintf("session:%s", sessionID)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("failed to get session: %w", err)
	}
	return userID, nil
}

func (s *RedisService) InvalidateSession(ctx context.Context, sessionID string) error {
	if err := s.client.Del(ctx, fmt.Sprintf("session:%s", sessionID)).Err(); err != nil {
		return fmt.Errorf("failed to invalidate session: %w", err)
	}
	return nil
}

// Cache operations
func (s *RedisService) FlushAll(ctx context.Context) error {
	if err := s.client.FlushAll(ctx).Err(); err != nil {
		return fmt.Errorf("failed to flush cache: %w", err)
	}
	return nil
}

func (s *RedisService) Ping(ctx context.Context) error {
	if err := s.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}
	return nil
}
