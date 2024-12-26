package redis

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Redis struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Database     int           `mapstructure:"database"`
	TTL          time.Duration `mapstructure:"ttl"`
	Username     string        `mapstructure:"username"`
	Password     string        `mapstructure:"password"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	WriteTimeOut time.Duration `mapstructure:"write_timeout"`
	ReadTimeOut  time.Duration `mapstructure:"read_timeout"`
	DialTimeOut  time.Duration `mapstructure:"dial_timeout"`
	TLSConfig    *struct {
		InsecureSkipVerify bool   `mapstructure:"insecure_skip_verify"`
		CertFilePath       string `mapstructure:"cert_file_path"`
	} `mapstructure:"tls_config,omitempty"`
}

type Cache struct {
	MaxItem int           `mapstructure:"max_item"`
	TTL     time.Duration `mapstructure:"ttl"`
}

func NewRedisClient(cfgRedis *Redis) *redis.Client {
	redisOpt := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfgRedis.Host, cfgRedis.Port),
		Password:     cfgRedis.Password,
		Username:     cfgRedis.Username,
		DB:           cfgRedis.Database,
		PoolSize:     cfgRedis.PoolSize,
		MinIdleConns: cfgRedis.MinIdleConns,
		WriteTimeout: cfgRedis.WriteTimeOut,
		ReadTimeout:  cfgRedis.ReadTimeOut,
		DialTimeout:  cfgRedis.DialTimeOut,
	}

	if cfgRedis.TLSConfig != nil {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: cfgRedis.TLSConfig.InsecureSkipVerify, //nolint:gosec
			MinVersion:         tls.VersionTLS12,
		}

		if !cfgRedis.TLSConfig.InsecureSkipVerify {
			caCert, err := os.ReadFile(cfgRedis.TLSConfig.CertFilePath)
			if err != nil {
				log.Fatal().Caller().Err(err).Send()
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			tlsConfig.RootCAs = caCertPool
		}
		redisOpt.TLSConfig = tlsConfig
	}

	// Connect to redis server
	client := redis.NewClient(redisOpt)
	log.Info().Caller().Msgf("Pinging to Redis Server: %s", cfgRedis.Host)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Error().Caller().Err(err).Msgf("Connect to Redis Server %s fail ", cfgRedis.Host)
	} else {
		log.Info().Caller().Msg("Connected to Redis Server")
	}
	return client
}
