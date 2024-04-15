package main

import (
	"log"

	"github.com/koliader/posts-gateway/internal/util"
	"github.com/koliader/posts-gateway/pkg/v1/handler/api"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	server, err := api.NewServer(config)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
