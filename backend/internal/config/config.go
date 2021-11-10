package config

import (
	"fmt"
	"log"
	"os"
)

const (
	defaultDNS        = ""
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
	defaultQuHost     = "10.80.32.66"
	defaultQuPort     = "5672"
	defaultQuUser     = "guest"
	defaultQuPass     = "guest"
	defaultCacheHost  = "web-cache"
	defaultCachePort  = "6380"
)

type Config struct {
	Host           string
	Port           string
	ServerMode     string
	PostgresAdress string
	QueueAdress    string
	CacheAdress    string
	Dns            string
	AuthType       AuthType
}

type AuthType struct {
	VKconfig      AuthConfig
	DiscordConfig AuthConfig
	GoogleConfig  AuthConfig
}

type AuthConfig struct {
	ClientID     string
	ClientSecret string
}

func Init() *Config {
	dns, ok := os.LookupEnv("DNS_SERVER")
	if !ok {
		dns = defaultDNS
	}

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

	var CHost string
	if CHost, ok = os.LookupEnv("CACHE_HOST"); !ok {
		CHost = defaultCacheHost
	}

	var CPort string
	if CPort, ok = os.LookupEnv("CACHE_PORT"); !ok {
		CPort = defaultCachePort
	}

	var VKClientID string
	if VKClientID, ok = os.LookupEnv("VK_CLIENT_ID"); !ok {
		log.Fatal("empty VK_CLIENT_ID")
	}

	var VKClientSecret string
	if VKClientSecret, ok = os.LookupEnv("VK_CLIENT_SECRET"); !ok {
		log.Fatal("empty VK_CLIENT_SECRET")
	}

	vkCFG := AuthConfig{
		ClientID:     VKClientID,
		ClientSecret: VKClientSecret,
	}

	var GoogleClientID string
	if GoogleClientID, ok = os.LookupEnv("GOOGLE_CLIENT_ID"); !ok {
		log.Fatal("empty GOOGLE_CLIENT_ID")
	}

	var GoogleClientSecret string
	if GoogleClientSecret, ok = os.LookupEnv("GOOGLE_CLIENT_SECRET"); !ok {
		log.Fatal("empty GOOGLE_CLIENT_SECRET")
	}

	googleCFG := AuthConfig{
		ClientID:     GoogleClientID,
		ClientSecret: GoogleClientSecret,
	}

	var DiscordClientID string
	if DiscordClientID, ok = os.LookupEnv("DISCORD_CLIENT_ID"); !ok {
		log.Fatal("empty DISCORD_CLIENT_ID")
	}

	var DiscordClientSecret string
	if DiscordClientSecret, ok = os.LookupEnv("DISCORD_CLIENT_SECRET"); !ok {
		log.Fatal("empty DISCORD_CLIENT_SECRET")
	}

	discordCFG := AuthConfig{
		ClientID:     DiscordClientID,
		ClientSecret: DiscordClientSecret,
	}
	authType := AuthType{vkCFG, discordCFG, googleCFG}

	postgresAddr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		DbHost, DbPort, DbUser, DbPassword, DbName, DbSslMode)
	queueAddress := fmt.Sprintf("amqp://%s:%s@%s:%s/", QuUser, QuPass, QuHost, QuPort)
	cacheAdress := fmt.Sprintf("%s:%s", CHost, CPort)

	return &Config{host, port, serverMode, postgresAddr, queueAddress, cacheAdress, dns, authType}
}
