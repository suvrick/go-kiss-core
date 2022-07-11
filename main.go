package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/suvrick/go-kiss-core/packets/meta"
	"github.com/suvrick/go-kiss-core/ws"

	_ "net/http/pprof"
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

	meta := meta.NewMeta()
	if meta.Error != nil {
		log.Fatalln(meta.Error.Error())
	}

	config := ws.GetDefaultGameSocketConfig()

	gs := ws.NewGameSocket(config)
	gs.Run()
	login_params := []interface{}{1000015, 32, 4, "200514254f3678c2f79cb18760ba048d", 0, ""}
	gs.Send(4, login_params)

	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)

	mux.HandleFunc("/packets", packetsHandler)

	mux.HandleFunc("/setting", settingHandler)

	fileServer := http.FileServer(http.Dir("./frontend/"))

	mux.Handle("/frontend/", http.StripPrefix("/frontend", fileServer))

	http.ListenAndServe(":8080", nil)

}
