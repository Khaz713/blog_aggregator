package main

import (
	config "blog_aggregator/internal/config"
	"log"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = cfg.SetUser("Khaz")
	if err != nil {
		log.Fatal(err)
		return
	}
	cfg, err = config.Read()
	if err != nil {
		log.Fatal(err)
		return
	}

}
