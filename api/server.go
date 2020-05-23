package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/parish/notes/store"
)

//Server defines our api server dependencies
type Server struct {
	*http.Server
}

//NewServer creates and configures an API server instance serving all application resource routes.
func NewServer(store store.Store, logger *log.Logger) (*Server, error) {
	api, err := New(true, store, logger)
	if err != nil {
		return nil, err
	}

	srv := http.Server{
		Addr:         ":8080",
		Handler:      api,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return &Server{&srv}, nil
}

//Start runs ListenAndServe on the http.Server with graceful shutdown.
func (srv *Server) Start() {
	log.Println("Starting API Server...")
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	log.Printf("Listening on %s\n", srv.Addr)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Println("Shutting down server... Reason:", sig)
	//teardown logic here

	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	log.Println("Server gracefully stopped.")
}
