package bingox

import (
	"context"
	haikunator "github.com/atrox/haikunatorgo/v2"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/shihtzu-systems/bingo/pkg/bingo"
	. "github.com/shihtzu-systems/bingo/pkg/bingoctl"
	"github.com/shihtzu-systems/redix"
	log "github.com/sirupsen/logrus"
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

	Redis redix.Redis

	Boxes  bingo.Boxes
	Serial string
}

func Serve(args ServeArgs) {
	r := mux.NewRouter()

	if args.Trace {
		log.SetLevel(log.TraceLevel)
	} else if args.Debug {
		log.SetLevel(log.DebugLevel)
	}

	sessionStore := sessions.NewCookieStore(args.SessionSecret)

	// root controller
	root := RootController{
		Redis:        args.Redis,
		SessionKey:   args.SessionKey,
		SessionStore: sessionStore,
		Boxes:        args.Boxes,
	}
	r.HandleFunc(RootPath(), root.HandleRoot)
	r.HandleFunc(RootPath("{letter:[bingo]}", "{index:[0-5]}"), root.HandleMark)
	r.HandleFunc(RootPath("recycle"), root.HandleRecycle)

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
