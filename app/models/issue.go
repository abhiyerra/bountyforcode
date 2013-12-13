package bountyforcode

import (
	"database/sql"
	"fmt"
	"github.com/abhiyerra/coinbase"
	"github.com/abhiyerra/scalpy"
	"log"
	"strings"
)

type Issue struct {
	Id string

	Project    string
	Repo       string
	Identifier string

	OriginalUrl string
	CreatedAt   string

	CoinbaseButtonCode string
}

func (i *Issue) GetCoinbaseButton() {

}

func FindIssue(id string) (i *Issue) {
	i = new(Issue)

	err := Db.QueryRow(`SELECT id, project, repo, identifier, coinbase_button_code FROM issues WHERE id = $1`, id).Scan(&i.Id, &i.Project, &i.Repo, &i.Identifier, &i.CoinbaseButtonCode)

	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("IssueId is %s\n", i.Id)
	}

	return i
}

func FindProjectIssues(project string) (issues []Issue) {
	var rows *sql.Rows
	var err error

	if project == "" {
		rows, err = Db.Query("SELECT id, repo, identifier, coinbase_button_code FROM issues")
	} else {
		rows, err = Db.Query("SELECT id, repo, identifier, coinbase_button_code FROM issues WHERE project = $1", strings.ToLower(project))
	}

	if err != nil {
		log.Fatal(err)
		return nil
	}

	for rows.Next() {
		var issue Issue

		if err := rows.Scan(&issue.Id, &issue.Repo, &issue.Identifier, &issue.CoinbaseButtonCode); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", issue)
		issues = append(issues, issue)
	}

	return issues
}

func NewIssue(scalp *scalpy.Scalp) (i *Issue) {
	i = new(Issue)

	project := scalp.Project
	repo := scalp.Repo

	err := Db.QueryRow(`SELECT id FROM issues WHERE hoster = 'github' AND project = $1 AND repo = $2 AND identifier = $3`, project, repo, scalp.IssueId).Scan(&i.Id)
	switch {
	case err == sql.ErrNoRows:

		button := coinbase.GetButton(&coinbase.ButtonRequest{
			Name:             "Abhi Yerra",
			Type:             "donation",
			PriceString:      "10.00",
			PriceCurrencyIso: "USD",
		})
		var coinbase_code string
		if button.Response.Success {
			coinbase_code = button.Response.Button.Code
		}

		Db.QueryRow(`INSERT INTO issues (hoster, project, repo, identifier, coinbase_button_code) VALUES ('github', $1, $2, $3, $4) RETURNING id`, project, repo, scalp.IssueId, coinbase_code).Scan(&i.Id)
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("IssueId is %s\n", i.Id)
	}

	return i
}
