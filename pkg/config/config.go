package config

import "os"

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
	SslMode  string
	ApiKey   string
}

func GetConfig() Config {
	return Config{
		Host:     os.Getenv("host"),
		Port:     os.Getenv("port"),
		User:     os.Getenv("user"),
		Password: os.Getenv("password"),
		Dbname:   os.Getenv("db"),
		SslMode:  os.Getenv("ssl"),
		ApiKey:   os.Getenv("key"),
	}
}

//$env:host = '185.200.241.2'
//$env:port = '5432'
//$env:user = 'vk'
//$env:password = 'vkdb'
//$env:db = 'vkdb'
//$env:ssl = 'disable'
//$env:key = '5995894659:AAG82B6kbmD17TmmPcKT5Zzqz4S6LpQolYQ'
