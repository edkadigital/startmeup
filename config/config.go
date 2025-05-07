package config

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	// StaticDir stores the name of the directory that will serve static files.
	StaticDir = "static"

	// StaticPrefix stores the URL prefix used when serving static files.
	StaticPrefix = "files"
)

type environment string

const (
	// EnvLocal represents the local environment.
	EnvLocal environment = "local"

	// EnvTest represents the test environment.
	EnvTest environment = "test"

	// EnvDevelop represents the development environment.
	EnvDevelop environment = "dev"

	// EnvStaging represents the staging environment.
	EnvStaging environment = "staging"

	// EnvQA represents the qa environment.
	EnvQA environment = "qa"

	// EnvProduction represents the production environment.
	EnvProduction environment = "prod"
)

// SwitchEnvironment sets the environment variable used to dictate which environment the application is
// currently running in.
// This must be called prior to loading the configuration in order for it to take effect.
func SwitchEnvironment(env environment) {
	if err := os.Setenv("PAGODA_APP_ENVIRONMENT", string(env)); err != nil {
		panic(err)
	}
}

type (
	// Config stores complete configuration.
	Config struct {
		HTTP     HTTPConfig
		App      AppConfig
		Cache    CacheConfig
		Database DatabaseConfig
		Files    FilesConfig
		Tasks    TasksConfig
		Mail     MailConfig
	}

	// HTTPConfig stores HTTP configuration.
	HTTPConfig struct {
		Hostname        string
		Port            uint16
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		IdleTimeout     time.Duration
		ShutdownTimeout time.Duration
		TLS             struct {
			Enabled     bool
			Certificate string
			Key         string
		}
	}

	// AppConfig stores application configuration.
	AppConfig struct {
		Name          string
		Host          string
		Environment   environment
		EncryptionKey string
		Timeout       time.Duration
		PasswordToken struct {
			Expiration time.Duration
			Length     int
		}
		EmailVerificationTokenExpiration time.Duration
	}

	// CacheConfig stores the cache configuration.
	CacheConfig struct {
		Capacity   int
		Expiration struct {
			StaticFile time.Duration
		}
	}

	// DatabaseConfig stores the database configuration.
	DatabaseConfig struct {
		Driver         string
		Connection     string
		TestConnection string
	}

	// FilesConfig stores the file system configuration.
	FilesConfig struct {
		Directory string
	}

	// TasksConfig stores the tasks configuration.
	TasksConfig struct {
		Goroutines      int
		ReleaseAfter    time.Duration
		CleanupInterval time.Duration
		ShutdownTimeout time.Duration
		PollInterval    time.Duration
	}

	// MailConfig stores the mail configuration.
	MailConfig struct {
		Hostname    string
		Port        uint16
		User        string
		Password    string
		FromAddress string
	}
)

// GetConfig loads and returns configuration.
func GetConfig() (Config, error) {
	var c Config

	// Load the config file.
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	// Load env variables.
	viper.SetEnvPrefix("pagoda")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return c, err
	}

	if err := viper.Unmarshal(&c); err != nil {
		return c, err
	}

	// Override database config with environment variables if provided
	overrideFromEnv(&c)

	return c, nil
}

// overrideFromEnv overrides configuration with environment variables
func overrideFromEnv(c *Config) {
	// Database configuration
	if dbDriver := os.Getenv("DB_DRIVER"); dbDriver != "" {
		c.Database.Driver = dbDriver
	}

	// Only use DATABASE_URL for connection string
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		c.Database.Connection = dbURL
	}

	// Also update test connection if not explicitly set
	if os.Getenv("DB_TEST_CONNECTION") == "" && c.Database.TestConnection == "" {
		c.Database.TestConnection = c.Database.Connection
	} else if testConn := os.Getenv("DB_TEST_CONNECTION"); testConn != "" {
		c.Database.TestConnection = testConn
	}
}
