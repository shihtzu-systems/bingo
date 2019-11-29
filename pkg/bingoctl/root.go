package bingoctl

import (
	"github.com/gorilla/sessions"
	"github.com/shihtzu-systems/bingo/pkg/bingo"
	"github.com/shihtzu-systems/bingo/pkg/bingosvc"
	"github.com/shihtzu-systems/bingo/pkg/bingoview"

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
	SessionStore sessions.Store
	Boxes        bingo.Boxes
}

func (c RootController) HandleRoot(w http.ResponseWriter, r *http.Request) {
	log.Debug("handling ", RootPath())
	board := bingosvc.NewBingo(c.Boxes)
	v := bingoview.Bingo{
		Board: board,
	}
	v.View(w)
}
