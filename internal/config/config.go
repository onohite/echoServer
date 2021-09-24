package config

import (
	"fmt"
	"os"
)

const (
	defaultHTTPPort   = "80"
	defaultHost       = ""
	Prod              = "prod"
	Dev               = "dev"
	defaultPostgreSQL = "10.80.32.66"
	defaultDBPort     = "5432"
	defaultDSSLMode   = "disable"
	defaultDBuser     = "user"
	defaultDBpass     = "pass"
	defaultDBname     = "db"
)

type Config struct {
	Host           string
	Port           string
	ServerMode     string
	PostgresAdress string
}

func Init() *Config {
	serverMode, ok := os.LookupEnv("SERVER_MODE")
	if !ok {
		serverMode = Dev
	}

	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = defaultHost
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultHTTPPort
	}

	var DbUser string
	if DbUser, ok = os.LookupEnv("DBUSER"); !ok {
		DbUser = "123"
	}

	var DbPassword string
	if DbPassword, ok = os.LookupEnv("DBPASSWORD"); !ok {
		DbPassword = defaultDBpass
	}

	var DbName string
	if DbName, ok = os.LookupEnv("DBNAME"); !ok {
		DbName = defaultDBname
	}

	var DbHost string
	if DbHost, ok = os.LookupEnv("DBHOST"); !ok {
		DbHost = defaultPostgreSQL
	}

	var DbPort string
	if DbPort, ok = os.LookupEnv("DBPORT"); !ok {
		DbPort = defaultDBPort
	}

	var DbSslMode string
	if DbSslMode, ok = os.LookupEnv("DBSSLMODE"); !ok {
		DbSslMode = defaultDSSLMode
	}

	postgresAddr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		DbHost, DbPort, DbUser, DbPassword, DbName, DbSslMode)

	return &Config{host, port, serverMode, postgresAddr}
}
