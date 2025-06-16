package main

import (
	"github.com/brienze1/notes-api/internal/infra/injections"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("./env/.env.localstack")
	if err != nil {
		err := godotenv.Load("../../env/.env.localstack")
		if err != nil {
			panic(err)
		}
	}
	server, err := injections.Wire().InitializeServer()
	if err != nil {
		panic(err)
	}

	err = server.Serve()
	if err != nil {
		panic(err)
	}
}
