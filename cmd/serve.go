package cmd

import (
	"fmt"
	"github.com/shihtzu-systems/bingo/pkg/bingo"
	"github.com/shihtzu-systems/bingo/pkg/bingox"
	"github.com/shihtzu-systems/redix"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go/config"
	"time"
)

var serveCommand = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		theme := viper.GetString("bingo.v1.theme")
		log.Debug("current theme: ", theme)

		var themedBoxes bingo.Boxes
		switch theme {
		case "cheesy christmas movies":
			fallthrough
		default:
			themedBoxes = christmasBoxes()
		}

		log.Debug("boxes: ", len(themedBoxes))

		bingox.Serve(bingox.ServeArgs{
			Serial: fmt.Sprintf("%s+on.%s.at.%s", version, datestamp, timestamp),
			Debug:  viper.GetBool("system.v2.debug"),

			SessionSecret: []byte(viper.GetString("server.v1.sessionSecret")),
			SessionKey:    viper.GetString("server.v1.sessionKey"),

			TraceConfig: config.Configuration{
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
			},

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
