package main

import (
	"url-shortoner/internal/config"
)

func main() {
	cfg := config.MustLoad()

	// TODO: init logger (slog)

	// TODO: init storage (sqlite)

	// TODO: init router (chi | chi render)

	// TODO: run server
}

func setupLogger(env string) {
	switch expr {}
}