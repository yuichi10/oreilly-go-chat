package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	fb "github.com/huandu/facebook"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

type authHandler struct {
	next http.Handler
}

func getGoogleOauthConfig() *oauth2.Config {
	fmt.Println(os.Getenv("GOOGLE_OAUTH_CLIENT_ID"))
	fmt.Println(os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"))
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "https://localhost:8080/auth/callback/google",
		Scopes:       []string{"https://picasaweb.google.com/data/", "https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
	}
}

func getFacebookOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("FACEBOOK_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("FACEBOOK_OAUTH_CLIENT_SECRET"),
		Endpoint:     facebook.Endpoint,
		RedirectURL:  "https://localhost:8080/auth/callback/facebook",
		Scopes:       []string{"user_about_me", "email"},
	}
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// 未認証
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(err.Error())
	} else {
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func Provider(provider string) *oauth2.Config {
	switch provider {
	case "google":
		return getGoogleOauthConfig()
	case "github":
		return nil
	case "facebook":
		return getFacebookOauthConfig()
	default:
		return nil
	}
}

func callbackHandler(w http.ResponseWriter, r *http.Request, provider string) {
	switch provider {
	case "google":
		code := r.FormValue("code")
		conf := Provider(provider)
		token, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗しました。")
		}

		client := conf.Client(oauth2.NoContext, token)
		userinfo, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo")
		if err != nil {
			log.Fatal("ユーザー情報の取得に失敗しました。")
		}
		defer userinfo.Body.Close()
		uinfo, _ := ioutil.ReadAll(userinfo.Body)
		fmt.Println(string(uinfo))

	case "github":
	case "facebook":
		code := r.FormValue("code")
		conf := Provider(provider)
		token, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗しました。")
		}
		if token.Valid() == false {
			panic(errors.New("vaild token"))
		}
		client := conf.Client(oauth2.NoContext, token)
		session := &fb.Session{
			Version:    "v2.8",
			HttpClient: client,
		}
		res, err := session.Get("/me?fields=id,name,email", nil)
		if err != nil {
			panic(err)
		}

		fmt.Println(res)
	default:
	}
}

// loginHandlerはサードパーティーのログインの処理を受け持ちます
// パスの形式: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		conf := Provider(provider)
		loginURL := conf.AuthCodeURL("test")
		w.Header().Set("Location", loginURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		callbackHandler(w, r, provider)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sには非対応です", action)
	}
}
