package bingoctl

import (
	"context"
	haikunator "github.com/atrox/haikunatorgo/v2"
	"github.com/gorilla/sessions"
	"github.com/opentracing/opentracing-go"
	"github.com/shihtzu-systems/bingo/pkg/bingosvc"
	"github.com/shihtzu-systems/bingo/pkg/loggerx"
	"github.com/shihtzu-systems/redix"
	"go.uber.org/zap"

	"net/http"
	"path"
)

const (
	rootBasePath = "/"
)

func RootPath(pieces ...string) string {
	return path.Join(append([]string{rootBasePath}, pieces...)...)
}

func NewRootController(logx loggerx.Factory,
	redis redix.Redis,
	sessionStore sessions.Store,
	sessionKey string) RootController {
	return RootController{
		logx:         logx,
		redis:        redis,
		sessionStore: sessionStore,
		sessionKey:   sessionKey,
	}
}

type RootController struct {
	logx  loggerx.Factory
	redis redix.Redis

	sessionStore sessions.Store
	sessionKey   string
}

func (c RootController) HandleRoot(w http.ResponseWriter, r *http.Request) {
	sessionId := c.Id(w, r)
	boardId, err := bingosvc.GetBoardId(sessionId, c.redis)
	if err != nil {
		c.logx.For(r.Context()).Fatal("unable to get board", zap.Error(err),
			zap.String("session_id", sessionId),
			zap.String("redis_address", c.redis.Address),
			zap.Int("redis_port", c.redis.Port),
			zap.Int("redis_database", c.redis.Database))
	}

	if boardId == "" {
		boardId = generateName()
		c.logx.For(r.Context()).Debug("generated board", zap.String("board_id", boardId))
		if err := bingosvc.SaveBoardId(sessionId, boardId, c.redis); err != nil {
			c.logx.For(r.Context()).Fatal("unable to save board", zap.Error(err),
				zap.String("session_id", sessionId),
				zap.String("board_id", boardId),
				zap.String("redis_address", c.redis.Address),
				zap.Int("redis_port", c.redis.Port),
				zap.Int("redis_database", c.redis.Database))
		}
	}
	w.Header().Set("Location", BoardPath(boardId))
	traceResponseHeaders(r.Context(), w)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (c RootController) Id(w http.ResponseWriter, r *http.Request) string {
	store, err := c.sessionStore.Get(r, c.sessionKey)
	if err != nil {
		c.logx.For(r.Context()).Fatal("unable to get session", zap.Error(err),
			zap.String("session_key", c.sessionKey))
	}

	name, exists := store.Values["name"]
	if !exists {
		name = generateName()
		store.Values["name"] = name
	}
	if err := store.Save(r, w); err != nil {
		c.logx.For(r.Context()).Fatal("unable to save session store", zap.Error(err),
			zap.String("session_key", c.sessionKey))
	}
	c.logx.For(r.Context()).Debug("checking session", zap.String("name", store.Values["name"].(string)))
	return store.Values["name"].(string)
}

func generateName() string {
	namer := haikunator.New()
	namer.TokenLength = 6
	return namer.Haikunate()
}

func traceResponseHeaders(ctx context.Context, w http.ResponseWriter) {
	_ = opentracing.GlobalTracer().Inject(
		opentracing.SpanFromContext(ctx).Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(w.Header()),
	)
}
