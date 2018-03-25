package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
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
		RedirectURL:  "http://localhost:8080/auth/callback/google",
		Scopes:       []string{"https://picasaweb.google.com/data/"},
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
		return nil
	default:
		return nil
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
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sには非対応です", action)
	}
}
