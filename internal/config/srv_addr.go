package config

import (
	"errors"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

type SrvAddr struct {
	Host string
	Port int
}

func NewSrvAddr() *SrvAddr {
	port := 8080
	srvAddr := &SrvAddr{Host: "localhost", Port: port}
	flag.Var(srvAddr, "a", "server address host:port")
	return srvAddr
}

func (a *SrvAddr) ApplyEnv() {
	serverAddress, ok := os.LookupEnv("RUN_ADDRESS")
	if !ok || serverAddress == "" {
		return
	}
	if err := a.Set(serverAddress); err != nil {
		log.Fatalf("error while set RUN_ADDRESS env: %s", err)
	}
}

func (a *SrvAddr) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

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
	a.Host = hp[0]
	a.Port = port
	return nil
}
