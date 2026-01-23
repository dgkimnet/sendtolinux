package config

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Bind        string
	Port        int
	Dir         string
	MaxUploadMB int
}

func ParseFlags() Config {
	defaults := defaultConfigFromEnv()

	bind := flag.String("bind", defaults.Bind, "Bind address for HTTP server")
	port := flag.Int("port", defaults.Port, "Port for HTTP server (0 for random)")
	dir := flag.String("dir", defaults.Dir, "Directory to store received files (empty for default)")
	maxUpload := flag.Int("max-upload-mb", defaults.MaxUploadMB, "Maximum upload size in MB")

	flag.Parse()

	return Config{
		Bind:        strings.TrimSpace(*bind),
		Port:        *port,
		Dir:         strings.TrimSpace(*dir),
		MaxUploadMB: *maxUpload,
	}
}

func defaultConfigFromEnv() Config {
	cfg := Config{
		Bind:        "0.0.0.0",
		Port:        8000,
		Dir:         "",
		MaxUploadMB: 100,
	}

	if val := strings.TrimSpace(os.Getenv("STL_BIND")); val != "" {
		cfg.Bind = val
	}
	if val := strings.TrimSpace(os.Getenv("STL_PORT")); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			cfg.Port = parsed
		}
	}
	if val := strings.TrimSpace(os.Getenv("STL_DIR")); val != "" {
		cfg.Dir = val
	}
	if val := strings.TrimSpace(os.Getenv("STL_MAX_UPLOAD_MB")); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			cfg.MaxUploadMB = parsed
		}
	}

	return cfg
}
