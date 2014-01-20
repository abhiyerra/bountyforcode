package bountyforcode

import (
	"github.com/octokit/go-octokit/octokit"
	"log"
	"time"
)

type User struct {
	Id                int       `db:"id" json:"id"`
	GithubUsername    string    `db:"github_username" json:"github_username"`
	GithubAccessToken string    `db:"github_access_token" json:"access_token"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
	//	ExpiresAt      string `db:"expires_at" json:"expires_at"`

}

func GithubUser(access_token string) (user *octokit.User) {
	client := octokit.NewClient(octokit.TokenAuth{access_token})

	url, _ := octokit.CurrentUserURL.Expand(octokit.M{})
	user, _ = client.Users(url).One()

	log.Printf("url %s %v", user.Login, user)

	return
}

func NewUser(github_username, access_token string) (u *User) {
	if access_token == "" {
		log.Printf("No access_token to add to db\n")
		return
	}

	u = &User{
		GithubUsername:    github_username,
		GithubAccessToken: access_token,
	}

	err := DbMap.SelectOne(u, "SELECT * FROM users WHERE github_username = $1", github_username)
	if err != nil {
		if err = DbMap.Insert(u); err != nil {
			log.Printf("Couldn't add github user %s %v", github_username, err)
			return nil
		}
	}

	u.GithubAccessToken = access_token
	DbMap.Update(u)

	return
}

func FindUser(id string) (u *User) {
	obj, err := DbMap.Get(User{}, id)
	if err != nil {
		log.Println(err)
	}
	u = obj.(*User)

	return
}
