package main

import (
	"fmt"
	"net/http"
)

func main() {

	//	var pow = []int{1, 2, 4, 8, 16, 32, 64, 128}
	//
	//	for i := 0; i < len(pow)-1; i++ {
	//		fmt.Println(pow[i])
	//	}

	fetchr := &Fetcher{}
	coins := fetchr.Get_Eth_At_Date(1692316800000, "0x1B0b90234f5a29CD4486a4E55f3A73B43E21Bc45")
	fmt.Println(coins)
	fetchr.Grab_etherium_transactions("0x2c1ba59d6f58433fb1eaee7d20b26ed83bda51a3", 0)

	server := http.NewServeMux()
	fetchhandler := &eth_API_Handler{fetchr: &Fetcher{}}
	fetchhandler2 := &eth_API_HandlerV2{fetchr: &Fetcher{}}
	fmt.Println("fetchr server is online")
	server.Handle("/gettransactions", fetchhandler)
	server.Handle("/getbalance", fetchhandler2)
	http.ListenAndServe(":8001", server)
}
