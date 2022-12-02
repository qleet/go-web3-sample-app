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
	Account   string
}

var ctx context.Context
var address string
var fbalance *big.Float
var ethValue *big.Float
var rpcendpoint = os.Getenv("RPCENDPOINT")
//var rpcendpoint="https://1rpc.io/eth"

func getBalance(client *ethclient.Client, context context.Context)  {
	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(context, account, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance)
	fbalance = new(big.Float)
	fbalance.SetString(balance.String())
	ethValue = new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	fmt.Println(ethValue)
}

func main() {

	ctx = context.Background()

	fmt.Println("RPCENDPOINT: ",rpcendpoint)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
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

	address = "0x71c7656ec7ab88b098defb751b7401b5f6d8976f"
	getBalance(ethClient, ctx)
	tmpl := template.Must(template.ParseFiles("forms.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			getBalance(ethClient, ctx)
			tmpl.Execute(w, struct{ Success bool; Account string;  Balance  *big.Float}{true, address, ethValue})
			return
		}

		details := BalanceDetails{
			Account:   r.FormValue("account"),
		}

		// do something with details
		_ = details

		tmpl.Execute(w, struct{ Success bool }{true})
	})

	fmt.Println("Running at: http://localhost:9090")
	http.ListenAndServe(":9090", nil)
}
