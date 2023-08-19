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
	//	fmt.Println(payload.Result[0].Timestamp)
	//	fmt.Println(payload.Result[1].Timestamp)
	zeroBalance := 0

	//	fmt.Println("last transaction0 :", payload.Result[0].Timestamp)
	//	fmt.Println("last transaction1 :", payload.Result[1].Timestamp)
	//	fmt.Println("last transaction2 :", payload.Result[2].Timestamp)
	//	fmt.Println("last transaction3 :", payload.Result[3].Timestamp)
	//	fmt.Println("last transaction4 :", payload.Result[4].Timestamp)

	//	fmt.Println("current balance: ", currentbalance)
	//	currentbalance2 := &currentbalance

	presentDay := time.Now().Unix()

	//	fmt.Println("epoch time:", presentDay)

	periods := []*period{}
	periods2 := []period{}
	//	fmt.Println("length of payload", len(payload.Result))
	valid_transactions := []transaction{}
	writtenNumber, _ := strconv.Atoi(payload.Result[0].Timestamp)
	// by doing this we add the latest period
	//	fmt.Println("1st stardate and present day", writtenNumber, int(presentDay))
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
			//	if payload.Result[i].VALUE != "0" {
			//		//fmt.Println("Value:", payload.Result[i].VALUE)
			//		//fmt.Println("not zero")
			//	}
			//valueint, _ := strconv.Atoi(payload.Result[i].VALUE)
			// periods are sandwiched with transactions
			startint, _ := strconv.Atoi(payload.Result[i+1].Timestamp)
			endint, _ := strconv.Atoi(payload.Result[i].Timestamp)

			//	("addresses", strings.ToUpper(), strings.ToUpper(), strings.ToUpper())

			period := period{
				value:     0,
				startdate: startint,
				enddate:   endint,
			}
			periods = append(periods, &period)
			valid_transactions = append(valid_transactions, payload.Result[i])
			//	fmt.Println("og addresses", valid_transactions[i].FROM, valid_transactions[i].TO, walletaddress)

			//	fmt.Println("sizes be like", len(periods), len(valid_transactions))
		}
		//		payload.Result

		// if it's the earliest transaction
		//	if payload.Result[i].IsError == "0" && i == len(payload.Result)-1 {
		//		endint, _ := strconv.Atoi(payload.Result[len(payload.Result)-1].Timestamp)
		//		period := period{
		//			value:     0,
		//			startdate: 0,
		//			enddate:   endint,
		//		}
		//		periods = append(periods, &period)
		//		valid_transactions = append(valid_transactions, payload.Result[i])
		//	}

	}
	// 1   2    3    4     5
	//	//       1    2    3    4     5
	//
	//	fmt.Println(periods)
	//	fmt.Println(valid_transactions)

	for i := 1; i < len(valid_transactions); i++ {

		currentPeriod := periods[i]

		valueint, _ := strconv.Atoi(valid_transactions[i].VALUE)
		//	fmt.Println("Current Balance", (*currentPeriod).value)
		//	fmt.Println("ValueInt ", valueint)
		//
		//	fmt.Println("addresses", strings.ToUpper(valid_transactions[i].FROM), strings.ToUpper(valid_transactions[i].TO), strings.ToUpper(walletaddress))
		//	fmt.Println(strings.ToUpper(valid_transactions[i].TO) == strings.ToUpper(walletaddress), strings.ToUpper(valid_transactions[i].FROM) == strings.ToUpper(walletaddress))
		////////////////////////////////////////////////////////////////////////////////////////////////////
		//		if strings.ToUpper(valid_transactions[i].TO) == strings.ToUpper(walletaddress) {
		//			(*currentPeriod).value = (*currentPeriod).value - valueint
		//			//		fmt.Println("minus", currentbalance-valueint)
		//		}
		//		if strings.ToUpper(valid_transactions[i].FROM) == strings.ToUpper(walletaddress) {
		//			(*currentPeriod).value = (*currentPeriod).value + valueint
		//			//		fmt.Println("plus", currentbalance+valueint)
		//		}
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		if strings.ToLower(valid_transactions[i].FROM) == strings.ToLower(walletaddress) {
			newvalue := currentPeriod.value + valueint
			currentbalance = newvalue

			newperiod := period{
				value:     newvalue,
				startdate: periods[i].startdate,
				enddate:   periods[i].enddate,
			}

			periods2 = append(periods2, newperiod)

			//		fmt.Println("plus", currentbalance+valueint)
		}
		if strings.ToLower(valid_transactions[i].TO) == strings.ToLower(walletaddress) {
			newvalue := currentPeriod.value - valueint
			currentbalance = newvalue
			newperiod := period{
				value:     newvalue,
				startdate: periods[i].startdate,
				enddate:   periods[i].enddate,
			}
			//	fmt.Println("value of current period", newperiod.value)

			periods2 = append(periods2, newperiod)

			//	fmt.Println("minus", currentbalance-valueint)
		}
		//	fmt.Println("current balance", currentbalance)

		//	currentPerValue = currentbalance

		//	(currentPeriod).value = currentbalance
		//	periods[i].value = currentbalance
		//	fmt.Println(periods[i].value)
		//
		//	fmt.Println(" Mutated Current Balance ", (*currentPeriod).value)
		//	fmt.Println("value of current slice element", periods[i].value)

	}
	//for i2 := range periods {
	//	fmt.Println("index is", i2)
	//	fmt.Println("$", periods[i2].value)
	//	fmt.Println("stardate", periods[i2].startdate)
	//	fmt.Println("enddate", periods[i2].enddate)
	//	fmt.Println("////////////////////////////////////////////////")
	//}

	// solution := new(int)

	for i := 0; i < len(periods2); i++ {
		//		fmt.Println("value:", periods[i].value)

		//	fmt.Println("startdate:", periods2[i].startdate, "enddate:", periods2[i].enddate, "timestamp:", timestamp)

		if periods2[i].startdate < timestamp && periods2[i].enddate > timestamp {
			//	fmt.Println(periods[i].value)
			zeroBalance = periods2[i].value
			fmt.Println("zerobalance:", zeroBalance)
			//	return periods[i].value
		}
	}
	//	fmt.Println("last value", periods2[0].value)
	//	fmt.Println("last sdate", periods2[0].startdate)
	//	fmt.Println("last edate", periods2[0].enddate)
	//
	//	fmt.Println("first value", periods2[len(periods)-1].value)
	//	fmt.Println("first sdate", periods2[len(periods)-1].startdate)
	//	fmt.Println("first edate", periods2[len(periods)-1].enddate)
	//
	//	fmt.Println("periods 2", periods2)

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
	inputstr := "https://api.etherscan.io/api?module=account&action=txlistinternal&address=" + wallet + "&startblock=" + block_string + "&endblock=" + endblock + "&page=1&offset=10000&sort=desc&apikey=" + apikey
	//	fmt.Println(inputstr)
	resp, _ := http.Get(inputstr)
	//resp, err := http.Get("https://api.etherscan.io/api?module=account&action=txlistinternal&address=0x2c1ba59d6f58433fb1eaee7d20b26ed83bda51a3&startblock=0&endblock=100000000000000&offset=10000&sort=desc&apikey=8SYNIAHH11X4UA7SGZI8JHICTTAZPZVJXE")
	// if err != nil {
	// 	panic(err)
	// }
	defer resp.Body.Close()
	content, _ := io.ReadAll(resp.Body)
	//	fmt.Println(string(content))
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
