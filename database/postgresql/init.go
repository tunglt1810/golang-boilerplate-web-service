package postgresql

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type TLSConfig struct {
	ClientCert []byte `mapstructure:"client_cert"`
	ClientKey  []byte `mapstructure:"client_key"`
	CACert     []byte `mapstructure:"ca_cert"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"database"`
	Schema   string `mapstructure:"schema"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"sslmode"`

	MaxIdleConnection    int           `mapstructure:"max_idle_conns"`
	MaxActiveConnection  int           `mapstructure:"max_active_conns"`
	MaxConnectionTimeout time.Duration `mapstructure:"max_conn_timeout"`

	SQLDebugger struct {
		ShowSQLStatement bool `mapstructure:"show_sql_statement"`
	} `mapstructure:"sql_debugger"`

	TLS             *TLSConfig `mapstructure:"tls"`
	MinimumPoolSize int        `mapstructure:"minimum_pool_size"`
}

func NewClient(cfg *Database) (*gorm.DB, error) {
	dgConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: fmt.Sprintf("%s.", cfg.Schema),
		},
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SSLMode,
	)

	client, err := gorm.Open(postgres.Open(dsn), dgConfig)
	if err != nil {
		return nil, err
	}

	db, err := client.DB()
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(cfg.MaxIdleConnection)
	db.SetMaxOpenConns(cfg.MaxActiveConnection)
	db.SetConnMaxIdleTime(cfg.MaxConnectionTimeout)
	if err != nil {
		log.Err(err).Caller().Msgf("Cannot connect to PostgresSQL %s", cfg.Host)
		return nil, err
	}

	log.Debug().Caller().Msgf("show sql %v", cfg.SQLDebugger.ShowSQLStatement)
	if cfg.SQLDebugger.ShowSQLStatement {
		client = client.Debug()
	}
	return client, nil
}
