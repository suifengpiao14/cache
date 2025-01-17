package cache_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.huishoubao.com/gopackage/cache"
)

func TestMemeryCache(t *testing.T) {
	data := map[string]string{"key1": "value1"}

	cache := cache.MemeryCache[map[string]string]()
	cache.Set("key1", data, 1)
	var data2 map[string]string

	exists, err := cache.Get("key1", &data2)
	require.NoError(t, err)
	require.True(t, exists)
	fmt.Println(data2)
}
func TestRemember(t *testing.T) {
	memoryCache := cache.MemeryCache[map[string]string]()
	key := "key1"
	var dst map[string]string
	cache.Remember(memoryCache, key, 20*time.Second, &dst, func() (map[string]string, error) {
		data := map[string]string{"key1": "value1"}
		fmt.Println("load data from db")
		return data, nil
	})
	fmt.Println(dst)
	var dst2 map[string]string
	cache.Remember(memoryCache, key, 20*time.Second, &dst2, func() (map[string]string, error) {
		data := map[string]string{"key1": "value1"}
		fmt.Println("load data from db 2")
		return data, nil
	})
	fmt.Println(dst2)

}
