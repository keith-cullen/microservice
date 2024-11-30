package config

import "log"

const (
	StoreDriverNameKey    = "StoreDriverName"
	StoreDatabaseFileKey  = "StoreDatabaseFile"
	StaticDirKey          = "StaticDirKey"
	StaticHtmlFileKey     = "StaticHtmlFile"
	CertKey               = "CertKey"
	PrivkeyKey            = "PrivKey"
	HttpAddrKey           = "HttpAddr"
	HttpsAddrKey          = "HttpsAddr"
	HttpCorsOriginKey     = "HttpCorsOrigin"
	HttpsCorsOriginKey    = "HttpsCorsOrigin"
	HttpTimeoutKey        = "HttpTimeout"
	HttpMaxHeaderBytesKey = "HttpMaxHeaderBytes"
	HttpReqsPerSecKey     = "HttpReqsPerSec"
)

var (
	Data = map[string]string{
		StoreDriverNameKey:    "sqlite3",
		StoreDatabaseFileKey:  "file:data/store.db?_fk=1",
		StaticDirKey:          "www",
		StaticHtmlFileKey:     "index.html",
		CertKey:               "./server_cert.pem",
		PrivkeyKey:            "./server_privkey.pem",
		HttpAddrKey:           "0.0.0.0:8080",
		HttpsAddrKey:          "0.0.0.0:4443",
		HttpCorsOriginKey:     "https://localhost",
		HttpsCorsOriginKey:    "https://localhost",
		HttpTimeoutKey:        "10",
		HttpMaxHeaderBytesKey: "4096",
		HttpReqsPerSecKey:     "10",
	}
)

func Get(key string) string {
	val, ok := Data[key]
	if !ok {
		log.Printf("configuration data missing for key: '%s'", key)
	}
	return val
}
