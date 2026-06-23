package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/redis/go-redis/v9"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/redis/go-redis.v9"
)

type Client struct {
	inner *redis.Client
}

func NewClient(host, port, password string, db int) (*Client, error) {
	addr := fmt.Sprintf("%s:%s", host, port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeRedisError, "failed to connect to redis: %v", err))
	}

	return &Client{inner: rdb}, nil
}

func (c *Client) Close() error {
	return c.inner.Close()
}

func (c *Client) EnableDDTrace(opts ...redistrace.ClientOption) {
	redistrace.WrapClient(c.inner, opts...)
}

func (c *Client) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return c.inner.Set(ctx, key, value, expiration).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.inner.Get(ctx, key).Result()
}

func (c *Client) SetJSON(ctx context.Context, key string, value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeRedisError, "failed to marshal value: %v", err))
	}
	return c.inner.Set(ctx, key, data, expiration).Err()
}

func (c *Client) GetJSON(ctx context.Context, key string, dest any) error {
	data, err := c.inner.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (c *Client) Delete(ctx context.Context, keys ...string) error {
	return c.inner.Del(ctx, keys...).Err()
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.inner.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (c *Client) SetNX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	return c.inner.SetNX(ctx, key, value, expiration).Result()
}

// MGetBytes fetches multiple keys in a single round-trip.
// Nil entries mean the key was not found in Redis.
func (c *Client) MGetBytes(ctx context.Context, keys ...string) ([][]byte, error) {
	results, err := c.inner.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	out := make([][]byte, len(results))
	for i, r := range results {
		if r == nil {
			continue
		}
		if s, ok := r.(string); ok {
			out[i] = []byte(s)
		}
	}
	return out, nil
}

func (c *Client) SetNXJSON(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeRedisError, "failed to marshal value: %v", err))
	}
	return c.inner.SetNX(ctx, key, data, expiration).Result()
}

func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.inner.Expire(ctx, key, expiration).Err()
}

func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.inner.TTL(ctx, key).Result()
}

func (c *Client) HSet(ctx context.Context, key string, field string, value string) error {
	return c.inner.HSet(ctx, key, field, value).Err()
}

func (c *Client) HGet(ctx context.Context, key string, field string) (string, error) {
	return c.inner.HGet(ctx, key, field).Result()
}

func (c *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.inner.HGetAll(ctx, key).Result()
}

func (c *Client) HSetJSON(ctx context.Context, key string, field string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return apperror.WrapWithMessage(err, apperror.NewMessage(apperror.CodeRedisError, "failed to marshal value: %v", err))
	}
	return c.inner.HSet(ctx, key, field, data).Err()
}

func (c *Client) HGetJSON(ctx context.Context, key string, field string, dest any) error {
	data, err := c.inner.HGet(ctx, key, field).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (c *Client) Keys(ctx context.Context, pattern string) ([]string, error) {
	return c.inner.Keys(ctx, pattern).Result()
}

func (c *Client) Scan(ctx context.Context, pattern string) ([]string, error) {
	var keys []string
	var cursor uint64
	for {
		var batch []string
		var err error
		batch, cursor, err = c.inner.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, batch...)
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}

func (c *Client) IncrBy(ctx context.Context, key string, value int64) error {
	return c.inner.IncrBy(ctx, key, value).Err()
}

// GetInt64 reads a key that holds an integer value without modifying it.
// Returns 0 if the key does not exist.
func (c *Client) GetInt64(ctx context.Context, key string) (int64, error) {
	result, err := c.inner.Get(ctx, key).Int64()
	if err != nil {
		if IsNotFound(err) {
			return 0, nil
		}
		return 0, err
	}
	return result, nil
}

// GetDel atomically reads and deletes a key that holds an integer value.
// Returns 0 if the key does not exist.
func (c *Client) GetDel(ctx context.Context, key string) (int64, error) {
	result, err := c.inner.GetDel(ctx, key).Int64()
	if err != nil {
		if IsNotFound(err) {
			return 0, nil
		}
		return 0, err
	}
	return result, nil
}

func IsNotFound(err error) bool {
	return err == redis.Nil
}
