package main

import (
	"flag"
	"log"
	"net/http"
)

type Config struct {
	addr string
}

var cfg Config

func main() {
	var err error
	flag.StringVar(&cfg.addr, "addr", ":5000", "HTTP address bound")

	mux := http.NewServeMux()

	err = http.ListenAndServe(cfg.addr, mux)
	log.Fatal(err)
}
