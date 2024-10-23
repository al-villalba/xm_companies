package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"producer/common"

	"github.com/sirupsen/logrus"
)

func main() {
	// define environment
	e := common.GetEnv()
	// uncomment to generate and show a jwt token at start
	// tokenString, _ := common.GenerateJwtTokenString("appuser", common.ROLE_WRITER)
	// logrus.Info(tokenString)

	// create the router
	r := SetupRouter()

	// init db (panic if it fails)
	db := common.GetDatabase()
	defer db.Close()

	// init kafka producer
	p, err := common.GetProducer()
	if err != nil {
		logrus.WithError(err).Error("Failed to create kafka producer")
		os.Exit(1)
	}
	defer p.Close()

	// start the server
	go func() {
		fmt.Println("Starting http server on :" + e.Addr)
		if err := http.ListenAndServe(":"+e.Addr, r); err != nil {
			fmt.Printf("Server failed to start: %v\n", err)
			os.Exit(1)
		}
	}()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals

	os.Exit(0)
}
