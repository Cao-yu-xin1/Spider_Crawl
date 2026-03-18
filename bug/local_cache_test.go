package cache

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestLocalCache_Clear(t *testing.T) {
	type fields struct {
		items map[string]*CacheItem
		mu    sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LocalCache{
				items: tt.fields.items,
				mu:    tt.fields.mu,
			}
			c.Clear()
		})
	}
}

func TestLocalCache_Decrement(t *testing.T) {
	type fields struct {
		items map[string]*CacheItem
		mu    sync.RWMutex
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LocalCache{
				items: tt.fields.items,
				mu:    tt.fields.mu,
			}
			if got := c.Decrement(tt.args.key); got != tt.want {
				t.Errorf("Decrement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalCache_Delete(t *testing.T) {
	type fields struct {
		items map[string]*CacheItem
		mu    sync.RWMutex
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LocalCache{
				items: tt.fields.items,
				mu:    tt.fields.mu,
			}
			c.Delete(tt.args.key)
		})
	}
}

func TestLocalCache_DeleteExpired(t *testing.T) {
	type fields struct {
		items map[string]*CacheItem
		mu    sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LocalCache{
				items: tt.fields.items,
				mu:    tt.fields.mu,
			}
			if got := c.DeleteExpired(); got != tt.want {
				t.Errorf("DeleteExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalCache_DeleteMulti(t *testing.T) {
	type fields struct {
		items map[string]*CacheItem
		mu    sync.RWMutex
	}
	type args struct {
		keys []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LocalCache{
				items: tt.fields.items,
				mu:    tt.fields.mu,
			}
			c.DeleteMulti(tt.args.keys)
		})
	}
}

func TestLocalCache_Get(t *testing.T) {
	type fields struct {
		items map[string]*CacheItem
		mu    sync.RWMutex
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LocalCache{
				items: tt.fields.items,
				mu:    tt.fields.mu,
			}
			got, got1 := c.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestLocalCache_GetMulti(t *testing.T) {
	type fields struct {
		items map[string]*CacheItem
		mu    sync.RWMutex
	}
	type args struct {
		keys []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LocalCache{
				items: tt.fields.items,
				mu:    tt.fields.mu,
			}
			if got := c.GetMulti(tt.args.keys); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMulti() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalCache_GetOrSet(t *testing.T) {
	type fields struct {
		items map[string]*CacheItem
		mu    sync.RWMutex
	}
	type args struct {
		key   string
		value interface{}
		ttl   time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LocalCache{
				items: tt.fields.items,
				mu:    tt.fields.mu,
			}
			if got := c.GetOrSet(tt.args.key, tt.args.value, tt.args.ttl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOrSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalCache_Has(t *testing.T) {
	type fields struct {
		items map[string]*CacheItem
		mu    sync.RWMutex
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LocalCache{
				items: tt.fields.items,
				mu:    tt.fields.mu,
			}
			if got := c.Has(tt.args.key); got != tt.want {
				t.Errorf("Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalCache_Increment(t *testing.T) {
	type fields struct {
		items map[string]*CacheItem
		mu    sync.RWMutex
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LocalCache{
				items: tt.fields.items,
				mu:    tt.fields.mu,
			}
			if got := c.Increment(tt.args.key); got != tt.want {
				t.Errorf("Increment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalCache_Keys(t *testing.T) {
	type fields struct {
		items map[string]*CacheItem
		mu    sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LocalCache{
				items: tt.fields.items,
				mu:    tt.fields.mu,
			}
			if got := c.Keys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalCache_Range(t *testing.T) {
	type fields struct {
		items map[string]*CacheItem
		mu    sync.RWMutex
	}
	type args struct {
		fn func(key string, value interface{}) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LocalCache{
				items: tt.fields.items,
				mu:    tt.fields.mu,
			}
			c.Range(tt.args.fn)
		})
	}
}

func TestLocalCache_Set(t *testing.T) {
	type fields struct {
		items map[string]*CacheItem
		mu    sync.RWMutex
	}
	type args struct {
		key   string
		value interface{}
		ttl   time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LocalCache{
				items: tt.fields.items,
				mu:    tt.fields.mu,
			}
			c.Set(tt.args.key, tt.args.value, tt.args.ttl)
		})
	}
}

func TestLocalCache_SetMulti(t *testing.T) {
	type fields struct {
		items map[string]*CacheItem
		mu    sync.RWMutex
	}
	type args struct {
		items map[string]interface{}
		ttl   time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LocalCache{
				items: tt.fields.items,
				mu:    tt.fields.mu,
			}
			c.SetMulti(tt.args.items, tt.args.ttl)
		})
	}
}

func TestLocalCache_Size(t *testing.T) {
	type fields struct {
		items map[string]*CacheItem
		mu    sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LocalCache{
				items: tt.fields.items,
				mu:    tt.fields.mu,
			}
			if got := c.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalCache_Values(t *testing.T) {
	type fields struct {
		items map[string]*CacheItem
		mu    sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   []interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LocalCache{
				items: tt.fields.items,
				mu:    tt.fields.mu,
			}
			if got := c.Values(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLocalCache(t *testing.T) {
	tests := []struct {
		name string
		want *LocalCache
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLocalCache(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLocalCache() = %v, want %v", got, tt.want)
			}
		})
	}
}
