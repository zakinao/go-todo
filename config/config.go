package config

type ConfigList struct {
	Port string
	SQLDriver string
	DbName string
	LogFile string
}

var Config ConfigList
