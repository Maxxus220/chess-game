package main

import (
	"github.com/starfederation/datastar-go/datastar"
	"time"
	"net/http"
	"fmt"
)

func root_handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/index.html")
}

func endpoint_handler(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(`<div id="hal">I'm sorry, Dave. I'm afraid I can't do that.</div>`)
	time.Sleep(1 * time.Second)
	sse.PatchElements(`<div id="hal">Waiting for an order...</div>`)
}

func main() {
	http.HandleFunc("/", root_handler)
	http.HandleFunc("/endpoint", endpoint_handler)
	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", nil)
}
