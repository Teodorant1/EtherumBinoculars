package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type date struct {
	year  int
	month int
	date  int
}

type transaction struct {
	Blocknumber string `json:"blockNumber"`
	VALUE       string `json:"value"`
	FROM        string `json:"from"`
	TO          string `json:"to"`
	IsError     string `json:"isError"`
	ErrCode     string `json:"errCode"`
	Timestamp   string `json:"timeStamp"`
	Hash        string `json:"hash"`
}

type api_Payload struct {
	Message string        `json:"message"`
	Status  string        `json:"status"`
	Result  []transaction `json:"result"`
}

type balance struct {
	Balance int `json:"balance"`
}
type period struct {
	value     int
	startdate int
	enddate   int
}

const apikey = "8SYNIAHH11X4UA7SGZI8JHICTTAZPZVJXE"

type Fetcher struct{}

func (Fetcher Fetcher) filterDates(timestamp int, payload api_Payload, currentbalance int) {
	presentDay := time.Now().Unix()
	periods := []*period{}

	writtenNumber, _ := strconv.Atoi(payload.Result[len(payload.Result)-1].Timestamp)

	periods = append(periods, &period{
		value:     currentbalance,
		startdate: writtenNumber,
		enddate:   int(presentDay),
	})

	for i:= len(payload.Result)-1; i>= 0 ; i-- {
		if payload.Result[i].IsError == "0" {
			period:= period{
				value:     0,
				startdate: 0,
				enddate:   0,
			}
			}
		}
	}

}

func (Fetcher Fetcher) Get_Eth_At_Date(timestamp int, wallet string) balance {
	reqStr := "https://api.etherscan.io/api?module=account&action=balance&address=" + wallet + "&tag=latest&apikey=YourApiKeyToken"
	resp, _ := http.Get(reqStr)
	defer resp.Body.Close()
	content, _ := io.ReadAll(resp.Body)
	balance := balance{}
	json.Unmarshal(content, &balance)

	transactions_payload := Fetcher.Grab_etherium_transactions(wallet, 1)

	Fetcher.filterDates(timestamp, transactions_payload, balance.Balance)

	return balance
}

func (Fetcher Fetcher) Grab_etherium_transactions(wallet string, block int) api_Payload {
	endblock := strconv.Itoa(block + 100000000000000)
	block_string := strconv.Itoa(block)
	inputstr := "https://api.etherscan.io/api?module=account&action=txlistinternal&address=" + wallet + "&startblock=" + block_string + "&endblock=" + endblock + "&page=1&offset=10000&sort=desc&apikey=" + apikey
	fmt.Println(inputstr)
	resp, _ := http.Get(inputstr)
	//resp, err := http.Get("https://api.etherscan.io/api?module=account&action=txlistinternal&address=0x2c1ba59d6f58433fb1eaee7d20b26ed83bda51a3&startblock=0&endblock=100000000000000&offset=10000&sort=asc&apikey=8SYNIAHH11X4UA7SGZI8JHICTTAZPZVJXE")
	// if err != nil {
	// 	panic(err)
	// }
	defer resp.Body.Close()
	content, _ := io.ReadAll(resp.Body)
	fmt.Println(string(content))
	fmt.Println("success")
	fmt.Println("//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////")
	payload := api_Payload{}

	json.Unmarshal(content, &payload)

	//	fmt.Println(payload)

	//	for _, transaction := range payload.Result {
	//		//	fmt.Println(transaction)
	//		fmt.Println(transaction.Timestamp)
	//		fmt.Println("blocknumber", transaction.Blocknumber)
	//	}

	return payload

}
