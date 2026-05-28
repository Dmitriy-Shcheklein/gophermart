// Package config конфигурация приложения
package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
)

var errConfig = errors.New("config init error")

type Config struct {
	port              int
	host              string
	dbDSN             string
	accrualSystemAddr string
}

// New конструктор для конфигурации
func New() (*Config, error) {
	srvAddr := NewSrvAddr()
	dsn := NewDSN()
	accrual := NewAccrual()
	flag.Parse()
	srvAddr.ApplyEnv()
	dsn.ApplyEnv()
	accrual.ApplyEnv()

	if dsn.value == "" {
		return nil, fmt.Errorf("invalid db dsn: %w", errConfig)
	}
	if accrual.value == "" {
		return nil, fmt.Errorf("invalid accrual system address: %w", errConfig)
	}

	return &Config{
		host:              srvAddr.host,
		port:              srvAddr.port,
		accrualSystemAddr: accrual.value,
		dbDSN:             dsn.value,
	}, nil
}

// GetSrvAddr Метод для получения адреса для запуска сервиса
func (c *Config) GetSrvAddr() string {
	return c.host + ":" + strconv.Itoa(c.port)
}

// DbDsn Метод для получения строки подключения к БД
func (c *Config) DbDsn() string {
	return c.dbDSN
}

// GetAccrualSrvAddr Метода для получения адреса системы расчетов баллов
func (c *Config) GetAccrualSrvAddr() string {
	return c.accrualSystemAddr
}
