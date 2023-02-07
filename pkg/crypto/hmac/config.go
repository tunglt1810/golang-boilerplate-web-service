package hmac

const (
	defaultHeader = "SM-HMAC-SHA256" //nolint
)

type Config struct {
	Algorithm        string   `mapstructure:"algorithm"`
	Secret           string   `mapstructure:"secret_key"`
	Headers          []string `mapstructure:"headers"`
	ValidateBody     bool     `mapstructure:"validate_body"`
	AllowedClockSkew int64    `mapstructure:"allowed_clock_skew"`
	EnforcedHeaders  []string `mapstructure:"enforced_headers"`
	AccessKey        string   `mapstructure:"access_key,omitempty"`
}
