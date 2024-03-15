package config

import (
	"log"
	"os"
	"strconv"
)

func GetCertDir() string {
	return getEnv("CERT_DIR")
}

func GetPaymentURL() string {
	return getEnv("PAYMENT_URL")
}

func GetEnv() string {
	return getEnv("ENV")
}

func GetDatabaseURL() string {
	return getEnv("DATABASE_URL")
}

func GetApplicationPort() int {
	portString := getEnv("PORT")
	port, err := strconv.Atoi(portString)

	if err != nil {
		log.Fatalf("port %s is invalid", portString)
	}

	return port
}

func getEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("%s env is missing", key)
	}

	return v
}
