package cache_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	errors "github.com/eclipse-xfsc/microservice-core-go/pkg/err"
	ptr "github.com/eclipse-xfsc/microservice-core-go/pkg/ptr"
	goacache "github.com/eclipse-xfsc/redis-cache-service/gen/cache"
	"github.com/eclipse-xfsc/redis-cache-service/internal/service/cache"
	"github.com/eclipse-xfsc/redis-cache-service/internal/service/cache/cachefakes"
)

func TestNew(t *testing.T) {
	svc := cache.New(nil, nil, zap.NewNop())
	assert.Implements(t, (*goacache.Service)(nil), svc)
}

func TestService_Get(t *testing.T) {
	const key1 = "key,namespace,scope"
	const key2 = "key,namespace,scope2"

	tests := []struct {
		name  string
		cache *cachefakes.FakeCache
		req   *goacache.CacheGetRequest

		res        interface{}
		errkind    errors.Kind
		errtext    string
		loggerText string
	}{
		{
			name:    "missing cache key",
			req:     &goacache.CacheGetRequest{},
			errkind: errors.BadRequest,
			errtext: "missing key",
		},
		{
			name: "error getting value from cache",
			req: &goacache.CacheGetRequest{
				Key:       "key",
				Namespace: ptr.String("namespace"),
				Scope:     ptr.String("scope"),
			},
			cache: &cachefakes.FakeCache{
				GetStub: func(ctx context.Context, key string) ([]byte, error) {
					return nil, errors.New("some error")
				},
			},
			errkind: errors.Unknown,
			errtext: "some error",
		},
		{
			name: "key not found in cache",
			req: &goacache.CacheGetRequest{
				Key:       "key",
				Namespace: ptr.String("namespace"),
				Scope:     ptr.String("scope"),
			},
			cache: &cachefakes.FakeCache{
				GetStub: func(ctx context.Context, key string) ([]byte, error) {
					return nil, errors.New(errors.NotFound)
				},
			},
			errkind: errors.NotFound,
			errtext: "key not found in cache",
		},
		{
			name: "value returned from cache is not json",
			req: &goacache.CacheGetRequest{
				Key:       "key",
				Namespace: ptr.String("namespace"),
				Scope:     ptr.String("scope"),
			},
			cache: &cachefakes.FakeCache{
				GetStub: func(ctx context.Context, key string) ([]byte, error) {
					return []byte("boom"), nil
				},
			},
			errkind: errors.Unknown,
			errtext: "cannot decode json value from cache",
		},
		{
			name: "json value is successfully returned from cache",
			req: &goacache.CacheGetRequest{
				Key:       "key",
				Namespace: ptr.String("namespace"),
				Scope:     ptr.String("scope"),
			},
			cache: &cachefakes.FakeCache{
				GetStub: func(ctx context.Context, key string) ([]byte, error) {
					return []byte(`{"test":"value"}`), nil
				},
			},
			res:     map[string]interface{}{"test": "value"},
			errtext: "",
		},
		{
			name: "multiple scope cache return error",
			req: &goacache.CacheGetRequest{
				Key:       "key",
				Namespace: ptr.String("namespace"),
				Scope:     ptr.String("scope,scope2"),
				Strategy:  ptr.String("last"),
			},
			cache: &cachefakes.FakeCache{
				GetStub: func(ctx context.Context, key string) ([]byte, error) {
					if key == key1 {
						return nil, fmt.Errorf("some error")
					}
					return []byte(`{"test":"value2"}`), nil
				},
			},
			errtext: "error getting value from cache",
		},
		{
			name: "multiple scope with merge flatten strategy",
			req: &goacache.CacheGetRequest{
				Key:       "key",
				Namespace: ptr.String("namespace"),
				Scope:     ptr.String("scope,scope2"),
				Strategy:  ptr.String("merge"),
			},
			cache: &cachefakes.FakeCache{
				GetStub: func(ctx context.Context, key string) ([]byte, error) {
					if key == key1 {
						return []byte(`{"test":"value"}`), nil
					}
					return []byte(`{"test":"value2"}`), nil
				},
			},
			res:     map[string]interface{}{"test_1": "value", "test_2": "value2"},
			errtext: "",
		},
		{
			name: "multiple scope with first flatten strategy",
			req: &goacache.CacheGetRequest{
				Key:       "key",
				Namespace: ptr.String("namespace"),
				Scope:     ptr.String("scope,scope2"),
				Strategy:  ptr.String("first"),
			},
			cache: &cachefakes.FakeCache{
				GetStub: func(ctx context.Context, key string) ([]byte, error) {
					if key == key1 {
						return []byte(`{"test":"value"}`), nil
					}
					return []byte(`{"test":"value2"}`), nil
				},
			},
			res:     map[string]interface{}{"test": "value"},
			errtext: "",
		},
		{
			name: "multiple scope with last flatten strategy",
			req: &goacache.CacheGetRequest{
				Key:       "key",
				Namespace: ptr.String("namespace"),
				Scope:     ptr.String("scope,not_existed_scope,scope2"),
				Strategy:  ptr.String("last"),
			},
			cache: &cachefakes.FakeCache{
				GetStub: func(ctx context.Context, key string) ([]byte, error) {
					if key == key1 {
						return []byte(`{"test":"value"}`), nil
					}
					return []byte(`{"test":"value2"}`), nil
				},
			},
			res:     map[string]interface{}{"test": "value2"},
			errtext: "",
		},
		{
			name: "multiple scope with last flatten strategy",
			req: &goacache.CacheGetRequest{
				Key:       "key",
				Namespace: ptr.String("namespace"),
				Scope:     ptr.String("scope,scope2"),
				Strategy:  ptr.String("last"),
			},
			cache: &cachefakes.FakeCache{
				GetStub: func(ctx context.Context, key string) ([]byte, error) {
					if key == key1 {
						return []byte(`{"test":"value"}`), nil
					}
					return []byte(`{"test":"value2"}`), nil
				},
			},
			res:     map[string]interface{}{"test": "value2"},
			errtext: "",
		},
		{
			name: "multiple scope return warn if the key doesn't exist",
			req: &goacache.CacheGetRequest{
				Key:       "key",
				Namespace: ptr.String("namespace"),
				Scope:     ptr.String("scope,scope2"),
				Strategy:  ptr.String("last"),
			},
			cache: &cachefakes.FakeCache{
				GetStub: func(ctx context.Context, key string) ([]byte, error) {
					if key == key1 {
						return []byte(`{"test":"value"}`), nil
					}
					if key == key2 {
						return []byte(`{"test":"value2"}`), nil
					}
					return nil, errors.New(errors.NotFound, "some error")
				},
			},
			res:        map[string]interface{}{"test": "value2"},
			loggerText: "key not found in cache",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			core, logs := observer.New(zap.WarnLevel)
			logger := zap.New(core)

			svc := cache.New(test.cache, nil, logger)
			res, err := svc.Get(context.Background(), test.req)

			if test.loggerText != "" && logs.Len() >= 1 {
				assert.Contains(t, logs.All()[0].Message, test.loggerText)
			}

			if err == nil {
				assert.Empty(t, test.errtext)
				assert.Equal(t, test.res, res)
			} else {
				assert.Nil(t, res)
				assert.Error(t, err)

				e, ok := err.(*errors.Error)
				assert.True(t, ok)
				assert.Equal(t, test.errkind, e.Kind)
				assert.Contains(t, e.Error(), test.errtext)
			}
		})
	}
}

