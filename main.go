package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"strconv"
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

func main() {
	host := flag.String("host", "", "address to host the site")
	port := flag.Int("port", 8080, "port to host the site")
	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("/", signin)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/signin", signin)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/home", home)
	mux.HandleFunc("/edit", edit)
	mux.HandleFunc("/request", request)
	err := http.ListenAndServe(*host+":"+strconv.Itoa(*port), session.LoadAndSave(mux))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
