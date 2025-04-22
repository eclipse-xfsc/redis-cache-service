package config

import "time"

type Config struct {
	HTTP    httpConfig
	Redis   redisConfig
	Nats    natsConfig
	Metrics metricsConfig
	Auth    authConfig

	LogLevel string `envconfig:"LOG_LEVEL" default:"INFO"`
}

// HTTP server configuration
type httpConfig struct {
	Host         string        `envconfig:"HTTP_HOST"`
	Port         string        `envconfig:"HTTP_PORT" default:"8080"`
	IdleTimeout  time.Duration `envconfig:"HTTP_IDLE_TIMEOUT" default:"120s"`
	ReadTimeout  time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"10s"`
	WriteTimeout time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"10s"`
}

type redisConfig struct {
	// Addr specifies network address of Redis server
	Addr    string        `envconfig:"REDIS_ADDR" required:"true"`
	User    string        `envconfig:"REDIS_USER" required:"true"`
	Pass    string        `envconfig:"REDIS_PASS" required:"true"`
	DB      int           `envconfig:"REDIS_DB" default:"0"`
	TTL     time.Duration `envconfig:"REDIS_EXPIRATION"` //  no default expiration, keys are set to live forever
	Cluster bool          `envconfig:"REDIS_CLUSTER" default:"false"`
}

type natsConfig struct {
	// Addr specifies network address of NATS server
	Addr string `envconfig:"NATS_ADDR" required:"true"`
	// Subject specifies NATS subject to publish events to
	Subject string `envconfig:"NATS_SUBJECT" default:"external"`
}

type metricsConfig struct {
	// Addr specifies the address to expose prometheus metrics
	Addr string `envconfig:"METRICS_ADDR" default:":2112"`
}

type authConfig struct {
	// Enabled specifies whether authentication is enabled
	Enabled         bool          `envconfig:"AUTH_ENABLED" default:"false"`
	JwkURL          string        `envconfig:"AUTH_JWK_URL"`
	RefreshInterval time.Duration `envconfig:"AUTH_REFRESH_INTERVAL" default:"1h"`
}
