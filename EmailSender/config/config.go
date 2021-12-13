package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	defaultPostgreSQL = "10.80.32.66"
	defaultDBPort     = "5432"
	defaultDSSLMode   = "disable"
	defaultDBuser     = "user"
	defaultDBpass     = "pass"
	defaultDBname     = "db"
	defaultTickerTime = "15"
)

type Config struct {
	PostgresAdress string
	TickerTime     int
}

func Init() *Config {
	var ok bool
	var DbUser string
	if DbUser, ok = os.LookupEnv("DBUSER"); !ok {
		DbUser = defaultDBuser
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

	var Ticker string
	if Ticker, ok = os.LookupEnv("Ticker"); !ok {
		Ticker = defaultTickerTime
	}

	resultTicker, _ := strconv.Atoi(Ticker)
	postgresAddr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		DbHost, DbPort, DbUser, DbPassword, DbName, DbSslMode)

	return &Config{postgresAddr, resultTicker}
}
