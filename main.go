package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BalanceDetails struct {
	Account   string
}

var address string
var fbalance *big.Float
var ethValue *big.Float

func getBalance(client *ethclient.Client)  {
	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), account, nil)
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
	rpcendpoint:=os.Getenv("RPCENDPOINT")

	fmt.Println("RPCENDPOINT: ",rpcendpoint)

	client, err := ethclient.Dial(rpcendpoint)
	if err != nil {
		log.Fatal(err)
	}

	address = "0x71c7656ec7ab88b098defb751b7401b5f6d8976f"
	getBalance(client)
	tmpl := template.Must(template.ParseFiles("forms.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			getBalance(client)
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
