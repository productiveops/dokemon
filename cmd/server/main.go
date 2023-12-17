package main

import (
	"dokemon/pkg/server"
	"os"
)

func main() {
	s := server.Server{}
	s.Init(
		os.Getenv("DB_CONNECTION_STRING"),
		os.Getenv("DATA_PATH"),
		os.Getenv("LOG_LEVEL"),
	)
	s.Run(os.Getenv("BIND_ADDRESS"))
}