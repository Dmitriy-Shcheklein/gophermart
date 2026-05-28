package config

import (
	"flag"
	"log"
	"os"
)

// DSN структруа для хранения конфигурации адреса подключения к БД
type DSN struct {
	value string
}

// NewDSN конструктор для конфигурации адреса подключения к БД
func NewDSN() *DSN {
	dsn := &DSN{}
	flag.Var(dsn, "d", "Database uri")
	return dsn
}

// ApplyEnv метод для применения данных из .ENV
func (a *DSN) ApplyEnv() {
	dsn, ok := os.LookupEnv("DATABASE_URI")
	if !ok || dsn == "" {
		return
	}
	if err := a.Set(dsn); err != nil {
		log.Fatalf("error while set DATABASE_URI env: %s", err)
	}
}

// String реализация для соответствия интерфейсу flag.Var
func (a *DSN) String() string {
	return a.value
}

// Set реализация для соответствия интерфейсу flag.Var
func (a *DSN) Set(s string) error {
	a.value = s
	return nil
}
