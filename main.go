package main

import (
	"bytes"
	"fmt"
	"github.com/starfederation/datastar-go/datastar"
	"html/template"
	"net/http"
)

const chessboardTmplString = `
<div id="chess-board">
	{{range $rank_index, $rank := .Ranks}}
		{{range $file_index, $file := $.Files}}
		<div class="chess-square" style="background-color: {{getSquareColor $rank_index $file_index}};"></div>
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
			return "green"
		} else {
			return "grey"
		}
	},
}).Parse(chessboardTmplString))

func chessboardHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Ranks []int
		Files []int
	}{
		Ranks: []int{0, 1, 2, 3, 4, 5, 6, 7},
		Files: []int{0, 1, 2, 3, 4, 5, 6, 7},
	}
	var templateResult bytes.Buffer
	chessboardTmpl.Execute(&templateResult, data)

	sse := datastar.NewSSE(w, r)
	sse.PatchElements(templateResult.String())
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
	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}
