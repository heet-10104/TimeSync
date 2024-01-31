package main

import (
	"log"
	"net/http"
	"strings"
)

func send_signup(w http.ResponseWriter, r *http.Request, err_msg string) {
	tdata := struct {
		AppName string
		ErrMsg  string
	}{AppName, err_msg}
	err := templates.ExecuteTemplate(w, "signup.html", tdata)
	if err != nil {
		log.Fatalln(err)
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	if session.Exists(r.Context(), "user") {
		http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
	}
	if r.Method == "POST" {
		r.ParseForm()
		username := r.Form.Get("username")
		pass := r.Form.Get("password")
		pass_confirm := r.Form.Get("password_1")
		username = strings.TrimSpace(username)
		if username == "" {
			send_signup(w, r, "Username cannot be empty")
			return
		}
		if pass != pass_confirm {
			send_signup(w, r, "Passwords do not match")
			return
		}
		if pass == "" {
			send_signup(w, r, "Password cannot be empty")
			return
		}

		if len(username) > 25 {
			send_signup(w, r, "username cannot be more than 25")
			return
		}

		var exists bool
		err := database.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=?)", username).Scan(&exists)
		if err != nil {
			log.Fatalln(err)
		}
		if exists {
			send_signup(w, r, "username already exists")
			return
		}

		pass, err = HashPassword(pass)
		if err != nil {
			log.Fatalln(err)
		}

		_, err = database.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, pass)
		if err != nil {
			log.Fatalln(err)
		}
		session.Put(r.Context(), "user", username)
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	send_signup(w, r, "")
	/*var Failed_1 bool
	if password_1 == "" {
		Failed_1 = true
	}
	var Failed bool
	if password != password_1 {
		Failed = true
	} else {
		Failed = false
	}
	tdata := struct {
		AppName  string
		Failed   bool
		Failed_1 bool
	}{AppName: AppName, Failed: Failed, Failed_1: Failed_1}

	err = t.Execute(w, tdata)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = database.Exec("INSERT INTO users(username, password) VALUES(?, ?)", new_username, password)
	if err != nil {
		log.Fatalln(err)
	}*/
}
