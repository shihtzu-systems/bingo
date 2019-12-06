// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tracerx

import (
	"fmt"
	"github.com/shihtzu-systems/bingo/pkg/loggerx"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	"github.com/uber/jaeger-client-go/transport"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

// Init creates a new instance of Jaeger tracer.
func Init(serviceName string, cfg config.Configuration, metricsFactory metrics.Factory, logger loggerx.Factory) opentracing.Tracer {
	// TODO(ys) a quick hack to ensure random generators get different seeds, which are based on current time.
	time.Sleep(100 * time.Millisecond)
	jaegerLogger := jaegerLoggerAdapter{logger.Bg()}
	var sender jaeger.Transport
	if strings.HasPrefix(cfg.Reporter.LocalAgentHostPort, "http://") {
		sender = transport.NewHTTPTransport(
			cfg.Reporter.LocalAgentHostPort,
			transport.HTTPBatchSize(1),
		)
	} else {
		if s, err := jaeger.NewUDPTransport(cfg.Reporter.LocalAgentHostPort, 0); err != nil {
			logger.Bg().Fatal("cannot initialize UDP sender", zap.Error(err))
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
		logger.Bg().Fatal("cannot initialize Jaeger Tracer", zap.Error(err))
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
