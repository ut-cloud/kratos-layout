package data

import (
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/model"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/gomodule/redigo/redis"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData, NewRegistrar, NewDiscovery,
	model.NewDB, model.NewRedisPool,
	NewGreeterRepo)

// Data .
type Data struct {
	Db    *gorm.DB
	RPool *redis.Pool
	// TODO wrapped database client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, db *gorm.DB, rPool *redis.Pool) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{Db: db, RPool: rPool}, cleanup, nil
}

func NewRegistrar(conf *conf.Registry) registry.Registrar {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(conf.GetNacos().GetAddress(), conf.GetNacos().GetPort()),
	}
	cc := constant.ClientConfig{
		NamespaceId: conf.GetNacos().GetNamespaceId(),
		AccessKey:   conf.GetNacos().GetAccessKey(),
		SecretKey:   conf.GetNacos().GetSecretKey(),
	}
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}
	r := nacos.New(client, nacos.WithGroup(conf.GetNacos().GetGroupName()))
	return r
}

func NewDiscovery(conf *conf.Registry) registry.Discovery {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(conf.GetNacos().GetAddress(), conf.GetNacos().GetPort()),
	}
	cc := constant.ClientConfig{
		NamespaceId: conf.GetNacos().GetNamespaceId(),
		AccessKey:   conf.GetNacos().GetAccessKey(),
		SecretKey:   conf.GetNacos().GetSecretKey(),
	}
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}
	r := nacos.New(client, nacos.WithGroup(conf.GetNacos().GetGroupName()))
	return r
}
