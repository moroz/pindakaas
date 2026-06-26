package config

import (
	"log"
	"os"
	"strconv"
)

func MustGetenv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		log.Fatalf("FATAL: Environment variable %s is not set!", name)
	}
	return val
}

func GetenvWithDefault(name, defaultValue string) string {
	val := os.Getenv(name)
	if val == "" {
		return defaultValue
	}
	return val
}

func MustParsePortNumber(val string) uint16 {
	parsed, err := strconv.ParseInt(val, 10, 16)
	if err != nil {
		log.Fatalf(`FATAL: Failed to parse port number from string "%c"`, val)
	}
	return uint16(parsed)
}

var SSHServerKeyPath = MustGetenv("SSH_SERVER_KEY_PATH")

var SSHPort = MustParsePortNumber(GetenvWithDefault("SSH_PORT", "2137"))
var HTTPPort = MustParsePortNumber(GetenvWithDefault("HTTP_PORT", "8080"))
var HTTPSPort = MustParsePortNumber(GetenvWithDefault("HTTPS_PORT", "8081"))
