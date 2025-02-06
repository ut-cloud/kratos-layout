package boot

import (
	"github.com/go-kratos/kratos-layout/internal/conf"
)

type Boot interface {
	Run()
	Setting()
}

func NewBootConf() *Conf {
	return &Conf{}
}

func NewBootLog(conf *conf.Bootstrap) *Log {
	return &Log{
		conf,
	}
}

func NewBootTrace(conf *conf.Bootstrap) *Trace {
	return &Trace{
		conf,
	}
}
