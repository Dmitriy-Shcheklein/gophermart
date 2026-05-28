package config

import (
	"errors"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

// SrvAddr структура для хранения конфигурации адреса запуска сервиса
type SrvAddr struct {
	host string
	port int
}

// NewSrvAddr конструктор для конфигурации адреса запуска сервиса
func NewSrvAddr() *SrvAddr {
	port := 8080
	srvAddr := &SrvAddr{host: "localhost", port: port}
	flag.Var(srvAddr, "a", "server address host:port")
	return srvAddr
}

// ApplyEnv метод для применения данных из .ENV
func (a *SrvAddr) ApplyEnv() {
	serverAddress, ok := os.LookupEnv("RUN_ADDRESS")
	if !ok || serverAddress == "" {
		return
	}
	if err := a.Set(serverAddress); err != nil {
		log.Fatalf("error while set RUN_ADDRESS env: %s", err)
	}
}

// String реализация для соответствия интерфейсу flag.Var
func (a *SrvAddr) String() string {
	return a.host + ":" + strconv.Itoa(a.port)
}

// Set реализация для соответствия интерфейсу flag.Var
func (a *SrvAddr) Set(s string) error {
	hp := strings.Split(s, ":")
	validLength := 2
	if len(hp) != validLength {
		return errors.New("need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	a.host = hp[0]
	a.port = port
	return nil
}
