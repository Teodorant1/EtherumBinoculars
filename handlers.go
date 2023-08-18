package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type eth_API_Handler struct {
	fetchr *Fetcher
}

type eth_API_HandlerV2 struct {
	fetchr *Fetcher
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
	Balanceattime int `json:"balanceattime"`
}

func (ethHandler *eth_API_HandlerV2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("fetching balance at date")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Context-Type", "application/json")
	body1, _ := io.ReadAll(r.Body)

	var request BalanceRequest
	err := json.Unmarshal(body1, &request)
	if err != nil {
		fmt.Print("JSON error", err)
	}
	fmt.Println("Address:", request.Address)
	fmt.Println("Date", request.Date)
	fetchr := ethHandler.fetchr

	balanceint := fetchr.Get_Eth_At_Date(request.Date, request.Address)
	replystruct := BalanceResponse{balanceint}
	fmt.Println(replystruct)
	reply, _ := json.Marshal(replystruct)
	//fmt.Println(reply)
	w.Write(reply)

}
func (ethHandler *eth_API_Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Print("fetching transactions")
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
	//fmt.Println(string(reply))
	w.Write((reply))
}
