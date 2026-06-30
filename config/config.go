package config

import (
	"crypto/sha512"
	"encoding/base64"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/bincyber/go-sqlcrypter"
	"github.com/moroz/pindakaas/internal/crypto"
	"golang.org/x/crypto/hkdf"
)

func MustGetenv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		log.Fatalf("FATAL: Environment variable %s is not set!", name)
	}
	return val
}

func MustGetenvBase64(name string) []byte {
	val := MustGetenv(name)
	binary, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		log.Fatalf("FATAL: Failed to decode environment variable %s as Base64!", name)
	}
	return binary
}

func GetenvWithDefault(name, defaultValue string) string {
	val := os.Getenv(name)
	if val == "" {
		return defaultValue
	}
	return val
}

func MustParsePortNumber(val string) uint16 {
	parsed, err := strconv.ParseUint(val, 10, 16)
	if err != nil {
		log.Fatalf(`FATAL: Failed to parse port number from string "%v"`, val)
	}
	return uint16(parsed)
}

func MustDeriveKey(base []byte, info string, lengthInBytes int) []byte {
	kdf := hkdf.New(sha512.New, base, nil, []byte(info))
	buf := make([]byte, lengthInBytes)
	if _, err := io.ReadFull(kdf, buf); err != nil {
		log.Fatalf("Failed to derive key (info: %s): %s", info, err)
	}
	return buf
}

func FormatHostPort(port uint16) string {
	return net.JoinHostPort("0.0.0.0", strconv.Itoa(int(port)))
}

func RequireInProduction(name string, defaultValue string) string {
	if IsProd {
		return MustGetenv(name)
	}
	return GetenvWithDefault(name, defaultValue)
}

func init() {
	crypterer, err := crypto.NewEncryptionProvider(DatabaseEncryptionKey, nil)
	if err != nil {
		log.Fatalf("Failed to initialize database encryption provider: %s", err)
	}
	sqlcrypter.Init(crypterer)
}

var IsProd = os.Getenv("GO_ENV") == "prod"

var SSHServerKeyPath = MustGetenv("SSH_SERVER_KEY_FILE")
var SSHPort = MustParsePortNumber(GetenvWithDefault("SSH_PORT", "42069"))

var HTTPPort = MustParsePortNumber(GetenvWithDefault("HTTP_PORT", "8080"))
var HTTPSPort = MustParsePortNumber(GetenvWithDefault("HTTPS_PORT", "8081"))

// DisableHTTP2, when set via DISABLE_HTTP2="true", prevents the HTTPS server
// from negotiating "h2" via ALPN. HTTP/2 connections are not hijackable, so
// WebSocket upgrades require HTTP/1.1; disable HTTP/2 if you proxy WebSockets.
var DisableHTTP2 = GetenvWithDefault("DISABLE_HTTP2", "") == "true"

var TLSCertFile = MustGetenv("TLS_CERT_FILE")
var TLSKeyFile = MustGetenv("TLS_KEY_FILE")

var BaseDomain = GetenvWithDefault("BASE_DOMAIN", "")
var DatabaseUrl = MustGetenv("DATABASE_URL")
var SecretKeyBase = MustGetenvBase64("SECRET_KEY_BASE")

var SessionKey = MustDeriveKey(SecretKeyBase, "Sessions", 32)
var DatabaseEncryptionKey = MustDeriveKey(SecretKeyBase, "ColumnLevelEncryption", 32)

var GoogleClientId = RequireInProduction("GOOGLE_CLIENT_ID", "")
var GoogleClientSecret = RequireInProduction("GOOGLE_CLIENT_SECRET", "")
var PublicUrl = RequireInProduction("PUBLIC_URL", "http://localhost:8080")
