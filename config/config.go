package config

import "os"

type Config struct {
    DBUser string
    DBPass string
    DBHost string
    DBPort string
    DBName string
}

func Load() *Config {
    return &Config{
        DBUser: getenv("DB_USER", "root"),
        DBPass: getenv("DB_PASS", "password"),
        DBHost: getenv("DB_HOST", "localhost"),
        DBPort: getenv("DB_PORT", "3306"),
        DBName: getenv("DB_NAME", "metro"),
    }
}

func getenv(key, fallback string) string {
    if val, ok := os.LookupEnv(key); ok {
        return val
    }
    return fallback
}
