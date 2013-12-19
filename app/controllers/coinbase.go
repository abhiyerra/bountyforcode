package bountyforcode

import (
	"fmt"
	. "github.com/abhiyerra/bountyforcode/app/models"
	"github.com/abhiyerra/coinbase"
	"io/ioutil"
	"log"
	"net/http"
)

func CoinbaseCallbackHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Coinbase callback fail %v", err)
	}

	order := coinbase.ParseCallback(string(body))

	bounty := FindBountyByCoinbaseButtonCode(order.Button.Id)
	if bounty == nil {
		log.Printf("Received unknown coinbase callback")
		fmt.Fprintf(w, "bad")
	}

	bounty.UpdateCoinbaseInfo(order)

	fmt.Fprintf(w, "ok")
}
