package main

import (
	"asterisk-dialer/asterisk"
	"asterisk-dialer/config"
	"flag"
	"log"
	"os"

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
	r.Use(injectAmi(ami), config.Inject(cnf.Config))

	initRoutes(r)

	r.Run(cnf.Listen)
}
