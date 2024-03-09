package config

import (
	"sync"

	"github.com/caarlos0/env/v9"
	"github.com/gofiber/fiber/v3/log"
)

var Env struct {
	Database struct {
		Host               string `env:"DATABASE_HOST" envDefault:"localhost"`
		Port               int    `env:"DATABASE_PORT" envDefault:"5432"`
		User               string `env:"DATABASE_USER" envDefault:"root"`
		Password           string `env:"DATABASE_PASSWORD" envDefault:"1234"`
		DB                 string `env:"DATABASE_DB" envDefault:"rinha2024q1"`
		Timezone           string `env:"DATABASE_TIMEZONE" envDefault:"America/Sao_Paulo"`
		MaxConnections     int    `env:"DATABASE_MAX_CONNECTIONS" envDefault:"30"`
		MaxIdleConnections int    `env:"DATABASE_MAX_IDLE_CONNECTIONS" envDefault:"10"`
		SSLMode            string `env:"DATABASE_SSL_MODE" envDefault:"disable"`
		AppName            string `env:"DATABASE_APP_NAME" envDefault:"rinha2024q1"`
		LogLevel           string `env:"DATABASE_LOG_LEVEL" envDefault:"ERROR"`
		ConnectionTimeout  int    `env:"DATABASE_CONNECTION_TIMEOUT" envDefault:"30"`
		CommandTimeout     int    `env:"DATABASE_COMMAND_TIMEOUT" envDefault:"30"`
	}

	Cache struct {
		Host string `env:"CACHE_HOST" envDefault:"localhost"`
		Port int    `env:"CACHE_PORT" envDefault:"6379"`
	}

	Server struct {
		Port      int       `env:"PORT" envDefault:"8080"`
		Prefork   bool      `env:"PREFORK_ENABLED" envDefault:"false"`
		LogLevel  log.Level `env:"LOG_LEVEL" envDefault:"2"`
		Profiling struct {
			Enabled        bool   `env:"PROFILING_ENABLED" envDefault:"false"`
			CpuFilePath    string `env:"PROFILING_CPU_FILEPATH" envDefault:"/app/prof/cpu.prof"`
			MemoryFilePath string `env:"PROFILING_MEMORY_FILEPATH" envDefault:"/app/prof/memory.prof"`
		}
	}
}

var once sync.Once

func init() {
	once.Do(func() {
		loadEnvVars()
	})
}

func loadEnvVars() error {
	if err := env.Parse(&Env); err != nil {
		return err
	}

	return nil
}
