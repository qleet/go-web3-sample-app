package main

import (
	"context"
	"crypto/tls"
	"fmt"
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
var fbalance *big.Float
var ethValue *big.Float
var rpcendpoint = os.Getenv("RPCENDPOINT")

func getBalance(address string, client *ethclient.Client, context context.Context) {
	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(context, account, nil)
	if err != nil {
		log.Fatal(err)
	}

	fbalance = new(big.Float)
	fbalance.SetString(balance.String())
	ethValue = new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	fmt.Println(ethValue)
}

func main() {

	ctx = context.Background()

	fmt.Println("RPCENDPOINT: ", rpcendpoint)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	rpcClient, err := rpc.DialHTTPWithClient(rpcendpoint, httpClient)
	if err != nil {
		log.Fatal(err)
	}
	ethClient := ethclient.NewClient(rpcClient)
	if err != nil {
		log.Fatal(err)
	}

	address = "0xeB2629a2734e272Bcc07BDA959863f316F4bD4Cf"
	getBalance(address, ethClient, ctx)
	tmpl := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		details := BalanceDetails{
			Address: r.FormValue("address"),
		}

		if r.Method == http.MethodPost {
			address = details.Address
		}

		getBalance(address, ethClient, ctx)
		tmpl.Execute(w, struct {
			Success bool
			Address string
			Balance *big.Float
		}{true, address, ethValue})
		tmpl.Execute(w, struct{ Success bool }{true})
	})

	port := ":8080"
	fmt.Println("Running at: http://localhost" + port)
	http.ListenAndServe(port, nil)
}
