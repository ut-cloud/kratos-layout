package data

import (
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/model"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	model.NewDB, model.NewRedisPool, NewData,
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
