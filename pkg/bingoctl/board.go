package bingoctl

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/shihtzu-systems/bingo/pkg/bingo"
	"github.com/shihtzu-systems/bingo/pkg/bingosvc"
	"github.com/shihtzu-systems/bingo/pkg/bingoview"
	"github.com/shihtzu-systems/bingo/pkg/loggerx"
	"github.com/shihtzu-systems/redix"
	"go.uber.org/zap"
	"strconv"

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
	Logger loggerx.Logger
	Redis  redix.Redis

	SessionStore sessions.Store
	SessionKey   string

	Boxes bingo.Boxes
}

func (c BoardController) HandleRoot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	var board bingo.Board
	if !bingosvc.BoardExists(id, c.Redis) {
		board, err := bingosvc.NewBoard(id, c.Boxes)
		if err != nil {
			c.Logger.Fatal("unable to create board", zap.Error(err))
		}
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
	vars := mux.Vars(r)

	id := vars["id"]
	if !bingosvc.BoardExists(id, c.Redis) {
		c.Logger.Fatal("no board to mark")
	}
	board := bingosvc.GetBoard(id, c.Redis)

	letter := vars["letter"]
	i, err := strconv.Atoi(vars["index"])
	if err != nil {
		c.Logger.Fatal("unable to convert index string to int", zap.Error(err))
	}

	c.Logger.Debug("toggle",
		zap.String("column", letter),
		zap.Int("row", i))
	board.Mark(letter, i)
	bingosvc.SaveBoard(board, c.Redis)

	traceResponseHeaders(r.Context(), w)
	w.Header().Set("Location", BoardPath(board.Id, "check"))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (c BoardController) HandleCheck(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	if !bingosvc.BoardExists(id, c.Redis) {
		c.Logger.Fatal("no board to mark")
	}
	board := bingosvc.GetBoard(id, c.Redis)

	c.Logger.Debug("check for bingo board",
		zap.String("id", board.Id))

	bingoed := bingosvc.CheckForBingo(&board)
	c.Logger.Debug("check for bingoed",
		zap.Bool("bingoed", bingoed))
	bingosvc.SaveBoard(board, c.Redis)

	traceResponseHeaders(r.Context(), w)
	w.Header().Set("Location", BoardPath(board.Id))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (c BoardController) HandleRecycle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	board, err := bingosvc.NewBoard(id, c.Boxes)
	if err != nil {
		c.Logger.Fatal("unable to create new board",
			zap.String("id", id),
			zap.Int("boxes_count", len(c.Boxes)))
	}

	bingosvc.SaveBoard(board, c.Redis)

	traceResponseHeaders(r.Context(), w)
	w.Header().Set("Location", BoardPath(board.Id))
	w.WriteHeader(http.StatusTemporaryRedirect)

}

func (c BoardController) Id(w http.ResponseWriter, r *http.Request) string {
	store, err := c.SessionStore.Get(r, c.SessionKey)
	if err != nil {
		c.Logger.Fatal("error occurred while getting session store", zap.Error(err))
	}

	name, exists := store.Values["name"]
	if !exists {
		name = generateName()
		store.Values["name"] = name
	}
	if err := store.Save(r, w); err != nil {
		c.Logger.Fatal("error occurred while saving session store", zap.Error(err))
	}
	c.Logger.Debug("checking session", zap.String("name", store.Values["name"].(string)))
	return store.Values["name"].(string)
}
