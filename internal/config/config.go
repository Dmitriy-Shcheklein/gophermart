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
		host:              srvAddr.Host,
		port:              srvAddr.Port,
		accrualSystemAddr: accrual.Value,
		dbDSN:             dsn.Value,
	}, nil
}

func (c *Config) GetSrvAddr() string {
	return c.host + ":" + strconv.Itoa(c.port)
}

func (c *Config) DbDsn() string {
	return c.dbDSN
}

func (c *Config) GetAccrualSrvAddr() string {
	return c.accrualSystemAddr
}
