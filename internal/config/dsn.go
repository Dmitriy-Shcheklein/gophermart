package config

import (
	"flag"
	"log"
	"os"
)

type DSN struct {
	Value string
}

func NewDSN() *DSN {
	dsn := &DSN{}

	if dbUri := os.Getenv("DATABASE_URI"); dbUri != "" {
		if err := dsn.Set(dbUri); err != nil {
			log.Fatalf("error while set DATABASE_URI env: %s", err)
		}
	}
	flag.Var(dsn, "d", "Database uri")

	return dsn
}

func (a *DSN) ApplyEnv() {
	dsn, ok := os.LookupEnv("DATABASE_URI")
	if !ok || dsn == "" {
		return
	}
	if err := a.Set(dsn); err != nil {
		log.Fatalf("error while set DATABASE_URI env: %s", err)
	}
}

func (a *DSN) String() string {
	return a.Value
}

func (a *DSN) Set(s string) error {
	a.Value = s
	return nil
}
