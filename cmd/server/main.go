package main

import (
	"os"

	"github.com/productiveops/dokemon/pkg/server"
)

func main() {
	s := server.NewServer(
		os.Getenv("DB_CONNECTION_STRING"),
		os.Getenv("DATA_PATH"),
		os.Getenv("LOG_LEVEL"),
		os.Getenv("SSL_ENABLED"),
		os.Getenv("STALENESS_CHECK"),
	)
	s.Run(os.Getenv("BIND_ADDRESS"))
}