package hmac

import (
	"sync"
)

var (
	Instance *HMAC
	HConfig  *Config
	once     sync.Once
)

//nolint:gochecknoinits
func init() {
	if Instance == nil && HConfig != nil {
		once.Do(func() {
			Instance = &HMAC{
				cfg: HConfig,
			}
		})
	}
}

func InitInstance(cfg *Config) {
	HConfig = cfg
	if Instance == nil {
		Instance = &HMAC{
			cfg: HConfig,
		}
	}
}
