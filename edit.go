package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func edit(w http.ResponseWriter, r *http.Request) {
	// TODO check duplicate slots of user
	username := session.GetString(r.Context(), "user")
	if username == "" {
		http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
		return
	}

	if r.Method == "POST" {
		var err string
		switch r.FormValue("type") {
		case "theory":
			err = save_thoery_course(r, username)
		case "lab":
			err = save_lab_course(r, username)
		}
		if err != "" {
			send_edit(w, r, err)
			return
		}
	}
	send_edit(w, r, "")
}

func save_thoery_course(r *http.Request, username string) string {
	course, err_msg := parse_course(r)
	if course.Id == "+" {
		if err_msg != "" {
			return err_msg
		}
		res, err := database.Exec("INSERT INTO courses(user, name, code, slot, x, tx, txx, venue) VALUES(?, ?, ?, ?, ?, ?, ?, ?)", username, course.Name, course.Code, course.Slot, course.X, course.TX, course.TXX, course.Venue)
		log.Println(res)
		if err != nil {
			log.Fatalln(err)
		}
		return ""
	}
	id, err := strconv.Atoi(course.Id)
	if err != nil {
		// Fail silently no need for an error message
		log.Println(err)
		return ""
	}
	if r.FormValue("action") == "Edit" {
		_, err = database.Exec("UPDATE courses SET user=?, name=?, slot=?, x=?, tx=?, txx=?, venue=? WHERE id=? AND user=?", username, course.Name, course.Slot, course.X, course.TX, course.TXX, course.Venue, id, username)
		if err != nil {
			log.Fatalln(err)
		}
	} else if r.FormValue("action") == "Delete" {
		_, err = database.Exec("DELETE FROM courses WHERE id=? AND user=?", id, username)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return ""
}

func save_lab_course(r *http.Request, username string) string {
	course, err_msg := parse_lab_course(r)
	if course.Id == "+" {
		if err_msg != "" {
			return err_msg
		}
		res, err := database.Exec("INSERT INTO lab_courses(user, name, code, slot, venue) VALUES(?, ?, ?, ?, ?)", username, course.Name, course.Code, course.Slot, course.Venue)
		log.Println(res)
		if err != nil {
			log.Fatalln(err)
		}
		return ""
	}
	id, err := strconv.Atoi(course.Id)
	if err != nil {
		// Fail silently no need for an error message
		log.Println(err)
		return ""
	}
	if r.FormValue("action") == "Edit" {
		_, err = database.Exec("UPDATE lab_courses SET user=?, name=?, slot=?, venue=? WHERE id=? AND user=?", username, course.Name, course.Slot, course.Venue, id, username)
		if err != nil {
			log.Fatalln(err)
		}
	} else if r.FormValue("action") == "Delete" {
		_, err = database.Exec("DELETE FROM lab_courses WHERE id=? AND user=?", id, username)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return ""
}

func parse_course(r *http.Request) (Course, string) {
	var course Course
	course.Id = r.FormValue("id")
	course.Name = r.FormValue("name")
	course.Slot = r.FormValue("slot")
	course.Code = r.FormValue("code")
	course.Venue = r.FormValue("venue")

	if course.Name == "" {
		return course, "Course Name cannot be empty"
	}

	if len(course.Code) == 0 {
		return course, "code cannot be empty"
	}

	if r.FormValue("type1") == "on" {
		course.X = true
	}
	if r.FormValue("type2") == "on" {
		course.TX = true
	}
	if r.FormValue("type3") == "on" {
		course.TXX = true
	}
	if !(course.X || course.TX || course.TXX) {
		return course, "you must select at least one slot"
	}
	course.Venue = strings.TrimSpace(course.Venue)
	if len(course.Venue) == 0 {
		return course, "venue cannot be empty"
	}
	return course, ""
}

func parse_lab_course(r *http.Request) (LabCourse, string) {
	var course LabCourse
	var err error
	course.Id = r.FormValue("id")
	course.Name = r.FormValue("name")
	course.Slot, err = strconv.Atoi(r.FormValue("slot"))
	course.Code = r.FormValue("code")
	course.Venue = r.FormValue("venue")
	if err != nil {
		return course, ""
	}
	fmt.Print(r.Form)

	if course.Name == "" {
		return course, "Course Name cannot be empty"
	}

	if len(course.Code) == 0 {
		return course, "code cannot be empty"
	}
	course.Venue = strings.TrimSpace(course.Venue)
	if len(course.Venue) == 0 {
		return course, "venue cannot be empty"
	}
	return course, ""
}

type Course struct {
	Id    string
	Name  string
	Code  string
	Slot  string
	X     bool
	TX    bool
	TXX   bool
	Venue string
}

type LabCourse struct {
	Id    string
	Name  string
	Code  string
	Slot  int
	Venue string
}

func send_edit(w http.ResponseWriter, r *http.Request, err_mesg string) {
	username := session.GetString(r.Context(), "user")
	if username == "" {
		http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
	}
	courses, err := get_user_courses(username)
	if err != nil {
		log.Fatalln(err)
	}
	labcourses, err := get_user_labcourses(username)
	if err != nil {
		log.Fatalln(err)
	}
	tdata := struct {
		ErrMsg     string
		AppName    string
		Courses    []Course
		LabCourses []LabCourse
		LabCounts  []LC
	}{ErrMsg: err_mesg, AppName: AppName, Courses: courses, LabCourses: labcourses, LabCounts: lab_counts(30)}
	err = templates.ExecuteTemplate(w, "edit.html", tdata)
	if err != nil {
		log.Fatalln(err)
	}
}

type LC struct {
	Id   int
	Disp string
}

func lab_counts(n int) []LC {
	lc := make([]LC, n)

	for i := 0; i < n; i++ {
		lc[i].Id = i + 1
		lc[i].Disp = "L" + strconv.Itoa((i)*2+1) + "+L" + strconv.Itoa((i+1)*2)
	}
	return lc
}