func TestService_Set(t *testing.T) {
	tests := []struct {
		name  string
		cache *cachefakes.FakeCache
		req   *goacache.CacheSetRequest

		res     interface{}
		errkind errors.Kind
		errtext string
	}{
		{
			name:    "missing cache key",
			req:     &goacache.CacheSetRequest{},
			errkind: errors.BadRequest,
			errtext: "missing key",
		},
		{
			name: "error setting value in cache",
			req: &goacache.CacheSetRequest{
				Key:       "key",
				Namespace: ptr.String("namespace"),
				Scope:     ptr.String("scope"),
				Data:      map[string]interface{}{"test": "value"},
			},
			cache: &cachefakes.FakeCache{
				SetStub: func(ctx context.Context, key string, value []byte, ttl time.Duration) error {
					return errors.New(errors.Timeout, "some error")
				},
			},
			errkind: errors.Timeout,
			errtext: "some error",
		},
		{
			name: "successfully set value in cache",
			req: &goacache.CacheSetRequest{
				Key:       "key",
				Namespace: ptr.String("namespace"),
				Scope:     ptr.String("scope"),
				Data:      map[string]interface{}{"test": "value"},
			},
			cache: &cachefakes.FakeCache{
				SetStub: func(ctx context.Context, key string, value []byte, ttl time.Duration) error {
					return nil
				},
			},
			errtext: "",
		},
		{
			name: "successfully set value in cache with TTL provided in request",
			req: &goacache.CacheSetRequest{
				Key:       "key",
				Namespace: ptr.String("namespace"),
				Scope:     ptr.String("scope"),
				Data:      map[string]interface{}{"test": "value"},
				TTL:       ptr.Int(60),
			},
			cache: &cachefakes.FakeCache{
				SetStub: func(ctx context.Context, key string, value []byte, ttl time.Duration) error {
					return nil
				},
			},
			errtext: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			svc := cache.New(test.cache, nil, zap.NewNop())
			err := svc.Set(context.Background(), test.req)
			if err == nil {
				assert.Empty(t, test.errtext)
			} else {
				assert.Error(t, err)
				e, ok := err.(*errors.Error)
				assert.True(t, ok)
				assert.Equal(t, test.errkind, e.Kind)
				assert.Contains(t, e.Error(), test.errtext)
			}
		})
	}
}

