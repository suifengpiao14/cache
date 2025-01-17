package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"sync"
	"time"

	reidsv8 "github.com/go-redis/redis/v8"
	reidsv9 "github.com/redis/go-redis/v9"

	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

type Cache interface {
	Get(key string, data any) (exists bool, err error)
	Set(key string, data any, duration time.Duration) error
}

var RedisV8Getter = OnceValue(func() *reidsv8.Client {
	err := errors.New("RedisV8Instance 未初始化")
	panic(err)
})

// 使用reids 作为缓存中间件，必须先初始化RedisV8Getter或者RedisV9Getter
var RedisV9Getter = OnceValue(func() *reidsv9.Client {
	err := errors.New("RedisV9Instance 未初始化")
	panic(err)
})

func Remember(cache Cache, key string, duration time.Duration, dst any, fetchFunc func() (any, error)) error {
	md5Key := Md5Lower(key)
	exists, err := cache.Get(md5Key, dst)
	if err != nil {
		return err
	}
	if exists { // 正常取到直接返回
		return nil
	}
	data, err := fetchFunc()
	if err != nil {
		return err
	}
	err = cache.Set(md5Key, data, duration)
	if err != nil {
		return err
	}
	reflect.Indirect(reflect.ValueOf(dst)).Set(reflect.Indirect(reflect.ValueOf(data)))
	return nil
}

func RedisV8Cache(client func() *reidsv8.Client) Cache {
	return _RedisV8Cache{}
}

type _RedisV8Cache struct {
}

func (r _RedisV8Cache) Get(key string, data any) (exists bool, err error) {
	ctx := context.Background()
	b, err := RedisV8Getter().Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, reidsv8.Nil) { // 是redis.Nil 错误，屏蔽错误，exists 返回false
			return false, nil
		}
		return false, err
	}
	err = json.Unmarshal(b, data)
	if err != nil {
		return false, err
	}
	return true, nil

}

func (r _RedisV8Cache) Set(key string, data any, duration time.Duration) (err error) {
	ctx := context.Background()
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = RedisV8Getter().Set(ctx, key, b, duration).Result()
	if err != nil {
		return err
	}
	return nil
}

func RedisV9Cache(client func() *reidsv9.Client) Cache {
	return _RedisV9Cache{}
}

type _RedisV9Cache struct {
}

func (r _RedisV9Cache) Get(key string, data any) (exists bool, err error) {
	ctx := context.Background()
	b, err := RedisV9Getter().Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, reidsv8.Nil) { // 是redis.Nil 错误，屏蔽错误，exists 返回false
			return false, nil
		}
		return false, err
	}
	err = json.Unmarshal(b, data)
	if err != nil {
		return false, err
	}
	return true, nil

}

func (r _RedisV9Cache) Set(key string, data any, duration time.Duration) (err error) {
	ctx := context.Background()
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = RedisV9Getter().Set(ctx, key, b, duration).Result()
	if err != nil {
		return err
	}
	return nil
}

var memeryCache = cache.New(1*time.Second, 10*time.Minute) //默认 1秒缓存的内存缓存实例 常用于单次请求,某个接口、sql 结果
type _MemeryCache struct {
}

func MemeryCache() Cache {
	return _MemeryCache{}
}

func (m _MemeryCache) Get(key string, dst any) (exists bool, err error) {
	resulany, found := memeryCache.Get(key)
	if !found {
		return false, nil
	}
	reflect.Indirect(reflect.ValueOf(dst)).Set(reflect.Indirect(reflect.ValueOf(resulany)))
	return true, nil
}

func (m _MemeryCache) Set(key string, data any, duration time.Duration) error {
	memeryCache.Set(key, data, duration)
	return nil
}

// 拷贝 sync.oncefunc.go 低版本go 不支持 go 1.21 版本才有，直接复制
// OnceValue returns a function that invokes f only once and returns the value
// returned by f. The returned function may be called concurrently.
//
// If f panics, the returned function will panic with the same value on every call.
func OnceValue[T any](f func() T) func() T {
	var (
		once   sync.Once
		valid  bool
		p      any
		result T
	)
	g := func() {
		defer func() {
			p = recover()
			if !valid {
				panic(p)
			}
		}()
		result = f()
		valid = true
	}
	return func() T {
		once.Do(g)
		if !valid {
			panic(p)
		}
		return result
	}
}

// Md5Lower md5 小写
func Md5Lower(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
