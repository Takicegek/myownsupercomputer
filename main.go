package main

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving index.html to ", r.RemoteAddr)

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		Write404(w, r)
		log.Println("error loading index.html!")
		return
	}

	if client, ok := clientList[r.RemoteAddr]; ok {
		client.LastUpdate = time.Now()
		clientList[r.RemoteAddr] = client
	} else {
		mutex.Lock()
		client := clientInfo{
			NodeID:     "SC-0.1-" + strconv.Itoa(GNodeID),
			Address:    net.ParseIP(r.RemoteAddr),
			LastUpdate: time.Now(),
		}
		GNodeID++
		mutex.Unlock()

		clientList[r.RemoteAddr] = client
	}

	tmpl.Execute(w, struct {
		NodeID   string
		NumNodes int
	}{clientList[r.RemoteAddr].NodeID, len(clientList)})
}

func main() {
	http.HandleFunc("/api/work.js", GenerateWork)
	http.Handle("/static/", http.FileServer(http.Dir("template/")))
	http.HandleFunc("/api/ws", ServeWS)
	http.HandleFunc("/", DefaultHandler)

	http.ListenAndServe(":8080", nil)
}
