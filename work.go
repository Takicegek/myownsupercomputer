package main

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

var GNodeID = 0
var mutex = sync.Mutex{}

type clientInfo struct {
	NodeID     string
	Address    net.IP
	LastUpdate time.Time
}

var clientList map[string]clientInfo

func init() {
	clientList = make(map[string]clientInfo)
}

func GenerateWork(w http.ResponseWriter, r *http.Request) {
	log.Println("Generating work for ", r.RemoteAddr)

	tmpl, err := template.ParseFiles("templates/work.js")
	if err != nil {
		Write404(w, r)
		return
	}

	if client, ok := clientList[r.RemoteAddr]; ok {
		// we've seen this guy before
		tmpl.Execute(w, struct {
			ServerHost string
			ServerPort string
			NodeID     string
		}{"localhost", "8080", client.NodeID})

		client.LastUpdate = time.Now()
		clientList[r.RemoteAddr] = client
	} else {
		// how'd this guy get here without going through our index page??
		Write404(w, r)

		log.Println("Attempted rogue connection from ", r.RemoteAddr)
	}
}
