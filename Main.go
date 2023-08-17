package main

import (
	"fmt"
	"net/http"
)

func main() {
	fetchr := Fetcher{}
	fetchr.Grab_etherium_transactions("0x2c1ba59d6f58433fb1eaee7d20b26ed83bda51a3", 600000)

	server := http.NewServeMux()
	fetchhandler := eth_API_Handler{fetchr: Fetcher{}}
	fmt.Println("fetchr server is online")
	server.Handle("/", fetchhandler)
	http.ListenAndServe(":8001", server)
}
