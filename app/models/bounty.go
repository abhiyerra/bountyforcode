package bountyforcode

import (
	"github.com/abhiyerra/coinbase"
)

type BountyState string

const (
	New       BountyState = "new"
	Paid      BountyState = "paid"
	Closed    BountyState = "closed"
	Cancelled BountyState = "cancelled"
)

type Bounty struct {
	Id      uint
	UserId  uint
	IssueId uint
	Amount  float32

	CoinbaseButtonCode string `db:"coinbase_button_code"`
	Status             BountyState
}

func NewBounty() (b *Bounty) {
	button := coinbase.GetButton(&coinbase.ButtonRequest{
		Name:             "Abhi Yerra",
		Type:             "donation",
		PriceString:      "10.00",
		PriceCurrencyIso: "USD",
	})

	if button.Response.Success {
		b.CoinbaseButtonCode = button.Response.Button.Code
	}

	return b

}

func BountiesOpen() {
	// SELECT * FROM issues INNER JOIN bounties ON bounties.issue_id = issue.id WHERE bounty_state in ('open', 'paid')
}

func BountiesRecent() {
	// SELECT * FROM bounties ORDER BY created_at DESC limit 10
}
