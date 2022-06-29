package asterisk

import (
	"log"
	"time"

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

	_, err = ami.Connect(socket)
	if err != nil {
		log.Fatalf("connect error: %v\n", err)
	}

	//Login
	uuid, err := ami.GetUUID()
	if err != nil {
		log.Fatalf("[FATAL] ami.GetUUID: %v\n", err)
	}

	err = ami.Login(socket, c.Username, c.Password, "Off", uuid)
	if err != nil {
		log.Fatalf("[FATAL] ami.Login: %v\n", err)
	}

	go amiHeartbeat(socket)

	log.Printf("login ok!\n")

	return socket
}

// Раз в минуту проверяем что подключение живо
func amiHeartbeat(socket *ami.Socket) {
	ticker := time.NewTicker(time.Minute)

	for {
		<-ticker.C

		uuid, err := ami.GetUUID()
		if err != nil {
			log.Fatalf("[FATAL] ami.GetUUID: %v\n", err)
		}

		err = ami.Ping(socket, uuid)
		if err != nil {
			log.Fatalf("[FATAL] ami.Ping: %v\n", err)
		}
	}
}
