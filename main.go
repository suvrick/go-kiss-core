package main

import (
	"log"
	"net/http"
	"text/template"
)

func settingHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/setting" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./frontend/templete/pages/setting.html",
		"./frontend/templete/layouts/base.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, Menu{
		IsMenu:   true,
		PathName: "/setting",
	})

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func packetsHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/packets" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./frontend/templete/pages/packets.html",
		"./frontend/templete/layouts/base.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, Menu{
		IsMenu:   true,
		PathName: "/packets",
	})

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./frontend/templete/pages/home.html",
		"./frontend/templete/layouts/base.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, Menu{
		IsMenu:   true,
		PathName: "/",
	})

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

type Menu struct {
	IsMenu   bool
	PathName string
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)

	mux.HandleFunc("/packets", packetsHandler)

	mux.HandleFunc("/setting", settingHandler)

	fileServer := http.FileServer(http.Dir("./frontend/"))

	mux.Handle("/frontend/", http.StripPrefix("/frontend", fileServer))

	http.ListenAndServe(":8080", mux)

}
