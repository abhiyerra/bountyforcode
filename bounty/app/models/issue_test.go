package bountyforcode

import (
	"testing"
)

func TestNewIssue(t *testing.T) {
	t.Errorf("Not implemented")
	// Create a new issue

	// Create a new issue that already exists

	// Create a
}

func TestFindIssue(t *testing.T) {
	setupTestDb()

	issue := FindIssue("1")

	if issue == nil {
		t.Fatalf("Issue is nil")
	}

	if issue.Hoster != "github" {
		t.Errorf("Hoster is not set")
	}

	if issue.Project != "abhiyerra" {
		t.Errorf("Project is not set")
	}

	if issue.Repo != "feedbackjs" {
		t.Errorf("Repo is not set")
	}

	if issue.Identifier != "2" {
		t.Errorf("Identifier is not set")
	}

	if issue := FindIssue("2"); issue != nil {
		t.Fatalf("Issue is not nil")
	}

}

func TestFindAllIssues(t *testing.T) {
	setupTestDb()

	issues := FindAllIssues()

	if len(issues) < 1 {
		t.Errorf("No issues found")
	}
}

func TestFindProjectIssues(t *testing.T) {
	setupTestDb()

	issues := FindProjectIssues("abhiyerra")
	if len(issues) != 1 {
		t.Errorf("No project issues found")
	}

	issues = FindProjectIssues("abcd")
	if len(issues) != 0 {
		t.Errorf("Returned invalid issues")
	}
}
