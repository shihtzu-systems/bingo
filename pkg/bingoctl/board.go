package bingoctl

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/shihtzu-systems/bingo/pkg/bingo"
	"github.com/shihtzu-systems/bingo/pkg/bingosvc"
	"github.com/shihtzu-systems/bingo/pkg/bingoview"
	"github.com/shihtzu-systems/redix"
	"strconv"

	log "github.com/sirupsen/logrus"
	"net/http"
	"path"
)

const (
	boardBasePath = "/board"
)

func BoardPath(pieces ...string) string {
	return path.Join(append([]string{boardBasePath}, pieces...)...)
}

type BoardController struct {
	Redis redix.Redis

	SessionStore sessions.Store
	SessionKey   string

	Boxes bingo.Boxes
}

func (c BoardController) HandleRoot(w http.ResponseWriter, r *http.Request) {
	log.Debug("handling ", BoardPath("{id}"))
	vars := mux.Vars(r)

	id := vars["id"]
	var board bingo.Board
	if !bingosvc.BoardExists(id, c.Redis) {
		board = bingosvc.NewBoard(id, c.Boxes)
		bingosvc.SaveBoard(board, c.Redis)
	} else {
		board = bingosvc.GetBoard(id, c.Redis)
	}

	v := bingoview.Bingo{
		Board: board,
	}
	v.View(w)
}

func (c BoardController) HandleMark(w http.ResponseWriter, r *http.Request) {
	log.Debug("handling ", BoardPath("{id}", "mark", "{letter}", "{index}"))
	vars := mux.Vars(r)

	id := vars["id"]
	if !bingosvc.BoardExists(id, c.Redis) {
		log.Fatal("no board to mark")
	}
	board := bingosvc.GetBoard(id, c.Redis)

	letter := vars["letter"]
	i, err := strconv.Atoi(vars["index"])
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("Toggle %s %d", letter, i)
	board.Mark(letter, i)
	bingosvc.SaveBoard(board, c.Redis)

	w.Header().Set("Location", BoardPath(board.Id, "check"))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (c BoardController) HandleCheck(w http.ResponseWriter, r *http.Request) {
	log.Debug("handling ", BoardPath("{id}", "check"))
	vars := mux.Vars(r)

	id := vars["id"]
	if !bingosvc.BoardExists(id, c.Redis) {
		log.Fatal("no board to check")
	}
	board := bingosvc.GetBoard(id, c.Redis)

	log.Debugf("Check %s for bingo", board.Id)
	bingoed := bingosvc.CheckForBingo(&board)
	log.Debug("bingo? ", bingoed)
	bingosvc.SaveBoard(board, c.Redis)

	w.Header().Set("Location", BoardPath(board.Id))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (c BoardController) HandleRecycle(w http.ResponseWriter, r *http.Request) {
	log.Debug("handling ", BoardPath("{id}", "recycle"))
	vars := mux.Vars(r)

	id := vars["id"]
	board := bingosvc.NewBoard(id, c.Boxes)
	bingosvc.SaveBoard(board, c.Redis)

	w.Header().Set("Location", BoardPath(board.Id))
	w.WriteHeader(http.StatusTemporaryRedirect)

}

func (c BoardController) Id(w http.ResponseWriter, r *http.Request) string {
	store, err := c.SessionStore.Get(r, c.SessionKey)
	if err != nil {
		log.Fatal(err)
	}

	name, exists := store.Values["name"]
	if !exists {
		name = generateName()
		store.Values["name"] = name
	}
	if err := store.Save(r, w); err != nil {
		log.Fatal(err)
	}
	log.Debug("id: ", store.Values["name"].(string))
	return store.Values["name"].(string)
}
