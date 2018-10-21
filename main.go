package main

import "github.com/MudaeH5A/4thinkbe/api"

func main() {
	server := api.New()
	server.Listen()
}
