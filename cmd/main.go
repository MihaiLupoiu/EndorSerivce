package main

import (
	"log"

	"github.com/hiring-seedtag/mihai-lupoiu-go-backend-test/internal/server"
)

func main() {
	a := server.App{}

	err := a.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	a.Run()
}
