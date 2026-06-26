package config

import (
	"log"
	"net"
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
	parsed, err := strconv.ParseInt(val, 10, 64)
	if err != nil || parsed > 65535 {
		log.Fatalf(`FATAL: Failed to parse port number from string "%v"`, val)
	}
	return uint16(parsed)
}

func FormatHostPort(port uint16) string {
	return net.JoinHostPort("0.0.0.0", strconv.Itoa(int(port)))
}

var SSHServerKeyPath = MustGetenv("SSH_SERVER_KEY_FILE")

var SSHPort = MustParsePortNumber(GetenvWithDefault("SSH_PORT", "2137"))
var HTTPPort = MustParsePortNumber(GetenvWithDefault("HTTP_PORT", "8080"))
var HTTPSPort = MustParsePortNumber(GetenvWithDefault("HTTPS_PORT", "8081"))

var TLSCertFile = MustGetenv("TLS_CERT_FILE")
var TLSKeyFile = MustGetenv("TLS_KEY_FILE")

var BaseDomain = GetenvWithDefault("BASE_DOMAIN", "")
var DatabaseUrl = MustGetenv("DATABASE_URL")
