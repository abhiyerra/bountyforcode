package bountyforcode

import (
	"github.com/abhiyerra/coinbase"
	"log"
	"strconv"
	"time"
)

const (
	BountyStateNew       string = "new"
	BountyStatePaid      string = "paid"
	BountyStateClosed    string = "closed"
	BountyStateCancelled string = "cancelled"
)

type Bounty struct {
	Id                  int       `db:"id" json:"id"`
	UserId              int       `db:"user_id" json:"user_id"`
	IssueId             int       `db:"issue_id" json:"issue_id"`
	Amount              float32   `db:"amount" json:"amount"`
	CoinbaseButtonCode  string    `db:"coinbase_button_code" json:"coinbase_button_code"`
	CoinbaseOrderId     string    `db:"coinbase_order_id" json:"-"`
	CoinbaseTotalBtc    int       `db:"coinbase_total_btc" json:"coinbase_total_btc"`
	CoinbaseCurrencyIso string    `db:"coinbase_currency_iso" json:"coinbase_currency_iso"`
	Status              string    `db:"status" json:"status"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`
}

func NewBounty(issue *Issue, user_id int) (b *Bounty) {
	b = &Bounty{
		UserId:  user_id,
		IssueId: issue.Id,
		Status:  BountyStateNew,
	}

	button := coinbase.GetButton(&coinbase.ButtonRequest{
		Name:             "Abhi Yerra", // TODO: This should be the Issue Id
		Type:             "donation",
		PriceString:      "10.00",
		PriceCurrencyIso: "USD",
		Custom:           strconv.Itoa(user_id),
		Style:            "donation_large",
		VariablePrice:    true,
		ChoosePrice:      true,
		// TODO: Include some prices
	})

	if button.Response.Success {
		b.CoinbaseButtonCode = button.Response.Button.Code
	}

	if err := DbMap.Insert(b); err != nil {
		log.Fatal(err)
	}

	return
}

func FindBountyByCoinbaseButtonCode(coinbase_code string) (b *Bounty) {
	err := DbMap.SelectOne(b, "SELECT * FROM bounties WHERE coinbase_button_code = $1", coinbase_code)
	if err != nil {
		log.Printf("FindBountyByCoinbaseCode failed %v\n", err)
		return nil
	}

	return
}

func (b *Bounty) UpdateCoinbaseInfo(order coinbase.Order) {
	b.CoinbaseOrderId = order.Id
	b.CoinbaseTotalBtc = order.TotalBtc.Cents
	b.CoinbaseCurrencyIso = order.TotalBtc.CurrencyIso

	count, err := DbMap.Update(b)
	if err != nil {
		log.Printf("FindBountyByCoinbaseCode failed %v\n", err)
	}

	log.Println("Rows updated:", count)
}
