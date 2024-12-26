package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"goland-boilerplate-web-service/pkg/crypto/hmac"
	"goland-boilerplate-web-service/pkg/database/postgresql"
	"goland-boilerplate-web-service/pkg/httpclient"

	"github.com/inhies/go-bytesize"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Schema struct {
	APIInfo struct {
		Version       string `mapstructure:"version"`
		LastUpdatedAt string `mapstructure:"last_updated_at"`
	} `mapstructure:"api_info"`

	Server Server `mapstructure:"server"`

	Database   postgresql.Database             `mapstructure:"postgres"`
	Sentry     Sentry                          `mapstructure:"sentry"`
	AuthServer httpclient.InternalClientConfig `mapstructure:"auth_server"`

	Logging Logging `mapstructure:"logging"`

	HMACInternal hmac.Config `mapstructure:"hmac_internal"`
}

type Server struct {
	Port    int `mapstructure:"port"` // port for http server
	Metrics int `mapstructure:"metrics"`
}

type Sentry struct {
	DSN              string  `mapstructure:"dsn"`
	ENV              string  `mapstructure:"env"`
	Debug            bool    `mapstructure:"debug"`
	TracesSampleRate float64 `mapstructure:"traces_sample_rate"`
}

type Logging struct {
	Level string `mapstructure:"level"`
}

var (
	Config      Schema
	callerCache sync.Map
)

// StringToByteSizeHookFunc returns a DecodeHookFunc that converts
// hex string to bytesize.ByteSize.
func StringToByteSizeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf(bytesize.B) {
			return data, nil
		}

		sDec, err := bytesize.Parse(data.(string))
		if err != nil {
			return nil, err
		}

		return sDec, nil
	}
}

func Init() {
	config := viper.New()
	config.SetConfigName("config")
	config.AddConfigPath(".")          // Look for config in current directory
	config.AddConfigPath("config/")    // Optionally look for config in the working directory.
	config.AddConfigPath("../config/") // Look for config needed for tests.
	config.AddConfigPath("../")        // Look for config needed for tests.

	config.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	config.AutomaticEnv()
	err := config.ReadInConfig() // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	err = config.Unmarshal(&Config, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			StringToByteSizeHookFunc(),
		),
	))
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	// setup logging
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		if caller, ok := callerCache.Load(file); ok {
			return caller.(string)
		}

		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		caller := short + ":" + strconv.Itoa(line)
		callerCache.Store(file, caller)
		return caller
	}
	loggingLevel, _ := zerolog.ParseLevel(Config.Logging.Level)
	zerolog.SetGlobalLevel(loggingLevel)
	zerolog.TimeFieldFormat = time.RFC3339Nano

	log.Debug().Caller().Msgf("%+v", Config)
}
