package main

import (
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {

	username := session.GetString(r.Context(), "user")
	if username == "" {
		http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
		return
	}

	send_home(w, r, username)

}

type Cell struct {
	Present bool
	Name    string
	Code    string
	Venue   string
}

type Schedule struct {
	Name string
	Data map[int]Cell
}

func send_home(w http.ResponseWriter, r *http.Request, person string) {

	friend := get_friends(person)
	day := r.FormValue("day")

	var schedules []Schedule

	schedules = append(schedules, get_schd(day, person))

	// append schedule of friends
	for _, j := range friend {
		schedules = append(schedules, get_schd(day, j))
	}

	tdata := struct {
		Username  string
		Schedules []Schedule
		Day       string
	}{Username: person, Schedules: schedules, Day: day}
	err := templates.ExecuteTemplate(w, "home.html", tdata)
	if err != nil {
		log.Fatalln(err)
	}
}

func get_schd(day string, person string) Schedule {

	MON := map[string]int{
		"A1":  0,
		"L1":  0,
		"F1":  1,
		"L2":  1,
		"D1":  2,
		"L3":  2,
		"TB1": 3,
		"L4":  3,
		"TG1": 4,
		"L5":  4,
		"S11": 5,
		"L6":  5,
		"A2":  6,
		"L31": 6,
		"F2":  7,
		"L32": 7,
		"D2":  8,
		"L33": 8,
		"TB2": 9,
		"L34": 9,
		"TG2": 10,
		"L35": 10,
		"S3":  11,
		"L36": 11,
	}

	TUE := map[string]int{
		"B1":   0,
		"L7":   0,
		"G1":   1,
		"L8":   1,
		"E1":   2,
		"L9":   2,
		"TC1":  3,
		"L10":  3,
		"TAA1": 4,
		"L11":  4,
		"":     5,
		"L12":  5,
		"B2":   6,
		"L37":  6,
		"G2":   7,
		"L38":  7,
		"E2":   8,
		"L39":  8,
		"TC2":  9,
		"L40":  9,
		"TAA2": 10,
		"L41":  10,
		"S1":   11,
		"L42":  11,
	}

	WED := map[string]int{
		"C1":   0,
		"L13":  0,
		"A1":   1,
		"L14":  1,
		"F1":   2,
		"L15":  2,
		"D1":   3,
		"L16":  3,
		"TBB1": 4,
		"L17":  4,
		"":     5,
		"L18":  5,
		"C2":   6,
		"L43":  6,
		"A2":   7,
		"L44":  7,
		"F2":   8,
		"L45":  8,
		"D2":   9,
		"L46":  9,
		"TBB2": 10,
		"L47":  10,
		"S4":   11,
		"L48":  11,
	}

	THU := map[string]int{
		"D1":   0,
		"L19":  0,
		"B1":   1,
		"L20":  1,
		"G1":   2,
		"L21":  2,
		"TE1":  3,
		"L22":  3,
		"TCC1": 4,
		"L23":  4,
		"":     5,
		"L24":  5,
		"D2":   6,
		"L49":  6,
		"B2":   7,
		"L50":  7,
		"G2":   8,
		"L51":  8,
		"E2":   9,
		"L52":  9,
		"TCC2": 10,
		"L53":  10,
		"S2":   11,
		"L54":  11,
	}

	FRI := map[string]int{
		"E1":   0,
		"L25":  0,
		"C1":   1,
		"L26":  1,
		"TA1":  2,
		"L27":  2,
		"TF1":  3,
		"L28":  3,
		"TDD1": 4,
		"L29":  4,
		"S15":  5,
		"L30":  5,
		"E2":   6,
		"L55":  6,
		"C2":   7,
		"L56":  7,
		"TA2":  8,
		"L57":  8,
		"TF2":  9,
		"L58":  9,
		"TDD2": 10,
		"L59":  10,
		"":     11,
		"L60":  11,
	}

	courses, err := get_user_courses(person)
	if err != nil {
		log.Fatalln(err)
	}

	Schd := map[int]Cell{
		0:  {},
		1:  {},
		2:  {},
		3:  {},
		4:  {},
		5:  {},
		6:  {},
		7:  {},
		8:  {},
		9:  {},
		10: {},
		11: {},
	}

	for _, j := range courses {
		b := j.X
		c := j.TX
		d := j.TXX
		var slt [3]string
		if b {
			slt[0] = j.Slot
		}

		if c {
			slt[1] = "T" + j.Slot
		}

		if d {
			slt[1] = "T" + string(j.Slot[0]) + j.Slot
		}

		if day == "MON" {
			for _, l := range slt {
				index, ok := MON[l]
				if ok {
					Schd[index] = Cell{
						Name:    j.Name,
						Venue:   j.Venue,
						Code:    j.Code,
						Present: true,
					}

				}
			}

		}

		if day == "TUE" {
			for _, l := range slt {
				index, ok := TUE[l]
				if ok {
					Schd[index] = Cell{
						Name:    j.Name,
						Venue:   j.Venue,
						Code:    j.Code,
						Present: true,
					}

				}
			}

		}

		if day == "WED" {
			for _, l := range slt {
				index, ok := WED[l]
				if ok {

					Schd[index] = Cell{
						Name:    j.Name,
						Venue:   j.Venue,
						Code:    j.Code,
						Present: true,
					}

				}
			}

		}

		if day == "THU" {
			for _, l := range slt {
				index, ok := THU[l]
				if ok {
					Schd[index] = Cell{
						Name:    j.Name,
						Venue:   j.Venue,
						Code:    j.Code,
						Present: true,
					}
				}
			}

		}

		if day == "FRI" {
			for _, l := range slt {
				index, ok := FRI[l]
				if ok {
					Schd[index] = Cell{
						Name:    j.Name,
						Venue:   j.Venue,
						Code:    j.Code,
						Present: true,
					}

				}
			}

		}
	}

	labcourses, err := get_user_labcourses(person)
	if err != nil {
		log.Fatalln(err)
	}

	for _, j := range labcourses {

		if day == "MON" {
			index, ok := MON["L"+strconv.Itoa((j.Slot-1)*2+1)]
			if ok {
				Schd[index] = Cell{
					Name:    j.Name,
					Venue:   j.Venue,
					Code:    j.Code,
					Present: true,
				}
				Schd[index+1] = Cell{
					Name:    j.Name,
					Venue:   j.Venue,
					Code:    j.Code,
					Present: true,
				}
			}
		}

		if day == "TUE" {
			index, ok := TUE["L"+strconv.Itoa((j.Slot-1)*2+1)]
			if ok {
				Schd[index] = Cell{
					Name:    j.Name,
					Venue:   j.Venue,
					Code:    j.Code,
					Present: true,
				}
				Schd[index+1] = Cell{
					Name:    j.Name,
					Venue:   j.Venue,
					Code:    j.Code,
					Present: true,
				}
			}
		}

		if day == "WED" {
			index, ok := WED["L"+strconv.Itoa((j.Slot-1)*2+1)]
			if ok {
				Schd[index] = Cell{
					Name:    j.Name,
					Venue:   j.Venue,
					Code:    j.Code,
					Present: true,
				}
				Schd[index+1] = Cell{
					Name:    j.Name,
					Venue:   j.Venue,
					Code:    j.Code,
					Present: true,
				}
			}
		}

		if day == "THU" {
			index, ok := THU["L"+strconv.Itoa((j.Slot-1)*2+1)]
			if ok {
				Schd[index] = Cell{
					Name:    j.Name,
					Venue:   j.Venue,
					Code:    j.Code,
					Present: true,
				}
				Schd[index+1] = Cell{
					Name:    j.Name,
					Venue:   j.Venue,
					Code:    j.Code,
					Present: true,
				}
			}
		}

		if day == "FRI" {
			index, ok := FRI["L"+strconv.Itoa((j.Slot-1)*2+1)]
			if ok {
				Schd[index] = Cell{
					Name:    j.Name,
					Venue:   j.Venue,
					Code:    j.Code,
					Present: true,
				}
				Schd[index+1] = Cell{
					Name:    j.Name,
					Venue:   j.Venue,
					Code:    j.Code,
					Present: true,
				}
			}
		}
	}
	return Schedule{
		Name: person,
		Data: Schd,
	}
}
