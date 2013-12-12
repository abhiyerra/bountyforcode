package bountyforcode

type Bounty struct {
}

func NewBounty() {

}

func BountiesOpen() {
	// SELECT * FROM issues INNER JOIN bounties ON bounties.issue_id = issue.id WHERE bounty_state in ('open', 'paid')
}

func BountiesRecent() {
	// SELECT * FROM bounties ORDER BY created_at DESC limit 10
}
