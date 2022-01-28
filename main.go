package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/kevin/working/handlers"
	"github.com/nicholasjackson/env"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var bindAddress = env.String("BIND_ADDRESS",false,":9090","Bind address for the server")

func main() {
	env.Parse()

	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	//create product handler
	ph := handlers.NewProducts(l)

	//create a new serve mux and register handlers
	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/",ph.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}",ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/",ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	s := &http.Server{
		Addr:  *bindAddress,
		Handler: sm,
		IdleTimeout: 5 * time.Second,
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <- sigChan
	l.Printf("Received terminate, graceful shutdown %v\n", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)


}
