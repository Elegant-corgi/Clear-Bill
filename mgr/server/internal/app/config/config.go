package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

const defaultConfigPath = "configs/config.toml"

type HTTP struct {
	Host            string
	Port            int
	ServeTimeout    int
	ShutdownTimeout int
}

type GORM struct {
	Host            string
	Port            int
	Database        string
	User            string
	Password        string
	Charset         string
	ParseTime       bool
	Loc             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
}

func (h HTTP) Address() string {
	host := h.Host
	if host == "" {
		host = "0.0.0.0"
	}

	port := h.Port
	if port == 0 {
		port = 8080
	}

	return fmt.Sprintf("%s:%d", host, port)
}

type Config struct {
	RunMode     string
	WWW         string
	PrintConfig bool
	HTTP        HTTP
	GORM        GORM

	HTTPAddr string `toml:"-" json:"-"`
	WebRoot  string `toml:"-" json:"-"`
}

var (
	C    = new(Config)
	once sync.Once
)

func Load() Config {
	configPath := envOrDefault("CLEAR_BILL_CONFIG", defaultConfigPath)
	MustLoad(configPath)
	return *C
}

func MustLoad(configPath string) {
	once.Do(func() {
		cfg := defaultConfig()

		content, err := os.ReadFile(configPath)
		if err != nil {
			panic(fmt.Errorf("read config file %q failed: %w", configPath, err))
		}

		if err := toml.Unmarshal(content, &cfg); err != nil {
			panic(fmt.Errorf("parse config file %q failed: %w", configPath, err))
		}

		cfg.HTTPAddr = envOrDefault("CLEAR_BILL_HTTP_ADDR", cfg.HTTP.Address())
		cfg.WebRoot = envOrDefault("CLEAR_BILL_WEB_ROOT", cfg.GetWWW())
		*C = cfg

		PrintWithJSON()
	})
}

func defaultConfig() Config {
	return Config{
		RunMode:     "debug",
		WWW:         "website",
		PrintConfig: false,
		HTTP: HTTP{
			Host:            "0.0.0.0",
			Port:            8080,
			ServeTimeout:    5,
			ShutdownTimeout: 5,
		},
		GORM: GORM{
			Host:            "127.0.0.1",
			Port:            3306,
			Database:        "clear_bill",
			User:            "root",
			Password:        "root",
			Charset:         "utf8mb4",
			ParseTime:       true,
			Loc:             "Local",
			MaxIdleConns:    10,
			MaxOpenConns:    50,
			ConnMaxLifetime: 300,
		},
	}
}

func (c Config) GetWWW() string {
	if c.WWW == "" {
		return "website"
	}

	return c.WWW
}

func (c Config) IsDebugMode() bool {
	return c.RunMode == "debug"
}

func PrintWithJSON() {
	if !C.PrintConfig {
		return
	}

	payload, err := json.MarshalIndent(C, "", "  ")
	if err != nil {
		os.Stderr.WriteString("[CONFIG] JSON marshal error: " + err.Error())
		return
	}

	os.Stdout.WriteString(string(payload) + "\n")
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
