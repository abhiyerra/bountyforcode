package bountyforcode

import (
	"database/sql"
	"github.com/abhiyerra/scalpy"
	"log"
	"strings"
	"time"
)

type Issue struct {
	Id          int            `db:"id" json:"id"`
	Hoster      string         `db:"hoster" json:"hoster"`
	Project     string         `db:"project" json:"project"`
	Repo        string         `db:"repo" json:"repo"`
	Identifier  string         `db:"identifier" json:"identifier"`
	OriginalUrl sql.NullString `db:"original_url" json:"original_url,omitempty"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updated_at"`
}

func FindIssue(id string) (i *Issue) {
	i = &Issue{}

	if err := DbMap.SelectOne(i, "SELECT * FROM issues WHERE id = $1", id); err != nil {
		log.Printf("FindIssue failed %v\n", err)
		return nil
	}

	return
}

func FindAllIssues() (issues []Issue) {
	_, err := DbMap.Select(&issues, "select * from issues")
	if err != nil {
		log.Printf("FindAllIssues failed %v\n", err)
	}

	return
}

func FindProjectIssues(project string) (issues []Issue) {
	if _, err := DbMap.Select(&issues, "SELECT * FROM issues WHERE project = $1", strings.ToLower(project)); err != nil {
		log.Printf("FindAllIssues failed %v\n", err)
	}

	return
}

func NewIssue(scalp *scalpy.Scalp) (i *Issue) {
	i = &Issue{
		Project:    scalp.Project,
		Repo:       scalp.Repo,
		Identifier: scalp.IssueId,
	}

	err := DbMap.SelectOne(&i, "select * from issues where hoster = 'github' AND project = $1 AND repo = $2 AND identifier = $3", i.Project, i.Repo, i.Identifier)
	if err != nil {
		log.Printf("NewIssue failed %v\n", err)
	}

	if err = DbMap.Insert(&i); err != nil {
		log.Printf("NewIssue failed %v\n", err)
	}

	return
}
