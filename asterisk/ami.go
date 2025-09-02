package asterisk

import (
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	socket, err := ami.NewSocket(ctx, c.Addr)
	if err != nil {
		log.Fatalf("socket error: %v\n", err)
	}

	_, err = ami.Connect(ctx, socket)
	if err != nil {
		log.Fatalf("connect error: %v\n", err)
	}

	//Login
	uuid, err := ami.GetUUID()
	if err != nil {
		log.Fatalf("[FATAL] ami.GetUUID: %v\n", err)
	}

	err = ami.Login(ctx, socket, c.Username, c.Password, "Off", uuid)
	if err != nil {
		log.Fatalf("[FATAL] ami.Login: %v\n", err)
	}

	go amiHeartbeat(socket)

	log.Printf("login ok!\n")

	return socket
}

// Раз в минуту проверяем что подключение живо
func amiHeartbeat(socket *ami.Socket) {
	ticker := time.NewTicker(30 * time.Second)

	for {
		<-ticker.C

		uuid, err := ami.GetUUID()
		if err != nil {
			log.Fatalf("[FATAL] ami.GetUUID: %v\n", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = ami.Ping(ctx, socket, uuid)
		if err != nil {
			log.Fatalf("[FATAL] ami.Ping: %v\n", err)
		}
		cancel()
	}
}
