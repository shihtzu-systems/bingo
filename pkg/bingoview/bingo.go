package bingoview

import (
	"github.com/shihtzu-systems/bingo/pkg/bingo"
	"html/template"
	"log"
	"net/http"
	"path"
)

var (
	// parts
	headerPartPath  = "parts/header.html"
	scriptsPartPath = "parts/scripts.html"

	bingoPath     = "bingo.html"
	templatesPath = "tpl"
)

type Bingo struct {
	Board bingo.Board
}

func (v Bingo) View(w http.ResponseWriter) {
	tpl, err := template.ParseFiles(
		path.Join(templatesPath, headerPartPath),
		path.Join(templatesPath, bingoPath),
		path.Join(templatesPath, scriptsPartPath),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := tpl.ExecuteTemplate(w, "bingo", v); err != nil {
		log.Fatal(err)
	}
}
