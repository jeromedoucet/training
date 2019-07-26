package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jeromedoucet/training/configuration"
	"github.com/jeromedoucet/training/controller"
)

var dbName string
var dbUser string
var dbPassword string
var dbHost string
var dbPort uint
var dbSslMode bool

var httpPort uint
var jwtSecret string

func main() {

	flag.StringVar(&dbName, "dbName", "", "Name of the database.")
	flag.StringVar(&dbUser, "dbUser", "postgres", "User of the database.")
	flag.StringVar(&dbPassword, "dbPassword", "", "Password of the database.")
	flag.StringVar(&dbHost, "dbHost", "", "Host of the database.")
	flag.UintVar(&dbPort, "dbPort", 5432, "Port of the database.")
	flag.BoolVar(&dbSslMode, "dbSslModeOn", false, "Ssl mode enabled.")

	flag.UintVar(&httpPort, "httpPort", 8080, "Port of the http server.")
	flag.StringVar(&jwtSecret, "jwtSecret", configuration.RandStringBytes(32), "Secret to use for jwt secret.")

	flag.Parse()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	conf := &configuration.GlobalConf{
		DbName:        dbName,
		User:          dbUser,
		Password:      dbPassword,
		Host:          dbHost,
		Port:          dbPort,
		SslMode:       dbSslMode,
		JwtSecret:     jwtSecret,
		JwtExpiration: 720 * time.Hour,
	}
	h := controller.InitRoutes(conf)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: h,
	}

	go func() {
		log.Printf("INFO >> Listening on http://0.0.0.0:%d\n", httpPort)
		if err := s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop
}
