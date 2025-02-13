package config

type Configuration interface {
	Init(config interface{}) error
}
