package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Starting")

	r := mux.NewRouter()
	r.Use(
		loggingMiddleware,
	)

	hostname, _ := os.Hostname()
	data := struct {
		Hostname string `json:"hostname"`
		Env      string `json:"env"`
		Username string `json:"username"`
		Password string `json:"password"`
		LogLevel string `json:"logLevel"`
	}{
		Hostname: hostname,
		Env:      os.Getenv("ENV"),
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
		LogLevel: os.Getenv("LOG_LEVEL"),
	}

	r.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	})
	r.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`OK`))
	})

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		log.Println("The service is ready to listen and serve")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("error occur while shutting down server:", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	wait := 60 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatal("error occur while shutting down:", err)
	}

	log.Println("The service is gracefully shutdown")
	os.Exit(0)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("received request")
		next.ServeHTTP(w, r)
	})
}
