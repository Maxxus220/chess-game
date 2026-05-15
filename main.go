package main

import (
	"bytes"
	"encoding/json"
	"github.com/starfederation/datastar-go/datastar"
	"html/template"
	"log"
	"net/http"
)

type DatastarSignals struct {
	SelectedSquare string `json:"selectedSquare"`
}

const chessboardTmplString = `
<div id="chess-board">
	{{range $rank_index, $rank := .Ranks}}
		{{range $file_index, $file := $.Files}}
		<div class="chess-square" data-on:click="$selectedSquare = {{$file}}{{$rank}}; @post('/api/select-square');"  style="background-color: {{getSquareColor $rank_index $file_index}};"></div>
		{{end}}
	{{end}}
</div>
`

var chessboardTmpl = template.Must(template.New("board").Funcs(template.FuncMap{
	"getSquareColor": func(rank int, file int) string {
		rank_is_even := rank%2 == 0
		file_is_even := file%2 == 0
		is_green := false

		if rank_is_even {
			is_green = file_is_even
		} else {
			is_green = !file_is_even
		}

		if is_green {
			return "#739552"
		} else {
			return "#ebecd0"
		}
	},
}).Parse(chessboardTmplString))

func chessboardHandler(w http.ResponseWriter, r *http.Request) {
	log.Default().Printf("Chessboard handler.\n")
	data := struct {
		Ranks []int
		Files []rune
	}{
		Ranks: []int{0, 1, 2, 3, 4, 5, 6, 7},
		Files: []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'},
	}
	var templateResult bytes.Buffer
	if err := chessboardTmpl.Execute(&templateResult, data); err != nil {
		log.Default().Printf("Error executing template: %v\n", err)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElements(templateResult.String())
}

func selectSquareHandler(w http.ResponseWriter, r *http.Request) {
	var signals DatastarSignals
	if err := json.NewDecoder(r.Body).Decode(&signals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Default().Printf("Here with val %s\n", signals.SelectedSquare)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, "./static/index.html")
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/api/chess-board", chessboardHandler)
	http.HandleFunc("/api/select-square", selectSquareHandler)
	log.Default().Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Default().Printf("Server failed: %s\n", err)
	}
}
