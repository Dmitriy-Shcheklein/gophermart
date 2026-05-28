package config

import (
	"flag"
	"log"
	"os"
)

// Accrual структура для хранения конфигурации адреса системы расчета баллов
type Accrual struct {
	value string
}

// NewAccrual конструктор для конфигурации адреса системы расчета баллов
func NewAccrual() *Accrual {
	accrual := &Accrual{}
	flag.Var(accrual, "r", "accrual system address")
	return accrual
}

// ApplyEnv метод для применения данных из .ENV
func (a *Accrual) ApplyEnv() {
	accrual, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	if !ok || accrual == "" {
		return
	}
	if err := a.Set(accrual); err != nil {
		log.Fatalf("error while set ACCRUAL_SYSTEM_ADDRESS env: %s", err)
	}
}

// String реализация для соответствия интерфейсу flag.Var
func (a *Accrual) String() string {
	return a.value
}

// Set реализация для соответствия интерфейсу flag.Var
func (a *Accrual) Set(s string) error {
	a.value = s
	return nil
}
