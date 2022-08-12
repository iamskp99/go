package main

import (
	"fmt"
	"net/http"
)

type hotdog int

func (m hotdog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	typeOfRequest := r.Method
	if typeOfRequest == "GET" {
		fmt.Fprintln(w, "Welcome Its a GET Request!")
	}

	if typeOfRequest == "POST" {
		fmt.Fprintln(w, "Welcome Its a POST Request!")
	}
}

func main() {
	var d hotdog
	http.ListenAndServe(":8080", d)
}
