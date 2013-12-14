package bountyforcode

import (
	"github.com/abhiyerra/coinbase"
	"log"
)

const (
	BountyStateNew       string = "new"
	BountyStatePaid      string = "paid"
	BountyStateClosed    string = "closed"
	BountyStateCancelled string = "cancelled"
)

type Bounty struct {
	Id                 int     `db:"id"`
	UserId             string  `db:"user_id"`
	IssueId            string  `db:"issue_id"`
	Amount             float32 `db:"amount"`
	CoinbaseButtonCode string  `db:"coinbase_button_code"`
	Status             string  `db:"status"`
}

func NewBounty(issue *Issue, user_id string) (b *Bounty) {
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
		Custom:           user_id,
		Style:            "donation_large",
		VariablePrice:    true,
		ChoosePrice:      true,
		// TODO: Include some prices
	})

	if button.Response.Success {
		b.CoinbaseButtonCode = button.Response.Button.Code
	}

	err := DbMap.Insert(b)
	if err != nil {
		log.Fatal(err)
	}

	return b

}

func BountiesOpen() {
	// SELECT * FROM issues INNER JOIN bounties ON bounties.issue_id = issue.id WHERE bounty_state in ('open', 'paid')
}

func BountiesRecent() {
	// SELECT * FROM bounties ORDER BY created_at DESC limit 10
}
