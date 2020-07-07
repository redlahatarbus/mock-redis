package mockredis

import (
	"context"
	"fmt"
	_ "strconv"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	redis "github.com/go-redis/redis/v8"
)

func TestDoSomething(t *testing.T) {
	mr, _ := miniredis.Run()
	defer mr.Close()
	c := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	err := DoSomething(context.Background(), c, map[string]interface{}{"key1": "value1"}, time.Hour)

	// Validate error
	if err != nil {
		t.Errorf("No error was expected")
	}

	// Validate hash key value
	if val := mr.HGet("myhash", "key1"); val != "value1" {
		t.Errorf("Actual value %v does not match expected value", val)
	}

	// Validate expiry
	if mr.FastForward(time.Hour + 1); mr.Exists("myhash") {
		t.Errorf("Key should have expired")
	}
}

func TestDoSomethingHSetFail(t *testing.T) {
	mr, _ := miniredis.Run()
	defer mr.Close()
	c := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	c.AddHook(&RedisTestHook{
		Target: "before",
		C:      c,
		Match:  `^(?i)expire `,
		Action: func(c *redis.Client) error {
			return fmt.Errorf("expire error")
		},
	})

	err := DoSomething(context.Background(), c, map[string]interface{}{"key1": "value1"}, time.Hour)

	// Validate error
	if err == nil || !strings.Contains(err.Error(), "expire error") {
		t.Errorf("Expected error `%v` does not match actual error `%v`", "expire error", err)
	}
}
