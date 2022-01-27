package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/kevin/working/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	ph := handlers.NewProducts(l)

	sm := mux.NewRouter()

	getRouter := sm.Methods("GET").Subrouter()
	getRouter.HandleFunc("/",ph.GetProducts)

	putRouter := sm.Methods("PUT").Subrouter()
	putRouter.HandleFunc("/",ph.UpdateProducts)

	s := &http.Server{
		Addr: ":9090",
		Handler: sm,
		IdleTimeout: 120 * time.Second,
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
