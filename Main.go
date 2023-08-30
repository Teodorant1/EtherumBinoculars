package main

import (
	"fmt"
	"net/http"
)

func main() {

	fetchr := Fetcher{}
	coins := fetchr.Get_Eth_At_Date(1688357340, "0x388c818ca8b9251b393131c08a736a67ccb19297")
	fmt.Println(coins)
	fetchr.Grab_etherium_transactions("0x2c1ba59d6f58433fb1eaee7d20b26ed83bda51a3", 0)

	server := http.NewServeMux()
	fetchhandler := &eth_API_Handler{fetchr: &Fetcher{}}
	fetchhandler2 := &eth_API_HandlerV2{fetchr: &Fetcher{}}
	fmt.Println("fetchr server is online")
	server.Handle("/gettransactions", fetchhandler)
	server.Handle("/getbalance", fetchhandler2)
	err := http.ListenAndServe(":8001", server)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
}
