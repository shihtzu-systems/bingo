package bingoctl

import (
	haikunator "github.com/atrox/haikunatorgo/v2"
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
	rootBasePath = "/"
)

func RootPath(pieces ...string) string {
	return path.Join(append([]string{rootBasePath}, pieces...)...)
}

type RootController struct {
	Redis redix.Redis

	SessionStore sessions.Store
	SessionKey   string

	Boxes bingo.Boxes
}

func (c RootController) HandleRoot(w http.ResponseWriter, r *http.Request) {
	log.Debug("handling ", RootPath())

	id := c.Id(w, r)
	log.Debug(id)
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

func (c RootController) HandleMark(w http.ResponseWriter, r *http.Request) {
	log.Debug("handling ", RootPath("letter", "index"))
	vars := mux.Vars(r)

	id := c.Id(w, r)
	log.Debug(id)
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

	w.Header().Set("Location", RootPath())
	w.WriteHeader(http.StatusTemporaryRedirect)

}

func (c RootController) HandleRecycle(w http.ResponseWriter, r *http.Request) {
	log.Debug("handling ", RootPath("recycle"))

	id := c.Id(w, r)
	log.Debug(id)
	board := bingosvc.NewBoard(id, c.Boxes)
	bingosvc.SaveBoard(board, c.Redis)

	w.Header().Set("Location", RootPath())
	w.WriteHeader(http.StatusTemporaryRedirect)

}

func (c RootController) Id(w http.ResponseWriter, r *http.Request) string {
	store, err := c.SessionStore.Get(r, c.SessionKey)
	if err != nil {
		log.Fatal(err)
	}

	name, exists := store.Values["name"]
	if !exists {
		name = generateSessionName()
		store.Values["name"] = name
	}
	if err := store.Save(r, w); err != nil {
		log.Fatal(err)
	}
	log.Debug("id: ", store.Values["name"].(string))
	return store.Values["name"].(string)
}

func generateSessionName() string {
	namer := haikunator.New()
	namer.TokenLength = 6
	return namer.Haikunate()
}
