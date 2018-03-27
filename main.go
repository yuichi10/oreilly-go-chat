package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/yuichi10/trace"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "application's address")
	flag.Parse()
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/token", jwtHandler)
	http.HandleFunc("/validate", jwtValidator)
	http.Handle("/room", r)
	// チャットルームを開始
	go r.run()
	// webサーバーを起動
	log.Println("start web server. PORT: ", *addr)
	if err := http.ListenAndServeTLS(*addr, "./ssl.crt", "./ssl.key", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
