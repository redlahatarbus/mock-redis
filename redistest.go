package mockredis

import (
	"context"
	"regexp"

	"github.com/go-redis/redis/v8"
)

type RedisTestHook struct {
	Target string                    // possible values are "before" or "after", based on where we would want to run the hook action
	C      *redis.Client             // the redis client
	Match  string                    // the regular expression against which we would check our command
	Action func(*redis.Client) error // the function to run if regular expression matches command
}

func (h *RedisTestHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if h.Target == "before" {
		if match, err := regexp.MatchString(h.Match, cmd.String()); err == nil && match {
			return nil, h.Action(h.C)
		}
	}
	return ctx, nil
}

func (h *RedisTestHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error { return nil }
func (h *RedisTestHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}
func (h *RedisTestHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}
