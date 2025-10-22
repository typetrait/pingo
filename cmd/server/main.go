package main

import (
	"log"

	"github.com/typetrait/pingo/cmd/server/networking"
)

func main() {
	s := networking.NewServer()
	err := s.Start()
	if err != nil {
		log.Fatal(err)
	}
}
