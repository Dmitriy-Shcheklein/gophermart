package config

import (
	"flag"
	"log"
	"os"
)

type Accrual struct {
	Value string
}

func NewAccrual() *Accrual {
	accrual := &Accrual{}
	flag.Var(accrual, "r", "accrual system address")
	return accrual
}

func (a *Accrual) ApplyEnv() {
	accrual, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	if !ok || accrual == "" {
		return
	}
	if err := a.Set(accrual); err != nil {
		log.Fatalf("error while set ACCRUAL_SYSTEM_ADDRESS env: %s", err)
	}
}

func (a *Accrual) String() string {
	return a.Value
}

func (a *Accrual) Set(s string) error {
	a.Value = s
	return nil
}
