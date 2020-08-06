package main

import (
	"os"

	log "github.com/towl/logger"
)

// Config is the configuration structure
type Config struct {
	Host       string
	Port       string
	WorkingDir string
	AuditDir   string
}

var config = &Config{}
var logger = log.GetLoggerFromEnv("", false)

func init() {
	config.loadEnv()
	logger.Info("Server config loaded successfully.")
}

func (c *Config) loadEnv() {
	c.Host = os.Getenv("HOST")
	c.Port = os.Getenv("PORT")
	c.WorkingDir = os.Getenv("WORKING_DIR")
	c.AuditDir = os.Getenv("AUDIT_DIR")
}
