package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type eth_API_Handler struct {
	fetchr Fetcher
}

type eth_API_HandlerV2 struct {
	fetchr Fetcher
}

type Request struct {
	Address string `json:"address"`
	Block   int    `json:"block"`
}

type BalanceRequest struct {
	Address string `json:"address"`
	Date    int    `json:"date"`
}

type BalanceResponse struct {
	balanceattime int
}

func (ethHandler eth_API_HandlerV2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Context-Type", "application/json")
	body1, _ := io.ReadAll(r.Body)

	var request BalanceRequest
	err := json.Unmarshal(body1, &request)
	if err != nil {
		fmt.Print("JSON error", err)
	}
	fetchr := ethHandler.fetchr

	response := fetchr.Get_Eth_At_Date(request.Date, request.Address)
	reply, _ := json.Marshal(response)
	w.Write((reply))

}
func (ethHandler eth_API_Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Context-Type", "application/json")

	body1, _ := io.ReadAll(r.Body)

	var request Request

	//	fmt.Println(string(body1))

	err := json.Unmarshal(body1, &request)

	if err != nil {
		fmt.Print("JSON error", err)
	}
	fmt.Println("address be like", request.Address)
	fmt.Println("block be like", request.Block)
	fetchr := ethHandler.fetchr
	api_Payload := fetchr.Grab_etherium_transactions(request.Address, request.Block)
	reply, _ := json.Marshal(api_Payload)
	w.Write((reply))
}
