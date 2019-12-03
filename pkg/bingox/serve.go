package bingox

import (
	"context"
	haikunator "github.com/atrox/haikunatorgo/v2"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/opentracing-contrib/go-gorilla/gorilla"
	"github.com/opentracing/opentracing-go"
	"github.com/shihtzu-systems/bingo/pkg/bingo"
	. "github.com/shihtzu-systems/bingo/pkg/bingoctl"
	"github.com/shihtzu-systems/redix"
	log "github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type ServeArgs struct {
	Trace bool
	Debug bool

	SessionSecret []byte
	SessionKey    string

	Redis       redix.Redis
	TraceConfig config.Configuration

	Boxes  bingo.Boxes
	Serial string
}

func Serve(args ServeArgs) {
	r := mux.NewRouter()

	if args.Debug {
		log.SetLevel(log.DebugLevel)
	}

	sessionStore := sessions.NewCookieStore(args.SessionSecret)

	// root controller
	root := RootController{
		Redis:        args.Redis,
		SessionKey:   args.SessionKey,
		SessionStore: sessionStore,
	}
	r.HandleFunc(RootPath(), root.HandleRoot)

	// board controller
	board := BoardController{
		Redis:        args.Redis,
		SessionKey:   args.SessionKey,
		SessionStore: sessionStore,
		Boxes:        args.Boxes,
	}
	r.HandleFunc(BoardPath("{id:[a-z0-9-]+}"), board.HandleRoot)
	r.HandleFunc(BoardPath("{id:[a-z0-9-]+}", "mark", "{letter:[bingo]}", "{index:[0-4]}"), board.HandleMark)
	r.HandleFunc(BoardPath("{id:[a-z0-9-]+}", "check"), board.HandleCheck)
	r.HandleFunc(BoardPath("{id:[a-z0-9-]+}", "recycle"), board.HandleRecycle)

	// static
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static/"))))

	// server startup
	namer := haikunator.New()
	namer.TokenLength = 0
	namer.Delimiter = " "
	name := namer.Haikunate()
	log.Printf("starting v%s as %s", args.Serial, name)
	log.Printf("listening on localhost:8080")
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// tracing
	tracer, closer, _ := args.TraceConfig.NewTracer(
		config.Logger(jaeger.StdLogger),
	)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		route.Handler(
			gorilla.Middleware(opentracing.GlobalTracer(), route.GetHandler()))
		return nil
	})

	// listen
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	// server teardown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Info("shutting down ", name)
	os.Exit(0)
}

func getTracer() (opentracing.Tracer, io.Closer, error) {
	//jaeger agent port
	jaegerHostPort := ":6831"

	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  jaegerHostPort,
		},
	}
	return cfg.New(
		"ExampleTracingMiddleware", //service name
	)
}
