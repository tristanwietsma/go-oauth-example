package main

import (
	"encoding/json"
	"fmt"
	"github.com/GolangDorks/endpoint"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

var (
	clientID     = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
)

func state() string {
	// generate and save a unique token here...
	return "abc123"
}

func checkState(token string) bool {
	// validate token is legit
	return true
}

func getAccess(code string) string {
	var v url.Values
	v = make(url.Values)
	v.Set("client_id", clientID)
	v.Set("client_secret", clientSecret)
	v.Set("code", code)
	v.Set("redirect_uri", "http://localhost:8080/oauth")
	resp, _ := http.Post("https://github.com/login/oauth/access_token?"+v.Encode(), "", nil)
	body, _ := ioutil.ReadAll(resp.Body)
	values, _ := url.ParseQuery(string(body))
	return values.Get("access_token")
}

func getID(code string) interface{} {
	token := getAccess(code)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Add("Authorization", "token "+token)
	resp, _ := client.Do(req)
	var data map[string]interface{}
	data = make(map[string]interface{})
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&data)
	return data["login"]
}

var login = endpoint.Endpoint{
	Path:   "/login",
	Method: "GET",
	Before: []endpoint.Middleware{},
	Control: func(ctx endpoint.Context) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

			var v url.Values
			v = make(url.Values)
			v.Set("client_id", clientID)
			v.Set("redirect_uri", "http://localhost:8080/oauth")
			v.Set("scope", "user")
			v.Set("state", state())

			u := url.URL{
				Scheme:   "https",
				Host:     "github.com",
				Path:     "/login/oauth/authorize",
				RawQuery: v.Encode(),
			}

			page := `<!doctype html><html><head><title>Golang OAuth</title></head><body><a href="` +
				u.String() + `">Login with GitHub</a></body></html>`
			w.Write([]byte(page))
		}
	},
}

var oauth = endpoint.Endpoint{
	Path:   "/oauth",
	Method: "GET",
	Before: []endpoint.Middleware{},
	Control: func(ctx endpoint.Context) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			vals := r.URL.Query()
			token := vals.Get("state")
			if checkState(token) {
				code := vals.Get("code")
				id := getID(code)
				w.Write([]byte(fmt.Sprintf("hello %v", id)))
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		}
	},
}

func main() {
	router := httprouter.New()
	router.Handle(login.Method, login.Path, login.Handler())
	router.Handle(oauth.Method, oauth.Path, oauth.Handler())
	log.Fatal(http.ListenAndServe(":8080", router))
}
