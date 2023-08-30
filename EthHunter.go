package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
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
	Balance string `json:"result"`
}
type period struct {
	value     int
	startdate int
	enddate   int
}

const apikey = "8SYNIAHH11X4UA7SGZI8JHICTTAZPZVJXE"

type Fetcher struct{}

// the first member of payload is the latest transaction and the last one is the first transaction
func (Fetcher *Fetcher) filterDates(timestamp int, payload api_Payload, currentbalance int, walletaddress string) int {

	zeroBalance := 0

	presentDay := time.Now().Unix()

	periods := []*period{}
	periods2 := []period{}
	valid_transactions := []transaction{}
	writtenNumber, _ := strconv.Atoi(payload.Result[0].Timestamp)
	// by doing this we add the latest period

	periods = append(periods, &period{
		value:     currentbalance,
		startdate: writtenNumber,
		enddate:   int(presentDay),
	})
	periods2 = append(periods2, period{
		value:     currentbalance,
		startdate: writtenNumber,
		enddate:   int(presentDay),
	})
	//starting from the latest transaction and working our way back
	for i := 0; i < len(payload.Result)-1; i++ {
		if payload.Result[i].IsError == "0" && i < len(payload.Result)-1 {

			startint, _ := strconv.Atoi(payload.Result[i+1].Timestamp)
			endint, _ := strconv.Atoi(payload.Result[i].Timestamp)

			period := period{
				value:     0,
				startdate: startint,
				enddate:   endint,
			}
			periods = append(periods, &period)
			valid_transactions = append(valid_transactions, payload.Result[i])

		}

	}

	for i := 1; i < len(valid_transactions); i++ {

		currentPeriod := periods[i]

		valueint, _ := strconv.Atoi(valid_transactions[i].VALUE)
		if strings.ToLower(valid_transactions[i].FROM) == strings.ToLower(walletaddress) {
			newvalue := currentPeriod.value + valueint
			currentbalance = newvalue

			newperiod := period{
				value:     newvalue,
				startdate: periods[i].startdate,
				enddate:   periods[i].enddate,
			}

			periods2 = append(periods2, newperiod)

		}
		if strings.ToLower(valid_transactions[i].TO) == strings.ToLower(walletaddress) {
			newvalue := currentPeriod.value - valueint
			currentbalance = newvalue
			newperiod := period{
				value:     newvalue,
				startdate: periods[i].startdate,
				enddate:   periods[i].enddate,
			}

			periods2 = append(periods2, newperiod)

		}

	}

	for i := 0; i < len(periods2); i++ {

		if periods2[i].startdate < timestamp && periods2[i].enddate > timestamp {
			zeroBalance = periods2[i].value
			fmt.Println("zerobalance:", zeroBalance)
		}
	}

	return zeroBalance
}

func (Fetcher *Fetcher) Get_Eth_At_Date(timestamp int, wallet string) int {
	reqStr := "https://api.etherscan.io/api?module=account&action=balance&address=" + wallet + "&tag=latest&apikey=" + apikey
	resp, _ := http.Get(reqStr)
	defer resp.Body.Close()
	content, _ := io.ReadAll(resp.Body)
	balance := balance{}
	json.Unmarshal(content, &balance)
	transactions_payload := Fetcher.Grab_etherium_transactions(wallet, 0)
	fmt.Println("Balance Json", string(content))

	actualBalanceInt, _ := strconv.Atoi(balance.Balance)
	fmt.Println("Balance", actualBalanceInt)
	fmt.Println("before filter dates", transactions_payload.Result[0].FROM, transactions_payload.Result[0].TO, wallet)
	balanceattime := Fetcher.filterDates(timestamp, transactions_payload, actualBalanceInt, wallet)

	return balanceattime
}

func (Fetcher *Fetcher) Grab_etherium_transactions(wallet string, block int) api_Payload {
	endblock := strconv.Itoa(block + 100000000000000)
	block_string := strconv.Itoa(block)
	inputstr := "https://api.etherscan.io/api?module=account&action=txlist&address=" + wallet + "&startblock=" + block_string + "&endblock=" + endblock + "&page=1&offset=10000&sort=desc&apikey=" + apikey
	resp, _ := http.Get(inputstr)

	defer resp.Body.Close()
	content, _ := io.ReadAll(resp.Body)
	payload := api_Payload{}
	json.Unmarshal(content, &payload)

	return payload

}
