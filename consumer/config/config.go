package config

import (
	"fmt"
	"os"
)

const (
	defaultServHost   = "web"
	defaultQuHost     = "consumer"
	defaultQuPort     = "5672"
	defaultQuUser     = "guest"
	defaultQuPass     = "guest"
	defaultPostgreSQL = "10.80.32.66"
	defaultDBPort     = "5432"
	defaultDSSLMode   = "disable"
	defaultDBuser     = "user"
	defaultDBpass     = "pass"
	defaultDBname     = "db"
)

type Config struct {
	ServerAdress   string
	QueueAdress    string
	PostgresAdress string
}

func Init() *Config {
	var Server string
	var ok bool
	if Server, ok = os.LookupEnv("SERVER_ADDR"); !ok {
		Server = defaultServHost
	}

	var QuHost string
	if QuHost, ok = os.LookupEnv("QUEUE_HOST"); !ok {
		QuHost = defaultQuHost
	}

	var QuPort string
	if QuPort, ok = os.LookupEnv("QUEUE_PORT"); !ok {
		QuPort = defaultQuPort
	}

	var QuUser string
	if QuUser, ok = os.LookupEnv("QUEUE_USER"); !ok {
		QuUser = defaultQuUser
	}

	var QuPass string
	if QuPass, ok = os.LookupEnv("QUEUE_PASS"); !ok {
		QuPass = defaultQuPass
	}

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

	queueAdress := fmt.Sprintf("amqp://%s:%s@%s:%s/", QuUser, QuPass, QuHost, QuPort)
	postgresAddr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		DbHost, DbPort, DbUser, DbPassword, DbName, DbSslMode)

	return &Config{Server, queueAdress, postgresAddr}
}
