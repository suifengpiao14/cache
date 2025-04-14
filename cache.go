package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
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

//缓存实例，默认使用内存缓存

var CacheInstance Cache = MemeryCache()

func Remember[T any](key string, dst *T, fetchFunc func(dst *T) (duration time.Duration, err error)) error {
	return RememberWithCacheInstance(CacheInstance, key, dst, fetchFunc)
}

func RememberInMemory[T any](key string, dst *T, fetchFunc func(data *T) (duration time.Duration, err error)) error {
	return RememberWithCacheInstance(MemeryCache(), key, dst, fetchFunc)
}

// RememberWithCacheInstance 传入缓存实例，key,目标对象和获取数据的函数，将过期时间转移到回调函数返回值中，方便过期时间由回调函数控制(如获取微信access_token时，过期时间由微信官方返回)
func RememberWithCacheInstance[T any](cache Cache, key string, dst *T, fetchFunc func(dst *T) (duration time.Duration, err error)) error {
	cacheKey := strings.ReplaceAll(key, " ", "_")
	if len(cacheKey) > 32 {
		prefix, suffix := cacheKey[:32], cacheKey[32:]
		cacheKey = fmt.Sprintf("%s_%s", prefix, Md5Lower(suffix))
	}
	exists, err := cache.Get(cacheKey, dst)
	if err != nil {
		return err
	}
	if exists { // 正常取到直接返回
		return nil
	}
	duration, err := fetchFunc(dst)
	if err != nil {
		return err
	}
	err = cache.Set(cacheKey, dst, duration)
	if err != nil {
		return err
	}
	//SetReflectValue(dst, data)
	return nil
}

func RedisV8Cache(client func() *reidsv8.Client) Cache {
	return _RedisV8Cache{
		client: sync.OnceValue(client),
	}
}

type _RedisV8Cache struct {
	client func() *reidsv8.Client
}

func (r _RedisV8Cache) Get(key string, data any) (exists bool, err error) {
	ctx := context.Background()
	b, err := r.client().Get(ctx, key).Bytes()
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
	_, err = r.client().Set(ctx, key, b, duration).Result()
	if err != nil {
		return err
	}
	return nil
}

func RedisV9Cache(client func() *reidsv9.Client) Cache {
	return _RedisV9Cache{
		client: sync.OnceValue(client),
	}
}

type _RedisV9Cache struct {
	client func() *reidsv9.Client
}

func (r _RedisV9Cache) Get(key string, data any) (exists bool, err error) {
	ctx := context.Background()
	b, err := r.client().Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, reidsv9.Nil) { // 是redis.Nil 错误，屏蔽错误，exists 返回false
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
	_, err = r.client().Set(ctx, key, b, duration).Result()
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
	SetReflectValue(dst, resulany)
	return true, nil
}

func (m _MemeryCache) Set(key string, data any, duration time.Duration) error {
	resulany := reflect.New(reflect.TypeOf(data).Elem()).Interface()
	SetReflectValue(resulany, data)
	memeryCache.Set(key, resulany, duration) // 这个地方需要复制一份data,其效果有待验证
	return nil
}

// Md5Lower md5 小写
func Md5Lower(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// SetReflectValue 设置反射值，如果类型不匹配，会自动转换
func SetReflectValue(dst any, src any) {
	rdst := reflect.Indirect(reflect.ValueOf(dst))
	rsrc := reflect.Indirect(reflect.ValueOf(src))
	if rsrc.CanConvert(rdst.Type()) {
		rsrc = rsrc.Convert(rdst.Type())
	}
	rdst.Set(rsrc)
}
