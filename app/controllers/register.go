package bountyforcode

import (
	"code.google.com/p/goauth2/oauth"
	"fmt"
	. "github.com/abhiyerra/bountyforcode/app/models"
	"log"
	"net/http"
)

var (
	GithubClientId     string
	GithubClientSecret string
	GithubRedirectUrl  string
	GithubAuthUrl      = "https://github.com/login/oauth/authorize"
	GithubTokenUrl     = "https://github.com/login/oauth/access_token"
	GithubConfig       *oauth.Config
)

func InitGithub() {
	GithubConfig = &oauth.Config{
		ClientId:     GithubClientId,
		ClientSecret: GithubClientSecret,
		AuthURL:      GithubAuthUrl,
		TokenURL:     GithubTokenUrl,
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Make sure token still exists.

	url := GithubConfig.AuthCodeURL("")
	http.Redirect(w, r, url, 302)
}

func RegisterAuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	if code == "" {
		fmt.Fprintf(w, "Failed to login")
		return
	}

	transport := &oauth.Transport{Config: GithubConfig}
	token, _ := transport.Exchange(code)

	log.Printf("New Token %v\n", token)
	user := NewUser(token.AccessToken)

	session, _ := Store.Get(r, "user")
	session.Values["UserId"] = user.Id
	session.Save(r, w)

	fmt.Fprintf(w, token.AccessToken)
}
