package boot

import (
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/ut-cloud/atlas-toolkit/utils"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Trace struct {
	conf *conf.Bootstrap
}

func (t *Trace) Run() *trace.TracerProvider {
	tp, err := utils.SetTracerProvider(t.conf.Trace.Endpoint, t.conf.Server.Environment, t.conf.Server.Name)
	if err != nil {
		panic(err)
	}
	return tp
}
