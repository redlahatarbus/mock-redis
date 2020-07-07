package mockredis

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
)

func DoSomething(ctx context.Context, c *redis.Client, hashMap map[string]interface{}, expiry time.Duration) error {
	key := "myhash"

	r1 := c.HSet(ctx, key, hashMap)
	if err := r1.Err(); err != nil {
		return err
	}

	r2 := c.Expire(ctx, key, expiry)
	if err := r2.Err(); err != nil {
		return err
	}

	return nil
}
