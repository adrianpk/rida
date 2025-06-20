package cfg

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type ClientsConfig struct {
	OttawaQty   int
	MontrealQty int
}

type PgConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type Config struct {
	APIKey   string
	HTTPPort string
	Pg       PgConfig
	Clients  ClientsConfig
}

func Load() *Config {
	apiKey := flag.String("api-key", getenv("RIDA_API_KEY", "demo-api-key"), "API key for simulated clients")
	ottawaQty := flag.Int("ottawa-clients", getenvInt("RIDA_OTTAWA_CLIENTS", 1), "Number of Ottawa clients")
	montrealQty := flag.Int("montreal-clients", getenvInt("RIDA_MONTREAL_CLIENTS", 2), "Number of Montreal clients")
	httpPort := flag.String("http-port", getenv("RIDA_HTTP_PORT", ":8080"), "HTTP server port (e.g. :8080)")
	pgHost := flag.String("pg-host", getenv("RIDA_PG_HOST", "localhost"), "Postgres host")
	pgPort := flag.String("pg-port", getenv("RIDA_PG_PORT", "5432"), "Postgres port")
	pgUser := flag.String("pg-user", getenv("RIDA_PG_USER", "postgres"), "Postgres user")
	pgPassword := flag.String("pg-password", getenv("RIDA_PG_PASSWORD", "postgres"), "Postgres password")
	pgDBName := flag.String("pg-dbname", getenv("RIDA_PG_DBNAME", "rida"), "Postgres database name")
	pgSSLMode := flag.String("pg-sslmode", getenv("RIDA_PG_SSLMODE", "disable"), "Postgres SSL mode")
	flag.Parse()

	return &Config{
		APIKey:   *apiKey,
		HTTPPort: *httpPort,
		Clients: ClientsConfig{
			OttawaQty:   *ottawaQty,
			MontrealQty: *montrealQty,
		},
		Pg: PgConfig{
			Host:     *pgHost,
			Port:     *pgPort,
			User:     *pgUser,
			Password: *pgPassword,
			DBName:   *pgDBName,
			SSLMode:  *pgSSLMode,
		},
	}
}

func (pg *PgConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		pg.Host, pg.Port, pg.User, pg.Password, pg.DBName, pg.SSLMode,
	)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

func getenvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}

	return fallback
}
