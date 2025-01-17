package cache_test

import (
	"fmt"
	"testing"
	"time"

	redisv8 "github.com/go-redis/redis/v8"
	redisv9 "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/cache"
)

type user struct {
	ID   int
	Name string
}

func testCase() (users []user, err error) {
	key := "usercache"
	users = make([]user, 0)
	err = cache.Remember(key, 20*time.Second, &users, func() (any, error) {
		fmt.Println("load from db")
		data := []user{{ID: 1, Name: "suifengpiao14"}, {ID: 2, Name: "suifengpiao15"}}
		return data, nil
	})
	return users, err
}

func TestMemeryCache(t *testing.T) {
	users, err := testCase()
	require.NoError(t, err)
	fmt.Println(users)
}

func TestRedisv9Cache(t *testing.T) {
	cache.CacheInstance = cache.RedisV9Cache(func() *redisv9.Client {
		return redisv9.NewClient(&redisv9.Options{
			Addr:     "10.0.11.125:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	})
	users, err := testCase()
	require.NoError(t, err)
	fmt.Println(users)

}

func TestRedisv8Cache(t *testing.T) {
	cache.CacheInstance = cache.RedisV8Cache(func() *redisv8.Client {
		return redisv8.NewClient(&redisv8.Options{
			Addr:     "10.0.11.125:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	})
	users, err := testCase()
	require.NoError(t, err)
	fmt.Println(users)

}
