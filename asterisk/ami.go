package asterisk

import (
	"log"

	"github.com/heltonmarx/goami/ami"
)

type (
	Ami struct {
		Username string `yaml:"user"`
		Password string `yaml:"password"`
		Addr     string `yaml:"addr"`
	}
)

func Connect(c Ami) *ami.Socket {
	socket, err := ami.NewSocket(c.Addr)
	if err != nil {
		log.Fatalf("socket error: %v\n", err)
	}
	if _, err := ami.Connect(socket); err != nil {
		log.Fatalf("connect error: %v\n", err)
	}

	//Login
	uuid, err := ami.GetUUID()
	if err != nil {
		log.Fatalf("Get UUID: %v\n", err)
	}
	if err := ami.Login(socket, c.Username, c.Password, "Off", uuid); err != nil {
		log.Fatalf("login error: %v\n", err)
	}
	log.Printf("login ok!\n")

	return socket
}
