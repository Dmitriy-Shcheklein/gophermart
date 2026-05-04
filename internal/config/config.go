package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
)

var errConfig = errors.New("config init error")

type Config struct {
	Port              int
	Host              string
	DbDSN             string
	AccrualSystemAddr string
}

type DbDSN struct {
	Value   string
	IsValid bool
}

func New() (*Config, error) {
	srvAddr := NewSrvAddr()
	dsn := NewDSN()
	accrual := NewAccrual()
	flag.Parse()
	srvAddr.ApplyEnv()
	dsn.ApplyEnv()
	accrual.ApplyEnv()

	if dsn.Value == "" {
		return nil, fmt.Errorf("invalid db dsn: %w", errConfig)
	}
	if accrual.Value == "" {
		return nil, fmt.Errorf("invalid accrual system address: %w", errConfig)
	}

	return &Config{
		Host:              srvAddr.Host,
		Port:              srvAddr.Port,
		AccrualSystemAddr: accrual.Value,
	}, nil
}

func (c *Config) GetSrvAddr() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

func (c *Config) DbDsn() string {
	return c.DbDSN
}
