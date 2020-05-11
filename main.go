package main

import (
	"os"
	"log"
	"flag"
	"fmt"
	"strings"
	"net/http"
	"html/template"

	"github.com/kaypee90/safetrans/version"
    "github.com/urfave/negroni"
)

type Page struct{
	Room string
}

func main() {
	var addr = flag.String("addr", "3000", "http service address")
    
	versionFlag := flag.Bool("version", false, "Version")
	serverFlag := flag.Bool("server", false,"Server")
	flag.Parse()

	if *versionFlag {
		fmt.Println("Build Date:", version.BuildDate)
        fmt.Println("Git Commit:", version.GitCommit)
        fmt.Println("Version:", version.Version)
        fmt.Println("Go Version:", version.GoVersion)
        fmt.Println("OS / Arch:", version.OsArch)
		return
	}

	if *serverFlag {
		
		log.Println("Safe Trans Application - Server")
	}
	
	mux := http.NewServeMux()

	hub := newHub()
	go hub.run()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Println("New Request!")
		serveWs(w, r)
	})
	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/chat", ChatHandler)

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(mux)

	message := "Server running on port: "
	result := fmt.Sprintf("%s%s", message, *addr) 
	fmt.Println(result)

	port := os.Getenv("PORT")
	if port == "" {
		port = *addr
	}
	

	http.ListenAndServe(":" + port, n)
	
}


func HomeHandler(w http.ResponseWriter, r *http.Request){
	
		w.Header().Set("Content-Type", "text/html")
		
		if t, err := template.ParseFiles("home.html"); err !=nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}else {
			
			t.Execute(w, nil)
		}
}

func ChatHandler(w http.ResponseWriter, r *http.Request){
	var roomId = ""
	if r.Method == "POST"{
		roomId = strings.TrimSpace(template.HTMLEscapeString(r.FormValue("privateId")))
		if(roomId == ""){
			http.Redirect(w, r, "/", http.StatusBadRequest)
		}
	}else{
		roomId = GenerateRoomId(10)
	}

	w.Header().Set("Content-Type", "text/html")
	
	if t, err := template.ParseFiles("chat.html"); err !=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}else {
		
		page := Page {Room: roomId}
		t.Execute(w, page)
	}
}
