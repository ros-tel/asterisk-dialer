package main

import (
	"flag"
	"log"
	"os"

	"asterisk-dialer/api"
	"asterisk-dialer/asterisk"
	"asterisk-dialer/config"

	"github.com/gin-gonic/gin"
)

var (
	// config the settings variable
	cnf = &conf{}

	config_file    = flag.String("config", "", "Usage: -config=<config_file>")
	templates_path = flag.String("templates", "templates", "Usage: -templates=<templates_path>")
	debug          = flag.Bool("debug", false, "Print debug information on stderr")
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	flag.Parse()

	getConfig()
	config.LoadTemplate(cnf.Config, *templates_path+string(os.PathSeparator)+"*.tpl")

	if *debug {
		log.Printf("CONFIG: %+v", cnf)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	ami := asterisk.Connect(cnf.Asterisk)

	r := gin.Default()
	r.Use(
		asterisk.Inject(ami),
		config.Inject(cnf.Config),
	)

	api.InitRoutes(r)

	r.Run(cnf.Listen)
}
