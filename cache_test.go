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
	key := "testCaseStruct"
	users = make([]user, 0) //
	err = cache.Remember(key, &users, func(dst *[]user) (duration time.Duration, err error) {
		fmt.Println("load from db")
		data := []user{{ID: 1, Name: "suifengpiao14"}, {ID: 2, Name: "suifengpiao15"}}
		*dst = data
		return 20 * time.Second, nil
	})
	return users, err
}

func testCaseInt64() (count int64, err error) {
	key := "testCaseInt64"
	err = cache.Remember(key, &count, func(dst *int64) (duration time.Duration, err error) {
		fmt.Println("load from db")
		*dst = 52
		return 20 * time.Second, nil
	})
	return count, err
}
func testCaseInt() (count int, err error) {
	key := "testCaseInt"
	err = cache.Remember(key, &count, func(dst *int) (duration time.Duration, err error) {
		fmt.Println("load from db")
		*dst = 52
		return 20 * time.Second, nil
	})
	return count, err
}
func testCaseBool() (exists bool, err error) {
	key := "testCaseBool"
	err = cache.Remember(key, &exists, func(dst *bool) (duration time.Duration, err error) {
		fmt.Println("load from db")
		*dst = true
		return 20 * time.Second, nil
	})
	return exists, err
}

func TestStruct(t *testing.T) {
	t.Run("redisv8", func(t *testing.T) {
		setRedisv8Cache()
		users1, err := testCaseStruct() // 第一次加载数据
		require.NoError(t, err)
		users2, err := testCaseStruct() //第二次从缓存中获取
		require.NoError(t, err)
		require.Equal(t, users1, users2)
	})

	t.Run("redisv9", func(t *testing.T) {
		setRedisv9Cache()
		users1, err := testCaseStruct() // 第一次加载数据
		require.NoError(t, err)
		users2, err := testCaseStruct() //第二次从缓存中获取
		require.NoError(t, err)
		require.Equal(t, users1, users2)
	})
	t.Run("memery", func(t *testing.T) {
		users1, err := testCaseStruct() // 第一次加载数据
		require.NoError(t, err)
		users2, err := testCaseStruct() //第二次从缓存中获取
		require.NoError(t, err)
		require.Equal(t, users1, users2)
	})
}

func TestInt64(t *testing.T) {
	t.Run("redisv8", func(t *testing.T) {
		setRedisv8Cache()
		count1, err := testCaseInt64() // 第一次加载数据
		require.NoError(t, err)
		count2, err := testCaseInt64() //第二次从缓存中获取
		require.NoError(t, err)
		require.Equal(t, count1, count2)
	})

	t.Run("redisv9", func(t *testing.T) {
		setRedisv9Cache()
		count1, err := testCaseInt64() // 第一次加载数据
		require.NoError(t, err)
		count2, err := testCaseInt64() //第二次从缓存中获取
		require.NoError(t, err)
		require.Equal(t, count1, count2)
	})
	t.Run("memery", func(t *testing.T) {
		count1, err := testCaseInt64() // 第一次加载数据
		require.NoError(t, err)
		count2, err := testCaseInt64() //第二次从缓存中获取
		require.NoError(t, err)
		require.Equal(t, count1, count2)
	})

}
func TestInt(t *testing.T) {
	t.Run("redisv8", func(t *testing.T) {
		setRedisv8Cache()
		count1, err := testCaseInt() // 第一次加载数据
		require.NoError(t, err)
		count2, err := testCaseInt() //第二次从缓存中获取
		require.NoError(t, err)
		require.Equal(t, count1, count2)
	})

	t.Run("redisv9", func(t *testing.T) {
		setRedisv9Cache()
		count1, err := testCaseInt() // 第一次加载数据
		require.NoError(t, err)
		count2, err := testCaseInt() //第二次从缓存中获取
		require.NoError(t, err)
		require.Equal(t, count1, count2)
	})
	t.Run("memery", func(t *testing.T) {
		count1, err := testCaseInt() // 第一次加载数据
		require.NoError(t, err)
		count2, err := testCaseInt() //第二次从缓存中获取
		require.NoError(t, err)
		require.Equal(t, count1, count2)
	})

}

func TestBool(t *testing.T) {
	t.Run("redisv8", func(t *testing.T) {
		setRedisv8Cache()
		count1, err := testCaseBool() // 第一次加载数据
		require.NoError(t, err)
		count2, err := testCaseBool() //第二次从缓存中获取
		require.NoError(t, err)
		require.Equal(t, count1, count2)
	})

	t.Run("redisv9", func(t *testing.T) {
		setRedisv9Cache()
		count1, err := testCaseBool() // 第一次加载数据
		require.NoError(t, err)
		count2, err := testCaseBool() //第二次从缓存中获取
		require.NoError(t, err)
		require.Equal(t, count1, count2)
	})
	t.Run("memery", func(t *testing.T) {
		count1, err := testCaseBool() // 第一次加载数据
		require.NoError(t, err)
		count2, err := testCaseBool() //第二次从缓存中获取
		require.NoError(t, err)
		require.Equal(t, count1, count2)
	})

}
