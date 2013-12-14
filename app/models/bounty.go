package bountyforcode

import (
	"github.com/abhiyerra/coinbase"
	"github.com/coopernurse/gorp"
	"log"
)

const (
	New       string = "new"
	Paid      string = "paid"
	Closed    string = "closed"
	Cancelled string = "cancelled"
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
		Status:  New,
	}

	button := coinbase.GetButton(&coinbase.ButtonRequest{
		Name:             "Abhi Yerra",
		Type:             "donation",
		PriceString:      "10.00",
		PriceCurrencyIso: "USD",
	})

	if button.Response.Success {
		b.CoinbaseButtonCode = button.Response.Button.Code
	}

	DbMap := &gorp.DbMap{Db: Db, Dialect: gorp.PostgresDialect{}}
	DbMap.AddTableWithName(Bounty{}, "bounties").SetKeys(true, "Id")

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
