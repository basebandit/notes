package main

import (
	"github.com/parish/notes/api"
)

func startServer() {
	cfg, err := getConfig()
	if err != nil {
		logger.Fatalf("failed to load configuration: %s", err)
	}

	store, err := getStore(cfg)
	if err != nil {
		logger.Fatalf("failed to init data store: %s", err)
	}

	api, err := api.NewServer(store, logger)
	if err != nil {
		logger.Fatalf("error starting api server: %s", err)
	}

	api.Start()
}
