package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	route.PathPrefix("/assets/").Handler(http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	route.HandleFunc("/", Home).Methods("GET")
	route.HandleFunc("/contact", Contact).Methods("GET")
	route.HandleFunc("/addMyProject", AddMyProject).Methods("GET")
	route.HandleFunc("/addMP", AddMP).Methods("POST")
	route.HandleFunc("/myProjectDetail/{index}", MyProjectDetail).Methods("GET")
	route.HandleFunc("/deleteMP/{index}", deleteMP).Methods("GET")

	fmt.Println("Server Running")
	http.ListenAndServe("localhost:5000", route)

}

func Home(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("error : " + err.Error()))
		return
	}

	response := map[string]interface{}{
		"MPs": dataMP,
	}

	tmpl.Execute(w, response)

}
func Contact(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("error : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)

}
func AddMyProject(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/addMyProject.html")

	if err != nil {
		w.Write([]byte("error : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)

}

type MP struct {
	Title       string
	Description string
	startDate   string
	endDate     string
	Duration    string
}

var dataMP = []MP{}

func AddMP(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var title = r.PostForm.Get("title")
	var description = r.PostForm.Get("description")
	var startDate = r.PostForm.Get("startDate")
	var endDate = r.PostForm.Get("endDate")

	layout := "2006-01-02"
	start_date, _ := time.Parse(layout, startDate)
	end_date, _ := time.Parse(layout, endDate)

	hours := end_date.Sub(start_date).Hours()
	days := hours / 24

	var duration string

	if days > 0 {
		duration = strconv.FormatFloat(days, 'f', 0, 64) + " days"
	}

	var newMP = MP{
		Title:       title,
		startDate:   startDate,
		endDate:     endDate,
		Duration:    duration,
		Description: description,
	}

	dataMP = append(dataMP, newMP)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func MyProjectDetail(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/myProjectDetail.html")

	if err != nil {
		w.Write([]byte("error : " + err.Error()))
		return
	}

	var MPDetail = MP{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, data := range dataMP {

		if i == index {
			MPDetail = MP{
				Title:       data.Title,
				startDate:   data.startDate,
				endDate:     data.endDate,
				Description: data.Description,
			}
		}
	}

	data := map[string]interface{}{
		"MP": MPDetail,
	}

	tmpl.Execute(w, data)

}

func deleteMP(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	dataMP = append(dataMP[:index], dataMP[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}

// func editMP(w http.ResponseWriter, r *http.Request)  {
// 	index, _ := strconv.Atoi(mux.Vars(r)["index"])

// }

// func duration(w http.ResponseWriter, r *http.Request)  {
// 	t1 := time.Date(startDate)
// 	t2 := startDate

// 	hs := t1.Sub(t2).Hours()

// 	hs, mf := math.Modf(hs)
// 	ms := mf * 60

// 	ms, sf := math.Modf(ms)
// 	ss := sf * 60

// }
