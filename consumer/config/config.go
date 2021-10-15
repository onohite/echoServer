package config

import (
	"fmt"
	"os"
)

const (
	defaultServHost  = "web"
	defaultQuHost    = "consumer"
	defaultQuPort    = "5672"
	defaultQuUser    = "guest"
	defaultQuPass    = "guest"
	defaultCacheHost = "web-cache"
	defaultCachePort = "6380"
)

type Config struct {
	ServerAdress string
	QueueAdress  string
	CacheAdress  string
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

	var CHost string
	if CHost, ok = os.LookupEnv("CACHE_HOST"); !ok {
		CHost = defaultCacheHost
	}

	var CPort string
	if CPort, ok = os.LookupEnv("CACHE_PORT"); !ok {
		CPort = defaultCachePort
	}

	queueAdress := fmt.Sprintf("amqp://%s:%s@%s:%s/", QuUser, QuPass, QuHost, QuPort)
	cacheAdress := fmt.Sprintf("%s:%s", CHost, CPort)

	return &Config{Server, queueAdress, cacheAdress}
}
