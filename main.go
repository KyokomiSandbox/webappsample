package main

import (
	"log"
	"net/http"
	"fmt"
	"strings"
	"html/template"
	"github.com/kyokomi/webappsample/session"
	_ "github.com/kyokomi/webappsample/memory"
)

var globalSessions *session.Manager
//この後init関数で初期化されます。
func init() {
	fmt.Println("init start")
	var err error
	globalSessions, err = session.NewManager("memory","gosessionid",3600)
	if err != nil {
		fmt.Println(err)
	}
	go globalSessions.GC()
	fmt.Println("init end")
}

func doLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //リクエストを取得するメソッド

	sess := globalSessions.SessionStart(w, r)
	r.ParseForm()

	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/login.tmpl.html")
		t.Execute(w,  sess.Get("username"))
	} else {
		username := template.HTMLEscapeString(r.Form.Get("username"))
		token := r.Form.Get("token")
		if token != "" {
			//tokenの合法性を検証します。
		} else {
			//tokenが存在しなければエラーを出します。
		}

		sess.Set("username", username)

		// r.FormValue["username"]と書くことでr.ParseForm()を省略可能
		fmt.Println("username length:", len(r.Form["username"][0]))
		fmt.Println("username:", username) //サーバ側に出力します。
		fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))

		http.Redirect(w, r, "/", 302)
	}
}

func doCount(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	ct := sess.Get("countnum")
	if ct == nil {
		sess.Set("countnum", 1)
	} else {
		sess.Set("countnum", (ct.(int) + 1))
	}
	t, _ := template.ParseFiles("templates/count.tmpl.html")
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, sess.Get("countnum"))
}

func doSayHelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	http.HandleFunc("/", doSayHelloName)
	http.HandleFunc("/login", doLogin)
	http.HandleFunc("/count", doCount)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

