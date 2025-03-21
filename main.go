package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"payments/bank"
	"time"
)

func main() {
	bank := bank.NewBank()
	bank.RegisterHandlers()

	srv := NewServer()

	stop := make(chan os.Signal, 1)

	go func() {
		httpError := srv.ListenAndServe()
		if httpError != nil {
			panic(fmt.Errorf("error starting http server %v ", httpError))
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		fmt.Println("error shutting down the server")
	}
}

func NewServer() *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}
}
