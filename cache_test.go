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

func setRedisv8Cache() {
	cache.CacheInstance = cache.RedisV8Cache(func() *redisv8.Client {
		return redisv8.NewClient(&redisv8.Options{
			Addr:     "10.0.11.125:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	})
}
func setRedisv9Cache() {
	cache.CacheInstance = cache.RedisV9Cache(func() *redisv9.Client {
		return redisv9.NewClient(&redisv9.Options{
			Addr:     "10.0.11.125:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	})
}

func testCaseStruct() (users []user, err error) {
	key := "usercache"
	users = make([]user, 0)
	err = cache.Remember(key, 20*time.Second, &users, func() (any, error) {
		fmt.Println("load from db")
		data := []user{{ID: 1, Name: "suifengpiao14"}, {ID: 2, Name: "suifengpiao15"}}
		return data, nil
	})
	return users, err
}

func testCaseInt64() (count int64, err error) {
	key := "usercache"
	err = cache.Remember(key, 20*time.Second, &count, func() (any, error) {
		fmt.Println("load from db")
		return 52, nil
	})
	return count, err
}
func testCaseInt() (count int, err error) {
	key := "usercache"
	err = cache.Remember(key, 20*time.Second, &count, func() (any, error) {
		fmt.Println("load from db")
		return 52, nil
	})
	return count, err
}
func testCaseBool() (exists bool, err error) {
	key := "usercache"
	err = cache.Remember(key, 20*time.Second, &exists, func() (any, error) {
		fmt.Println("load from db")
		return true, nil
	})
	return exists, err
}

func TestStruct(t *testing.T) {
	t.Run("redisv8", func(t *testing.T) {
		setRedisv8Cache()
		users, err := testCaseStruct()
		require.NoError(t, err)
		fmt.Println(users)
	})

	t.Run("redisv9", func(t *testing.T) {
		setRedisv9Cache()
		users, err := testCaseStruct()
		require.NoError(t, err)
		fmt.Println(users)
	})
	t.Run("memery", func(t *testing.T) {
		users, err := testCaseStruct()
		require.NoError(t, err)
		fmt.Println(users)
	})
}

func TestInt64(t *testing.T) {
	t.Run("redisv8", func(t *testing.T) {
		setRedisv8Cache()
		count, err := testCaseInt64()
		require.NoError(t, err)
		fmt.Println(count)
	})

	t.Run("redisv9", func(t *testing.T) {
		setRedisv9Cache()
		count, err := testCaseInt64()
		require.NoError(t, err)
		fmt.Println(count)
	})
	t.Run("memery", func(t *testing.T) {
		count, err := testCaseInt64()
		require.NoError(t, err)
		fmt.Println(count)
	})

}
func TestInt(t *testing.T) {
	t.Run("redisv8", func(t *testing.T) {
		setRedisv8Cache()
		count, err := testCaseInt()
		require.NoError(t, err)
		fmt.Println(count)
	})

	t.Run("redisv9", func(t *testing.T) {
		setRedisv9Cache()
		count, err := testCaseInt()
		require.NoError(t, err)
		fmt.Println(count)
	})
	t.Run("memery", func(t *testing.T) {
		count, err := testCaseInt()
		require.NoError(t, err)
		fmt.Println(count)
	})

}

func TestBool(t *testing.T) {
	t.Run("redisv8", func(t *testing.T) {
		setRedisv8Cache()
		count, err := testCaseBool()
		require.NoError(t, err)
		fmt.Println(count)
	})

	t.Run("redisv9", func(t *testing.T) {
		setRedisv9Cache()
		count, err := testCaseBool()
		require.NoError(t, err)
		fmt.Println(count)
	})
	t.Run("memery", func(t *testing.T) {
		count, err := testCaseBool()
		require.NoError(t, err)
		fmt.Println(count)
	})

}
