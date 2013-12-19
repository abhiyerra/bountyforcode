package bountyforcode

import (
	"encoding/json"
	"fmt"
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

	Store = sessions.NewCookieStore([]byte("something-very-secret"))
}

func GetSessionUserId(r *http.Request) string {
	session, _ := Store.Get(r, "user")
	user_id := session.Values["UserId"]

	if str, ok := user_id.(string); ok {
		return str
	} else {
		return ""
	}
}

func SetSessionUserId(w http.ResponseWriter, r *http.Request, user_id string) {
	session, _ := Store.Get(r, "user")
	session.Values["UserId"] = user_id
	session.Save(r, w)
}

func RenderJson(w http.ResponseWriter, page interface{}) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("%v", page)

	b, err := json.Marshal(page)
	if err != nil {
		log.Println("error:", err)
		fmt.Fprintf(w, "")
	}

	w.Write(b)
}
