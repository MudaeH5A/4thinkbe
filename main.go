package main

import "github.com/arxdsilva/4thinkbe/api"

func main() {
	server := api.New()
	server.Listen()
}
