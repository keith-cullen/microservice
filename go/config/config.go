package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	DatabaseDriverNameKey = "DatabaseDriverName"
	DatabaseFileKey       = "DatabaseFile"
	CertKey               = "Cert"
	PrivkeyKey            = "Privkey"
	HttpAddrKey           = "HttpAddr"
	HttpsAddrKey          = "HttpsAddr"
	HttpCorsOriginKey     = "HttpCorsOrigin"
	HttpsCorsOriginKey    = "HttpsCorsOrigin"
	HttpTimeoutKey        = "HttpTimeout"
	HttpMaxHeaderBytesKey = "HttpMaxHeaderBytes"
	HttpReqsPerSecKey     = "HttpReqsPerSec"
)

var (
	Data = map[string]string{}
)

func Open(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read configuration file: %s: %w", filename, err)
	}
	if err = yaml.Unmarshal([]byte(content), &Data); err != nil {
		return fmt.Errorf("failed to parse configuration file: %s: %w", filename, err)
	}
	return nil
}

func Get(key string) string {
	val, ok := Data[key]
	if !ok {
		log.Printf("configuration data missing for key: '%s'", key)
	}
	return val
}