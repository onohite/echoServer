package config

import (
	"fmt"
	"os"
)

const (
	defaultServHost = "web"
	defaultQuHost   = "localhost"
	defaultQuPort   = "5672"
	defaultQuUser   = "guest"
	defaultQuPass   = "guest"
)

type Config struct {
	ServerAdress string
	QueueAdress  string
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

	queueAddress := fmt.Sprintf("amqp://%s:%s@%s:%s/", QuUser, QuPass, QuHost, QuPort)

	return &Config{Server, queueAddress}
}
