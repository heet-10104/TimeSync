package main

import (
	"log"
	"net/http"
)

func request(w http.ResponseWriter, r *http.Request) {
	username := session.GetString(r.Context(), "user")
	if username == "" {
		http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
		return
	}

	if r.Method == "POST" {
		// either request accept reject or remove
		switch r.FormValue("action") {
		case "Request":
			errmsg := make_request(username, r.FormValue("user"))
			send_request(w, r, errmsg)
			return
		case "Accept":
			errmsg := accept_request(username, r.FormValue("user"))
			send_request(w, r, errmsg)
			return
		case "Reject":
			errmsg := reject_request(username, r.FormValue("user"))
			send_request(w, r, errmsg)
			return
		case "Remove":
			errmsg := remove_friend(username, r.FormValue("user"))
			send_request(w, r, errmsg)
			return
		}
	}
	send_request(w, r, "")
}

func make_request(sender string, reciever string) string {
	var exists bool
	err := database.QueryRow("SELECT EXISTS(SELECT 1 FROM request WHERE (sender=? and reciever=?) or (sender=? and reciever=?))", sender, reciever, reciever, sender).Scan(&exists)
	if err != nil {
		log.Fatalln(err)
	}
	if exists {
		return "Request already sent"
	}

	_, err = database.Exec("INSERT INTO request(sender, reciever, status) VALUES (?,?,0)", sender, reciever)
	if err != nil {
		log.Fatalln(err)
	}
	return ""
}
func accept_request(username string, sender string) string {
	_, err := database.Exec("UPDATE request SET status=1 WHERE sender=? AND reciever=?", sender, username)
	if err != nil {
		log.Fatalln(err)
	}
	return ""
}

func reject_request(username string, sender string) string {
	_, err := database.Exec("DELETE FROM request WHERE sender=? AND reciever=?", sender, username)
	if err != nil {
		log.Fatalln(err)
	}
	return ""
}

func remove_friend(username string, friend string) string {
	_, err := database.Exec("DELETE FROM request WHERE (sender=? AND reciever=?) OR (sender=? AND reciever=?)", username, friend, friend, username)
	if err != nil {
		log.Fatalln(err)
	}
	return ""
}

func send_request(w http.ResponseWriter, r *http.Request, errmsg string) {
	username := session.GetString(r.Context(), "user")
	if username == "" {
		http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
		return
	}

	pending := get_pending(username)
	friends := get_friends(username)

	data := struct {
		AppName  string
		Username string
		Errmsg   string
		Pending  []string
		Friends  []string
	}{
		AppName:  AppName,
		Username: username,
		Errmsg:   errmsg,
		Pending:  pending,
		Friends:  friends,
	}
	err := templates.ExecuteTemplate(w, "request.html", data)
	if err != nil {
		log.Fatalln(err)
	}
}

func get_pending(username string) []string {
	rows, err := database.Query("SELECT sender FROM request WHERE reciever=? AND status=0", username)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	var pending []string
	for rows.Next() {
		var sender string
		err = rows.Scan(&sender)
		if err != nil {
			log.Fatalln(err)
		}
		pending = append(pending, sender)
	}
	return pending
}

