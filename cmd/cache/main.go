package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
	"golang.org/x/sync/errgroup"

	auth "github.com/eclipse-xfsc/microservice-core-go/pkg/auth"
	graceful "github.com/eclipse-xfsc/microservice-core-go/pkg/graceful"
	goacache "github.com/eclipse-xfsc/redis-cache-service/gen/cache"
	goahealth "github.com/eclipse-xfsc/redis-cache-service/gen/health"
	goacachesrv "github.com/eclipse-xfsc/redis-cache-service/gen/http/cache/server"
	goahealthsrv "github.com/eclipse-xfsc/redis-cache-service/gen/http/health/server"
	goaopenapisrv "github.com/eclipse-xfsc/redis-cache-service/gen/http/openapi/server"
	"github.com/eclipse-xfsc/redis-cache-service/gen/openapi"
	"github.com/eclipse-xfsc/redis-cache-service/internal/clients/event"
	"github.com/eclipse-xfsc/redis-cache-service/internal/clients/redis"
	"github.com/eclipse-xfsc/redis-cache-service/internal/config"
	"github.com/eclipse-xfsc/redis-cache-service/internal/service"
	"github.com/eclipse-xfsc/redis-cache-service/internal/service/cache"
	"github.com/eclipse-xfsc/redis-cache-service/internal/service/health"
)

var Version = "0.0.0+development"

func main() {
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("cannot load configuration: %v", err)
	}

	logger, err := createLogger(cfg.LogLevel)
	if err != nil {
		log.Fatalln(err)
	}
	defer logger.Sync() //nolint:errcheck

	logger.Info("start cache service", zap.String("version", Version), zap.String("goa", goa.Version()))

	// create redis client
	redis := redis.New(cfg.Redis.Addr, cfg.Redis.User, cfg.Redis.Pass, cfg.Redis.DB, cfg.Redis.TTL, cfg.Redis.Cluster)

	// create event client
	events, err := event.New(cfg.Nats.Addr, cfg.Nats.Subject)
	if err != nil {
		log.Fatalf("failed to create events client: %v", err)
	}
	defer events.CLose(context.Background())

	// create services
	var (
		cacheSvc  goacache.Service
		healthSvc goahealth.Service
	)
	{
		cacheSvc = cache.New(redis, events, logger)
		healthSvc = health.New(Version)
	}

	// create endpoints
	var (
		cacheEndpoints   *goacache.Endpoints
		healthEndpoints  *goahealth.Endpoints
		openapiEndpoints *openapi.Endpoints
	)
	{
		cacheEndpoints = goacache.NewEndpoints(cacheSvc)
		healthEndpoints = goahealth.NewEndpoints(healthSvc)
		openapiEndpoints = openapi.NewEndpoints(nil)
	}

	// Provide the transport specific request decoder and response encoder.
	// The goa http package has built-in support for JSON, XML and gob.
	// Other encodings can be used by providing the corresponding functions,
	// see goa.design/implement/encoding.
	var (
		dec = goahttp.RequestDecoder
		enc = goahttp.ResponseEncoder
	)

	// Build the service HTTP request multiplexer and configure it to serve
	// HTTP requests to the service endpoints.
	mux := goahttp.NewMuxer()

	var (
		cacheServer   *goacachesrv.Server
		healthServer  *goahealthsrv.Server
		openapiServer *goaopenapisrv.Server
	)
	{
		cacheServer = goacachesrv.New(cacheEndpoints, mux, dec, enc, nil, errFormatter)
		healthServer = goahealthsrv.New(healthEndpoints, mux, dec, enc, nil, errFormatter)
		openapiServer = goaopenapisrv.New(openapiEndpoints, mux, dec, enc, nil, errFormatter, nil, nil)
	}

	// Apply Authentication middleware if enabled
	if cfg.Auth.Enabled {
		m, err := auth.NewMiddleware(cfg.Auth.JwkURL, cfg.Auth.RefreshInterval, httpClient())
		if err != nil {
			log.Fatalf("failed to create authentication middleware: %v", err)
		}
		cacheServer.Use(m.Handler())
	}

	// Configure the mux.
	goacachesrv.Mount(mux, cacheServer)
	goahealthsrv.Mount(mux, healthServer)
	goaopenapisrv.Mount(mux, openapiServer)

	// expose metrics
	go exposeMetrics(cfg.Metrics.Addr, logger)

	var handler http.Handler = mux
	srv := &http.Server{
		Addr:         cfg.HTTP.Host + ":" + cfg.HTTP.Port,
		Handler:      handler,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		if err := graceful.Shutdown(ctx, srv, 20*time.Second); err != nil {
			logger.Error("server shutdown error", zap.Error(err))
			return err
		}
		return errors.New("server stopped successfully")
	})
	if err := g.Wait(); err != nil {
		logger.Error("run group stopped", zap.Error(err))
	}

	logger.Info("bye bye")
}

func createLogger(logLevel string, opts ...zap.Option) (*zap.Logger, error) {
	var level = zapcore.InfoLevel
	if logLevel != "" {
		err := level.UnmarshalText([]byte(logLevel))
		if err != nil {
			return nil, err
		}
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(level)
	config.DisableStacktrace = true
	config.EncoderConfig.TimeKey = "ts"
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	return config.Build(opts...)
}

func errFormatter(ctx context.Context, e error) goahttp.Statuser {
	return service.NewErrorResponse(ctx, e)
}

func exposeMetrics(addr string, logger *zap.Logger) {
	promMux := http.NewServeMux()
	promMux.Handle("/metrics", promhttp.Handler())
	logger.Info(fmt.Sprintf("exposing prometheus metrics at %s/metrics", addr))
	if err := http.ListenAndServe(addr, promMux); err != nil { //nolint:gosec
		logger.Error("error exposing prometheus metrics", zap.Error(err))
	}
}

func httpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			TLSHandshakeTimeout: 10 * time.Second,
			IdleConnTimeout:     60 * time.Second,
		},
		Timeout: 30 * time.Second,
	}
}
