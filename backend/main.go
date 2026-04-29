package main

import (
	"github.com/starfederation/datastar-go/datastar"
	"time"
	"net/http"
	"fmt"
	"html/template"
	"bytes"
)

const chessboardTmplString = `
<div id="chess-board" class="grid grid-cols-8 w-64 h-64 border">
	{{range $rank := .Ranks}}
		{{range $file := $.Files}}
			<div class="w-8 h-8 border">
			text
			</div>
		{{end}}
	{{end}}
</div>
`
var chessboardTmpl = template.Must(template.New("board").Parse(chessboardTmplString))


func chessboardHandler(w http.ResponseWriter, r *http.Request) {

	data := struct {
		Ranks []int
		Files []int
	}{
		Ranks: []int{0,1,2,3,4,5,6,7,},
		Files: []int{0,1,2,3,4,5,6,7,},
	}
	var templateResult bytes.Buffer
	chessboardTmpl.Execute(&templateResult, data)

	sse := datastar.NewSSE(w, r)
	sse.PatchElements(templateResult.String())
}

func endpointHandler(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(`<div id="hal">I'm sorry, Dave. I'm afraid I can't do that.</div>`)
	time.Sleep(1 * time.Second)
	sse.PatchElements(`<div id="hal">Waiting for an order...</div>`)
}

func main() {
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)
	http.HandleFunc("/endpoint", endpointHandler)
	http.HandleFunc("/api/chess-board", chessboardHandler)
	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", nil)
}
