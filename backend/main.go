package main

import (
	"fmt"
	"net/http"
	"time"

	"bank-service/config"
	"bank-service/src"
	"bank-service/src/libs/env"

	"bank-service/src/libs/logger"
)

func main() {
	config.SetupCommonDependencies()
	defer config.TearDownCommonDependencies()
	handler := src.SetupHandler()

	host := fmt.Sprint(":", env.BankServiceRestPort)
	srv := &http.Server{
		Handler:      *handler,
		Addr:         host,
		WriteTimeout: 20 * time.Second,
		ReadTimeout:  20 * time.Second,
	}
	go srv.ListenAndServe()
	logger.GetInstance().Info("Server Listening on ", host)

	select {} //Infinite waiting
}
