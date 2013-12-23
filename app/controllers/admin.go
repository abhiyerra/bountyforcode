package bountyforcode

import (
	"fmt"
	"net/http"
)

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", GetSessionUser(r))

	// button := coinbase.GetButton(&coinbase.ButtonRequest{
	// 	Name: "Abhi Yerra",
	// 	Type: "donation",
	// 	PriceString: "10.00",
	// 	PriceCurrencyIso: "USD",
	// })

	// fmt.Printf("%v", button.Response.Button.Code)
	// fmt.Fprintf(w, "hi")
}
