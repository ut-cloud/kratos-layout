package boot

import (
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"os"
)

type Log struct {
	conf *conf.Bootstrap
}

func (b *Log) Run() log.Logger {
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", b.conf.Server.Id,
		"service.name", b.conf.Server.Name,
		"service.version", b.conf.Server.Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	return logger
}
