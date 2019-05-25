package main

import (
	"fmt"
	"net/http"
)

func testConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, your connection is fine. %s!", r.URL.Path[1:])
}