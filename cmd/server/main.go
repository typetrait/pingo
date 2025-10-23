package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/typetrait/pingo/cmd/server/networking"
)

func main() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	s := networking.NewServer()
	err := s.Start()
	if err != nil {
		log.Fatal(err)
	}
}
