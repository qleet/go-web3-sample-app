package main

import (
	"context"
	"crypto/tls"
	"html/template"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type BalanceDetails struct {
	Address string
}

var ctx context.Context
var address string
var ethValue *big.Float
var rpcendpoint = os.Getenv("RPCENDPOINT")

func getHttpClient(tr *http.Transport) *http.Client {
	// reuse connections https://stuartleeks.com/posts/connection-re-use-in-golang-with-http-client/
	return &http.Client{Transport: tr}
}

func getEthClient(httpClient *http.Client) *ethclient.Client {
	rpcClient, err := rpc.DialHTTPWithClient(rpcendpoint, httpClient)
	if err != nil {
		log.Fatal(err)
	}
	ethClient := ethclient.NewClient(rpcClient)
	if err != nil {
		log.Fatal(err)
	}

	return ethClient
}

func getBalance(address string, ethclient *ethclient.Client, context context.Context) *big.Float {
	account := common.HexToAddress(address)
	balance, err := ethclient.BalanceAt(context, account, nil)
	if err != nil {
		log.Fatal(err)
	}

	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	lev := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	log.Printf("Account: %s Balance: %v ETH\n", address, lev)
	return lev
}
func main() {
	ctx = context.Background()

	port := ":8080"
	log.Println("go-web3-sample-app is running at: http://localhost" + port)

	log.Println("RPCENDPOINT:", rpcendpoint)

	// set http Transport
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// set http Client
	httpClient := getHttpClient(transport)
	// set eth Client
	ethClient := getEthClient(httpClient)

	// initial address
	address = "0xeB2629a2734e272Bcc07BDA959863f316F4bD4Cf"
	// get initial balance for the initial address
	ethValue = getBalance(address, ethClient, ctx)
	tmpl := template.Must(template.ParseFiles("index.html"))

	// handle HTTP requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// call getBalance only on POST otherwise use ethValue from previous getBalance call
		if r.Method == http.MethodPost {
			// get details from HTML form
			details := BalanceDetails{
				Address: r.FormValue("address"),
			}
			// set address to value from HTML form
			address = details.Address
			// get ETH balance into ethValue
			ethValue = getBalance(address, ethClient, ctx)
		}

		tmpl.Execute(w, struct {
			Success bool
			Address string
			Balance *big.Float
		}{true, address, ethValue})

		//tmpl.Execute(w, struct{ Success bool }{true})
	})

	http.ListenAndServe(port, nil)
}
