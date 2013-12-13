package bountyforcode

import (
	"database/sql"
	"fmt"
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
}

func FindIssue(id string) (i *Issue) {
	i = new(Issue)

	err := Db.QueryRow(`SELECT id, project, repo, identifier FROM issues WHERE id = $1`, id).Scan(&i.Id, &i.Project, &i.Repo, &i.Identifier)

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
		rows, err = Db.Query("SELECT id, repo, identifier FROM issues")
	} else {
		rows, err = Db.Query("SELECT id, repo, identifier FROM issues WHERE project = $1", strings.ToLower(project))
	}

	if err != nil {
		log.Fatal(err)
		return nil
	}

	for rows.Next() {
		var issue Issue

		if err := rows.Scan(&issue.Id, &issue.Repo, &issue.Identifier); err != nil {
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
		Db.QueryRow(`INSERT INTO issues (hoster, project, repo, identifier) VALUES ('github', $1, $2, $3) RETURNING id`, project, repo, scalp.IssueId).Scan(&i.Id)
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("IssueId is %s\n", i.Id)
	}

	return i
}
