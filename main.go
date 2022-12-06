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
var ethValue *big.Float
var rpcendpoint = os.Getenv("RPCENDPOINT")

func getHttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
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

	fmt.Printf("Account: %s Balance: %v ETH\n", address, lev)
	return lev
}
func main() {
	ctx = context.Background()

	port := ":8080"
	fmt.Println("go-web3-sample-app is running at: http://localhost" + port)

	fmt.Println("RPCENDPOINT:", rpcendpoint)

	httpClient := getHttpClient()
	ethClient := getEthClient(httpClient)

	address = "0xeB2629a2734e272Bcc07BDA959863f316F4bD4Cf"
	ethValue = getBalance(address, ethClient, ctx)
	tmpl := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		details := BalanceDetails{
			Address: r.FormValue("address"),
		}

		if r.Method == http.MethodPost {
			address = details.Address
			ethValue = getBalance(address, ethClient, ctx)
		}

		tmpl.Execute(w, struct {
			Success bool
			Address string
			Balance *big.Float
		}{true, address, ethValue})
		tmpl.Execute(w, struct{ Success bool }{true})
	})

	http.ListenAndServe(port, nil)
}
