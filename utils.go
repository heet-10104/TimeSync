package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func get_user_courses(username string) ([]Course, error) {
	rows, err := database.Query("SELECT id,name,code,slot,x,tx,txx,venue FROM courses WHERE user=?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var courses []Course

	for rows.Next() {
		var course Course
		if err := rows.Scan(&course.Id, &course.Name, &course.Code, &course.Slot, &course.X, &course.TX, &course.TXX, &course.Venue); err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	return courses, nil
}

func get_user_labcourses(username string) ([]LabCourse, error) {
	rows, err := database.Query("SELECT id,name,code,slot ,venue FROM lab_courses WHERE user=?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var courses []LabCourse

	for rows.Next() {
		var course LabCourse
		if err := rows.Scan(&course.Id, &course.Name, &course.Code, &course.Slot, &course.Venue); err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	return courses, nil
}

func get_friends(username string) []string {
	rows, err := database.Query("SELECT sender FROM request WHERE reciever=? AND status=1", username)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	var friends []string
	for rows.Next() {
		var friend string
		err = rows.Scan(&friend)
		if err != nil {
			log.Fatalln(err)
		}
		friends = append(friends, friend)
	}
	rows2, err := database.Query("SELECT reciever FROM request WHERE sender=? AND status=1", username)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows2.Close()
	for rows2.Next() {
		var friend string
		err = rows2.Scan(&friend)
		if err != nil {
			log.Fatalln(err)
		}
		friends = append(friends, friend)
	}
	seen := make(map[string]bool)
	result := []string{}
	for _, v := range friends {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
