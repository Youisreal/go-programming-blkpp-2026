package config

import (
	"errors"
	"net/url"
	"os"
)

type Config struct {
	DatabaseURL string
	Addr        string
}

func Load() (Config, error) {
	values, err := readEnvFile(".env")
	if err != nil {
		return Config{}, err
	}
	get := func(key string) string {
		if value, ok := os.LookupEnv(key); ok {
			return value
		}
		return values[key]
	}
	dsn, err := DatabaseURL(get("DATABASE_URL"))
	if err != nil {
		return Config{}, err
	}
	addr := get("ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	return Config{DatabaseURL: dsn, Addr: addr}, nil
}

func DatabaseURL(value string) (string, error) {
	if value == "" {
		return "", errors.New("DATABASE_URL belum diisi; isi .env di folder kerja atau environment OS")
	}
	u, err := url.Parse(value)
	if err != nil || (u.Scheme != "postgres" && u.Scheme != "postgresql") || u.Hostname() == "" || u.User == nil || u.User.Username() == "" || u.Path == "" || u.Path == "/" {
		return "", errors.New("DATABASE_URL wajib berupa connection string PostgreSQL lengkap")
	}
	params := u.Query()
	switch params.Get("sslmode") {
	case "":
		params.Set("sslmode", "require")
	case "require", "verify-ca", "verify-full":
	default:
		return "", errors.New("DATABASE_URL harus memakai sslmode=require, verify-ca, atau verify-full")
	}
	// Parameter ssl=true dari beberapa contoh tidak didukung lib/pq.
	params.Del("ssl")
	u.RawQuery = params.Encode()
	return u.String(), nil
}
