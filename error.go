package main

import "net/http"

func Write404(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("<html><body><h1>Oops!</h1> <h3>Page Not Found</h3>"))
}
