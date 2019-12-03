package bingoctl

import (
	haikunator "github.com/atrox/haikunatorgo/v2"
	"github.com/gorilla/sessions"
	"github.com/shihtzu-systems/bingo/pkg/bingosvc"
	"github.com/shihtzu-systems/redix"
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
}

func (c RootController) HandleRoot(w http.ResponseWriter, r *http.Request) {
	sessionId := c.Id(w, r)
	boardId := bingosvc.GetBoardId(sessionId, c.Redis)
	if boardId == "" {
		boardId = generateName()
		bingosvc.SaveBoardId(sessionId, boardId, c.Redis)
	}
	w.Header().Set("Location", BoardPath(boardId))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (c RootController) Id(w http.ResponseWriter, r *http.Request) string {
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

func generateName() string {
	namer := haikunator.New()
	namer.TokenLength = 6
	return namer.Haikunate()
}
