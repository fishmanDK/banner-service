package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string         `yaml:"env" env-default:"local"`
	Postgres   PostgresConfig `yaml:"postgres"`
	HTTPServer `yaml:"http_server"`
	Clients    Clients `yaml:"clients"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type PostgresConfig struct {
	PgUser     string `yaml:"pg_user"`
	PgDatabase string `yaml:"pg_database"`
	PgHost     string `yaml:"pg_host"`
	PgPort     string `yaml:"pg_port"`
	PgSslmode  string `yaml:"pg_sslmode"`
	PgPassword string `yaml:"pg_password"`
}

type Clients struct {
	TriggerHookService TriggerHookAdapter `yaml:"triggerHookService"`
}

type TriggerHookAdapter struct {
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retriesCount"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = fetchConfigPath()
		if configPath == "" {
			log.Fatal("CONFIG_PATH is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

func (pg PostgresConfig) String() string {
	s := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", pg.PgHost, pg.PgPort, pg.PgUser, pg.PgDatabase, pg.PgPassword, pg.PgSslmode)
	return s
}
