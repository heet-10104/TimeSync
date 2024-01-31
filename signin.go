package main

import (
	"log"
	"net/http"
)

func send_sigin(w http.ResponseWriter, r *http.Request, failed bool) {
	tdata := struct {
		AppName string
		Failed  bool
	}{AppName: AppName, Failed: failed}
	err := templates.ExecuteTemplate(w, "signin.html", tdata)
	if err != nil {
		log.Fatalln(err)
	}
}

func signin(w http.ResponseWriter, r *http.Request) {
	if session.Exists(r.Context(), "user") {
		http.Redirect(w, r, "/edit", http.StatusTemporaryRedirect)
	}
	if r.Method == "GET" {
		send_sigin(w, r, false)
	} else if r.Method == "POST" {
		r.ParseForm()
		//username := r.Form["username"][0]
		username := r.Form.Get("username")
		password := r.Form.Get("password")

		row := database.QueryRow("SELECT password FROM users WHERE username=?", username)

		var psw string
		if err := row.Scan(&psw); err == nil {
			if CheckPasswordHash(password, psw) {
				session.Put(r.Context(), "user", username)
				http.Redirect(w, r, "/edit", http.StatusSeeOther)
				return
			}
		}
		send_sigin(w, r, true)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}
