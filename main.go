
package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	_ "github.com/mattn/go-sqlite3"
)

const AppName string = "Project X"

var database *sql.DB
var templates = template.Must(template.ParseGlob("template/*.html"))
var session = scs.New()

func init() {
	//defer database.Close()
	var err error
	database, err = sql.Open("sqlite3", "./database.sqlite3")
	if err != nil {
		log.Fatalln(err)
	}
	session.Lifetime = 6 * time.Hour
}

func sayHi(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hi")
	fmt.Println(r.URL)
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", sayHi)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/signin", signin)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/home", home)
	mux.HandleFunc("/edit", edit)
	mux.HandleFunc("/request", request)
	err := http.ListenAndServe(":8080", session.LoadAndSave(mux))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
