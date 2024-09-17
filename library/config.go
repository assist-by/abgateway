package library

import "os"

type Config struct {
	KafkaBroker         string
	Host                string
	Port                string
	RegistrationTopic   string
	ServiceDiscoveryURL string
}

func Load() *Config {
	return &Config{
		KafkaBroker:         getEnv("KAFKA_BROKER", "kafka:9092"),
		Host:                getEnv("HOST", "autro-api-gateway"),
		Port:                getEnv("PORT", "50050"),
		RegistrationTopic:   getEnv("REGISTRATION_TOPIC", "service-registration"),
		ServiceDiscoveryURL: getEnv("SERVICE_DISCOVERY_URL", "http://abready:8500"),
	}
}

func getEnv(key, temp string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return temp
}
