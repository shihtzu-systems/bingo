package cmd

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/shihtzu-systems/bingo/pkg/bingo"
	"github.com/shihtzu-systems/bingo/pkg/bingox"
	"github.com/shihtzu-systems/bingo/pkg/loggerx"
	"github.com/shihtzu-systems/bingo/pkg/tracerx"
	"github.com/shihtzu-systems/redix"

	jprom "github.com/uber/jaeger-lib/metrics/prometheus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	"github.com/uber/jaeger-client-go/transport"
	"github.com/uber/jaeger-lib/metrics"
	"strings"
	"time"
)

var serveCommand = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		// initialize
		theme := viper.GetString("bingo.v1.theme")

		var themedBoxes bingo.Boxes
		switch theme {
		case "cheesy christmas movies":
			fallthrough
		default:
			themedBoxes = christmasBoxes()
		}

		// tracing
		cfg := config.Configuration{
			ServiceName: viper.GetString("jaeger.v1.serviceName"),
			Sampler: &config.SamplerConfig{
				Type:  viper.GetString("jaeger.v1.samplerType"),
				Param: viper.GetFloat64("jaeger.v1.samplerParam"),
			},
			Reporter: &config.ReporterConfig{
				LogSpans:            false,
				BufferFlushInterval: 1 * time.Second,
				LocalAgentHostPort:  viper.GetString("jaeger.v1.agentHostPort"),
				CollectorEndpoint:   viper.GetString("jaeger.v1.collectorEndpoint"),
			},
		}
		metricsFactory := jprom.New().Namespace(metrics.NSOptions{
			Name: "bingo",
			Tags: nil,
		})
		opentracing.SetGlobalTracer(tracerx.Init("bingo", cfg, metricsFactory, loggerx.NewFactory(logger)))

		// serve
		bingox.Serve(bingox.ServeArgs{
			Serial: fmt.Sprintf("%s+on.%s.at.%s", version, datestamp, timestamp),
			Debug:  viper.GetBool("system.v2.debug"),

			SessionSecret: []byte(viper.GetString("server.v1.sessionSecret")),
			SessionKey:    viper.GetString("server.v1.sessionKey"),

			Logger: logger,

			Redis: redix.Redis{
				Address:  viper.GetString("redis.v1.address"),
				Port:     viper.GetInt("redis.v1.port"),
				Database: viper.GetInt("redis.v1.database"),
			},

			Boxes: themedBoxes,
		})
	},
}

func init() {

	rootCmd.AddCommand(serveCommand)
}

func christmasBoxes() (out bingo.Boxes) {
	contents := []string{
		"Main Character Returns to Small Town",
		"Storm",
		"Winter Carnival/Festival",
		"Concert",
		"Wise Old Women/Man/Couple",
		"Single Parent",
		"Sob Story",
		"Christmas Theme Name for Character",
		"Going out of Business",
		"Christmas Play",
		"Town with Christmas-themed Name",
		"Hunky Santa",
		"Fake Engagement/Marriage",
		"Travel Setbacks",
		"Dead Parent/Spouse",
		"Main Character Dislikes Holidays",
		"Odd Couple Share a Bed",
		"Odd Couple Teamup",
		"Celebrity Cameo",
		"Real Santa",
		"Busy Career Woman",
		"Movie Title Pun",
		"Decorating Montage",
		"Disapproving Parent",
		"Magical Item",
		"Highschool Sweethearts with Bad Breakup",
		"Sick/Dying Relative",
		"Parent/Child heart to heart",
		"Sidekick is gay",
		"Sidekick is non-white",
		"Childhood memory",
		"Interrupted kiss",
		"Lighting of the town tree",
		"No wifi",
	}
	for _, content := range contents {
		out = append(out, bingo.Box{
			Content: content,
			Marked:  false,
		})
	}
	return out
}

func Init(serviceName string, metricsFactory metrics.Factory, logger loggerx.Factory, backendHostPort string) opentracing.Tracer {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}
	// TODO(ys) a quick hack to ensure random generators get different seeds, which are based on current time.
	time.Sleep(100 * time.Millisecond)
	jaegerLogger := jaegerLoggerAdapter{logger.Bg()}
	var sender jaeger.Transport
	if strings.HasPrefix(backendHostPort, "http://") {
		sender = transport.NewHTTPTransport(
			backendHostPort,
			transport.HTTPBatchSize(1),
		)
	} else {
		if s, err := jaeger.NewUDPTransport(backendHostPort, 0); err != nil {
			logger.Bg().Fatal("cannot initialize UDP sender")
		} else {
			sender = s
		}
	}
	tracer, _, err := cfg.New(
		serviceName,
		config.Reporter(jaeger.NewRemoteReporter(
			sender,
			jaeger.ReporterOptions.BufferFlushInterval(1*time.Second),
			jaeger.ReporterOptions.Logger(jaegerLogger),
		)),
		config.Logger(jaegerLogger),
		config.Metrics(metricsFactory),
		config.Observer(rpcmetrics.NewObserver(metricsFactory, rpcmetrics.DefaultNameNormalizer)),
	)
	if err != nil {
		logger.Bg().Fatal("cannot initialize Jaeger Tracer")
	}
	return tracer
}

type jaegerLoggerAdapter struct {
	logger loggerx.Logger
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}
