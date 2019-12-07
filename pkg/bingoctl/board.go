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
	Logger loggerx.Factory
	Redis  redix.Redis

	SessionStore sessions.Store
	SessionKey   string

	Boxes bingo.Boxes
}

func (c BoardController) HandleRoot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	var board bingo.Board
	exists, err := bingosvc.BoardExists(id, c.Redis)
	if err != nil {
		c.Logger.For(r.Context()).Fatal("unable to check if board exists", zap.Error(err),
			zap.String("board_id", id),
			zap.String("redis_address", c.Redis.Address),
			zap.Int("redis_port", c.Redis.Port),
			zap.Int("redis_database", c.Redis.Database))
	}
	if !exists {
		board, err := bingosvc.NewBoard(id, c.Boxes)
		if err != nil {
			c.Logger.For(r.Context()).Fatal("unable to create new board", zap.Error(err),
				zap.String("board_id", id),
				zap.Int("boxes", len(c.Boxes)))
		}
		if err := bingosvc.SaveBoard(board, c.Redis); err != nil {
			c.Logger.For(r.Context()).Fatal("unable to save board", zap.Error(err),
				zap.String("board_id", id),
				zap.String("redis_address", c.Redis.Address),
				zap.Int("redis_port", c.Redis.Port),
				zap.Int("redis_database", c.Redis.Database))
		}
	} else {
		board, err = bingosvc.GetBoard(id, c.Redis)
		if err != nil {
			c.Logger.For(r.Context()).Fatal("unable to create new board", zap.Error(err),
				zap.String("board_id", id),
				zap.String("redis_address", c.Redis.Address),
				zap.Int("redis_port", c.Redis.Port),
				zap.Int("redis_database", c.Redis.Database))
		}
	}

	v := bingoview.Bingo{
		Board: board,
	}
	v.View(w)
}

func (c BoardController) HandleMark(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	exists, err := bingosvc.BoardExists(id, c.Redis)
	if err != nil {
		c.Logger.For(r.Context()).Fatal("unable to check if board exists", zap.Error(err),
			zap.String("board_id", id),
			zap.String("redis_address", c.Redis.Address),
			zap.Int("redis_port", c.Redis.Port),
			zap.Int("redis_database", c.Redis.Database))
	}
	if !exists {
		c.Logger.For(r.Context()).Fatal("unable to find board", zap.Error(err),
			zap.String("board_id", id),
			zap.String("redis_address", c.Redis.Address),
			zap.Int("redis_port", c.Redis.Port),
			zap.Int("redis_database", c.Redis.Database))
	}

	board, err := bingosvc.GetBoard(id, c.Redis)
	if err != nil {
		c.Logger.For(r.Context()).Fatal("unable to get board", zap.Error(err),
			zap.String("board_id", id),
			zap.String("redis_address", c.Redis.Address),
			zap.Int("redis_port", c.Redis.Port),
			zap.Int("redis_database", c.Redis.Database))
	}

	letter := vars["letter"]
	i, err := strconv.Atoi(vars["index"])
	if err != nil {
		c.Logger.For(r.Context()).Fatal("unable to convert index string to int", zap.Error(err),
			zap.String("index", vars["index"]),
			zap.String("column", letter),
			zap.String("row", vars["index"]))
	}

	c.Logger.For(r.Context()).Debug("toggle",
		zap.String("column", letter),
		zap.Int("row", i))
	board.Mark(letter, i)
	if err := bingosvc.SaveBoard(board, c.Redis); err != nil {
		c.Logger.For(r.Context()).Fatal("unable to save board", zap.Error(err),
			zap.String("board_id", id),
			zap.String("redis_address", c.Redis.Address),
			zap.Int("redis_port", c.Redis.Port),
			zap.Int("redis_database", c.Redis.Database))
	}

	traceResponseHeaders(r.Context(), w)
	w.Header().Set("Location", BoardPath(board.Id, "check"))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (c BoardController) HandleCheck(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	exists, err := bingosvc.BoardExists(id, c.Redis)
	if err != nil {
		c.Logger.For(r.Context()).Fatal("unable to check if board exists", zap.Error(err),
			zap.String("board_id", id),
			zap.String("redis_address", c.Redis.Address),
			zap.Int("redis_port", c.Redis.Port),
			zap.Int("redis_database", c.Redis.Database))
	}
	if !exists {
		c.Logger.For(r.Context()).Fatal("unable to find board", zap.Error(err),
			zap.String("board_id", id),
			zap.String("redis_address", c.Redis.Address),
			zap.Int("redis_port", c.Redis.Port),
			zap.Int("redis_database", c.Redis.Database))
	}

	board, err := bingosvc.GetBoard(id, c.Redis)
	if err != nil {
		c.Logger.For(r.Context()).Fatal("unable to get board", zap.Error(err),
			zap.String("board_id", id),
			zap.String("redis_address", c.Redis.Address),
			zap.Int("redis_port", c.Redis.Port),
			zap.Int("redis_database", c.Redis.Database))
	}

	c.Logger.For(r.Context()).Debug("check for bingo board",
		zap.String("id", board.Id))

	bingoed := bingosvc.CheckForBingo(&board)
	c.Logger.For(r.Context()).Debug("check for bingoed",
		zap.Bool("bingoed", bingoed))

	if err := bingosvc.SaveBoard(board, c.Redis); err != nil {
		c.Logger.For(r.Context()).Fatal("unable to save board", zap.Error(err),
			zap.String("board_id", id),
			zap.String("redis_address", c.Redis.Address),
			zap.Int("redis_port", c.Redis.Port),
			zap.Int("redis_database", c.Redis.Database))
	}

	traceResponseHeaders(r.Context(), w)
	w.Header().Set("Location", BoardPath(board.Id))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (c BoardController) HandleRecycle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	board, err := bingosvc.NewBoard(id, c.Boxes)
	if err != nil {
		c.Logger.For(r.Context()).Fatal("unable to create new board", zap.Error(err),
			zap.String("board_id", id),
			zap.Int("boxes", len(c.Boxes)))
	}

	if err := bingosvc.SaveBoard(board, c.Redis); err != nil {
		c.Logger.For(r.Context()).Fatal("unable to save board", zap.Error(err),
			zap.String("board_id", id),
			zap.String("redis_address", c.Redis.Address),
			zap.Int("redis_port", c.Redis.Port),
			zap.Int("redis_database", c.Redis.Database))
	}

	traceResponseHeaders(r.Context(), w)
	w.Header().Set("Location", BoardPath(board.Id))
	w.WriteHeader(http.StatusTemporaryRedirect)

}

func (c BoardController) Id(w http.ResponseWriter, r *http.Request) string {
	store, err := c.SessionStore.Get(r, c.SessionKey)
	if err != nil {
		c.Logger.For(r.Context()).Fatal("error occurred while getting session store", zap.Error(err))
	}

	name, exists := store.Values["name"]
	if !exists {
		name = generateName()
		store.Values["name"] = name
	}
	if err := store.Save(r, w); err != nil {
		c.Logger.For(r.Context()).Fatal("error occurred while saving session store", zap.Error(err))
	}
	c.Logger.For(r.Context()).Debug("checking session", zap.String("name", store.Values["name"].(string)))
	return store.Values["name"].(string)
}
