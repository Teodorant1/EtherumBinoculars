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

// the first member of payload is the latest transaction and the last one is the first transaction
func (Fetcher Fetcher) filterDates(timestamp int, payload api_Payload, currentbalance int, walletaddress string) int {
	//	fmt.Println(payload.Result[0].Timestamp)
	//	fmt.Println(payload.Result[1].Timestamp)
	presentDay := time.Now().Unix()
	periods := []*period{}
	valid_transactions := []transaction{}
	writtenNumber, _ := strconv.Atoi(payload.Result[0].Timestamp)
	// by doing this we add the latest period
	periods = append(periods, &period{
		value:     currentbalance,
		startdate: writtenNumber,
		enddate:   int(presentDay),
	})

	//starting from the latest transaction and working our way back
	for i := 0; i < len(payload.Result)-1; i++ {
		if payload.Result[i].IsError == "0" && i < len(payload.Result)-1 {
			//valueint, _ := strconv.Atoi(payload.Result[i].VALUE)
			// periods are sandwiched with transactions
			startint, _ := strconv.Atoi(payload.Result[i+1].Timestamp)
			endint, _ := strconv.Atoi(payload.Result[i].Timestamp)
			period := period{
				value:     0,
				startdate: startint,
				enddate:   endint,
			}
			periods = append(periods, &period)
			valid_transactions = append(valid_transactions, payload.Result[i])
			//	fmt.Println("sizes be like", len(periods), len(valid_transactions))
		}
		//		payload.Result

		// if it's the earliest transaction
		if payload.Result[i].IsError == "0" && i == len(payload.Result)-1 {
			endint, _ := strconv.Atoi(payload.Result[len(payload.Result)-1].Timestamp)

			period := period{
				value:     0,
				startdate: 0,
				enddate:   endint,
			}
			periods = append(periods, &period)
			valid_transactions = append(valid_transactions, payload.Result[i])
		}

	}
	// 1   2    3    4     5
	//	//       1    2    3    4     5
	//
	//	fmt.Println(periods)
	//	fmt.Println(valid_transactions)

	for i := 1; i < len(valid_transactions); i++ {
		valueint, _ := strconv.Atoi(valid_transactions[i].VALUE)
		if valid_transactions[i].TO == walletaddress {
			currentbalance = currentbalance - valueint
		}
		if valid_transactions[i].FROM == walletaddress {
			currentbalance = currentbalance + valueint
		}
		periods[i].value = currentbalance

	}
	//for i2 := range periods {
	//	fmt.Println("index is", i2)
	//	fmt.Println("$", periods[i2].value)
	//	fmt.Println("stardate", periods[i2].startdate)
	//	fmt.Println("enddate", periods[i2].enddate)
	//	fmt.Println("////////////////////////////////////////////////")
	//}

	// solution := new(int)

	for i := range periods {
		if periods[i].startdate < timestamp && periods[i].enddate > timestamp {
			fmt.Println(periods[i].value)
			return periods[i].value
			//	return periods[i].value
		}
	}

	return 0
}

func (Fetcher Fetcher) Get_Eth_At_Date(timestamp int, wallet string) int {
	reqStr := "https://api.etherscan.io/api?module=account&action=balance&address=" + wallet + "&tag=latest&apikey=YourApiKeyToken"
	resp, _ := http.Get(reqStr)
	defer resp.Body.Close()
	content, _ := io.ReadAll(resp.Body)
	balance := balance{}
	json.Unmarshal(content, &balance)
	transactions_payload := Fetcher.Grab_etherium_transactions(wallet, 1)
	balanceattime := Fetcher.filterDates(timestamp, transactions_payload, balance.Balance, wallet)

	return balanceattime
}

func (Fetcher Fetcher) Grab_etherium_transactions(wallet string, block int) api_Payload {
	endblock := strconv.Itoa(block + 100000000000000)
	block_string := strconv.Itoa(block)
	inputstr := "https://api.etherscan.io/api?module=account&action=txlistinternal&address=" + wallet + "&startblock=" + block_string + "&endblock=" + endblock + "&page=1&offset=10000&sort=desc&apikey=" + apikey
	//	fmt.Println(inputstr)
	resp, _ := http.Get(inputstr)
	//resp, err := http.Get("https://api.etherscan.io/api?module=account&action=txlistinternal&address=0x2c1ba59d6f58433fb1eaee7d20b26ed83bda51a3&startblock=0&endblock=100000000000000&offset=10000&sort=desc&apikey=8SYNIAHH11X4UA7SGZI8JHICTTAZPZVJXE")
	// if err != nil {
	// 	panic(err)
	// }
	defer resp.Body.Close()
	content, _ := io.ReadAll(resp.Body)
	fmt.Println(string(content))
	//fmt.Println("success")
	//fmt.Println("//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////")
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
