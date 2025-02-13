package config

type AutoConfiguration interface {
	Init(config interface{}) error
}
