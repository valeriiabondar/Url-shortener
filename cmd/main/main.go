package main

import (
	"fmt"
	"urlShortener/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)
}
