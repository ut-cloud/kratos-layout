package main

import (
	"github.com/go-kratos/kratos-layout/internal/boot"
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	_ "go.uber.org/automaxprocs"
)

func newApp(conf *conf.Server, logger log.Logger, gs *grpc.Server, r registry.Registrar) *kratos.App {
	return kratos.New(
		kratos.ID(conf.Id),
		kratos.Name(conf.Name),
		kratos.Version(conf.Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
		),
		kratos.Registrar(r),
	)
}

func main() {
	bc, rc := boot.NewBootConf().Run()
	logger := boot.NewBootLog(bc).Run()
	tp := boot.NewBootTrace(bc).Run()

	app, cleanup, err := wireApp(bc.Server, bc.Data, rc, tp, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
