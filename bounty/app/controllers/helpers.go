package bountyforcode

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	. "github.com/abhiyerra/bountyforcode/bounty/app/models"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

var (
	HtmlDir        string
	Store          *sessions.CookieStore
	SecretStoreKey string
)

type StatusResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func InitSessionStore() {
	if SecretStoreKey == "" {
		log.Fatal("SecretStoreKey isn't set")
	}

	gob.Register(&User{})

	Store = sessions.NewCookieStore([]byte("something-very-secret"))
}

func GetSessionUser(r *http.Request) (u *User) {
	session, _ := Store.Get(r, "user")
	user := session.Values["User"]

	log.Printf("session %v", user)

	u, ok := user.(*User)
	if !ok {
		return nil
	}

	return
}

func SetSessionUserId(w http.ResponseWriter, r *http.Request, user *User) {
	session, _ := Store.Get(r, "user")
	session.Values["User"] = user
	session.Save(r, w)
}

func RenderJson(w http.ResponseWriter, page interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Printf("%v", page)

	b, err := json.Marshal(page)
	if err != nil {
		log.Println("error:", err)
		fmt.Fprintf(w, "")
	}

	w.Write(b)
}
