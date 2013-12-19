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

	redirect_url := r.FormValue("redirect")
	session, _ := Store.Get(r, "user")
	session.Values["RedirectUrl"] = redirect_url
	session.Save(r, w)

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

	// TODO: Get the user information. Login the user if they already exist in te system.
	if user := NewUser(token.AccessToken); user != nil {
		SetSessionUserId(w, r, user.Id)
	}

	session, _ := Store.Get(r, "user")
	redirect_url := session.Values["RedirectUrl"]
	str, ok := redirect_url.(string)
	if !ok {
		str = "http://bountyforcode.com" // TODO: Make this a const
	}

	http.Redirect(w, r, str, 302)
}

func UserSessionHandler(w http.ResponseWriter, r *http.Request) {
	RenderJson(w, GetSessionUserId(r) != "")
}
