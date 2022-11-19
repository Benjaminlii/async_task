package config

type RocketMQConfig struct {
	NameServers []string
	Topic       string
}

type RedisConfig struct {
	Address  string
	Password string
}
