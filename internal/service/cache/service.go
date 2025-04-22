package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	errors "github.com/eclipse-xfsc/microservice-core-go/pkg/err"
	"github.com/eclipse-xfsc/redis-cache-service/gen/cache"
)

//go:generate counterfeiter . Cache
//go:generate counterfeiter . Events

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}

type Events interface {
	Send(ctx context.Context, key string) error
}

type Service struct {
	cache  Cache
	events Events
	logger *zap.Logger
}

func New(cache Cache, events Events, logger *zap.Logger) *Service {
	return &Service{
		cache:  cache,
		events: events,
		logger: logger,
	}
}

func (s *Service) Get(ctx context.Context, req *cache.CacheGetRequest) (interface{}, error) {
	logger := s.logger.With(zap.String("operation", "get"))

	if req.Key == "" {
		logger.Error("bad request: missing key")
		return nil, errors.New(errors.BadRequest, "missing key")
	}

	var scopes []string
	if req.Scope != nil {
		scopes = strings.Split(*req.Scope, ",")
	}

	if len(scopes) > 1 {
		return s.getWithMultipleScopes(ctx, req, scopes)
	}

	decodedValue, err := s.get(ctx, req.Key, req.Namespace, req.Scope)
	if err != nil {
		logger.Error("error getting value from cache", zap.Error(err))
		return nil, err
	}

	return decodedValue, nil
}

func (s *Service) Set(ctx context.Context, req *cache.CacheSetRequest) error {
	logger := s.logger.With(zap.String("operation", "set"))

	if req.Key == "" {
		logger.Error("bad request: missing key")
		return errors.New(errors.BadRequest, "missing key")
	}

	// create key from the input fields
	key := makeCacheKey(req.Key, req.Namespace, req.Scope)
	// encode payload to json bytes for storing in cache
	value, err := json.Marshal(req.Data)
	if err != nil {
		logger.Error("error encode payload to json", zap.Error(err))
		return errors.New(errors.BadRequest, "cannot encode payload to json", err)
	}

	// set cache ttl if provided in request
	var ttl time.Duration
	if req.TTL != nil {
		ttl = time.Duration(*req.TTL) * time.Second
	}

	if err := s.cache.Set(ctx, key, value, ttl); err != nil {
		logger.Error("error storing value in cache", zap.Error(err))
		return errors.New("error storing value in cache", err)
	}

	return nil
}

// SetExternal sets an external JSON value in the cache and provide an event for the input.
func (s *Service) SetExternal(ctx context.Context, req *cache.CacheSetRequest) error {
	logger := s.logger.With(zap.String("operation", "setExternal"))

	// set value in cache
	if err := s.Set(ctx, req); err != nil {
		logger.Error("error setting external input in cache", zap.Error(err))
		return errors.New("error setting external input in cache", err)
	}

	// create key from the input fields
	key := makeCacheKey(req.Key, req.Namespace, req.Scope)

	// send an event for the input
	if err := s.events.Send(ctx, key); err != nil {
		logger.Error("error sending an event for external input", zap.Error(err))
		return errors.New("error sending an event for external input", err)
	}

	return nil
}

func makeCacheKey(key string, namespace, scope *string) string {
	k := key
	if namespace != nil && *namespace != "" {
		k += "," + *namespace
	}
	if scope != nil && *scope != "" {
		k += "," + *scope
	}
	return k
}

func (s *Service) getWithMultipleScopes(ctx context.Context, req *cache.CacheGetRequest, scopes []string) (map[string]interface{}, error) {
	keyValues := map[string][]interface{}{}
	result := map[string]interface{}{}

	for _, scope := range scopes {
		scope := strings.TrimSpace(scope)
		decodedValue, err := s.get(ctx, req.Key, req.Namespace, &scope)
		if err != nil {
			if errors.Is(errors.NotFound, err) {
				s.logger.Warn(err.Error())
				continue
			}
			return nil, err
		}

		switch d := decodedValue.(type) {
		case map[string]interface{}:
			addValue(d, keyValues)
		case []map[string]interface{}:
			for _, data := range d {
				addValue(data, keyValues)
			}
		default:
			s.logger.Warn("decode value is of unknown type")
			continue
		}
	}

	switch {
	case req.Strategy == nil || *req.Strategy == "merge":
		result = mergeAll(keyValues)

	case *req.Strategy == "first":
		for key, value := range keyValues {
			result[key] = value[0]
		}

	case *req.Strategy == "last":
		for key, value := range keyValues {
			result[key] = value[len(value)-1]
		}
	}

	return result, nil
}

func addValue(decodedValue map[string]interface{}, des map[string][]interface{}) {
	for k, v := range decodedValue {
		if dataArr, contains := des[k]; contains {
			dataArr = append(dataArr, v)
			des[k] = dataArr
			continue
		}

		des[k] = []interface{}{v}
	}
}

// mergeAll merges all values for a key if more than one value is available.
func mergeAll(data map[string][]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for key, value := range data {
		if len(value) == 1 {
			result[key] = value[0]
			continue
		}

		for i, v := range value {
			index := fmt.Sprintf("%s_%d", key, i+1)
			result[index] = v
		}
	}
	return result
}

func (s *Service) get(ctx context.Context, key string, namespace *string, scope *string) (interface{}, error) {
	// create key from the input fields
	cacheKey := makeCacheKey(key, namespace, scope)
	data, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		if errors.Is(errors.NotFound, err) {
			return nil, errors.New(errors.NotFound, "key not found in cache", err)
		}
		return nil, errors.New("error getting value from cache", err)
	}

	decodedValue, err := unmarshalCacheData(data)
	if err != nil {
		return nil, errors.New("cannot decode json value from cache", err)
	}

	return decodedValue, nil
}

func unmarshalCacheData(data []byte) (interface{}, error) {
	var keyValueArray []map[string]interface{}
	var keyValue map[string]interface{}

	err := json.Unmarshal(data, &keyValue)
	if err != nil {
		err := json.Unmarshal(data, &keyValueArray)
		if err != nil {
			return nil, err
		}
		return keyValueArray, nil
	}
	return keyValue, nil
}