func TestService_SetExternal(t *testing.T) {
	tests := []struct {
		name   string
		cache  *cachefakes.FakeCache
		events *cachefakes.FakeEvents
		req    *goacache.CacheSetRequest

		res     interface{}
		errkind errors.Kind
		errtext string
	}{
		{
			name: "error setting external input in cache",
			req: &goacache.CacheSetRequest{
				Key:       "key",
				Namespace: ptr.String("namespace"),
				Scope:     ptr.String("scope"),
				Data:      map[string]interface{}{"test": "value"},
			},
			cache: &cachefakes.FakeCache{
				SetStub: func(ctx context.Context, key string, value []byte, ttl time.Duration) error {
					return errors.New(errors.Timeout, "some error")
				},
			},
			errkind: errors.Timeout,
			errtext: "some error",
		},
		{
			name: "error sending an event for external input",
			req: &goacache.CacheSetRequest{
				Key:       "key",
				Namespace: ptr.String("namespace"),
				Scope:     ptr.String("scope"),
				Data:      map[string]interface{}{"test": "value"},
			},
			cache: &cachefakes.FakeCache{
				SetStub: func(ctx context.Context, key string, value []byte, ttl time.Duration) error {
					return nil
				},
			},
			events: &cachefakes.FakeEvents{SendStub: func(ctx context.Context, s string) error {
				return errors.New(errors.Unknown, "failed to send event")
			}},
			errkind: errors.Unknown,
			errtext: "failed to send event",
		},
		{
			name: "successfully set value in cache and send an event to events",
			req: &goacache.CacheSetRequest{
				Key:       "key",
				Namespace: ptr.String("namespace"),
				Scope:     ptr.String("scope"),
				Data:      map[string]interface{}{"test": "value"},
			},
			cache: &cachefakes.FakeCache{
				SetStub: func(ctx context.Context, key string, value []byte, ttl time.Duration) error {
					return nil
				},
			},
			events: &cachefakes.FakeEvents{SendStub: func(ctx context.Context, s string) error {
				return nil
			}},
			errtext: "",
		},
		{
			name: "successfully set value in cache with TTL provided in request and send an event to events",
			req: &goacache.CacheSetRequest{
				Key:       "key",
				Namespace: ptr.String("namespace"),
				Scope:     ptr.String("scope"),
				Data:      map[string]interface{}{"test": "value"},
				TTL:       ptr.Int(60),
			},
			cache: &cachefakes.FakeCache{
				SetStub: func(ctx context.Context, key string, value []byte, ttl time.Duration) error {
					return nil
				},
			},
			events: &cachefakes.FakeEvents{SendStub: func(ctx context.Context, s string) error {
				return nil
			}},
			errtext: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			svc := cache.New(test.cache, test.events, zap.NewNop())
			err := svc.SetExternal(context.Background(), test.req)
			if err == nil {
				assert.Empty(t, test.errtext)
			} else {
				assert.Error(t, err)
				e, ok := err.(*errors.Error)
				assert.True(t, ok)
				assert.Equal(t, test.errkind, e.Kind)
				assert.Contains(t, e.Error(), test.errtext)
			}
		})
	}
}
