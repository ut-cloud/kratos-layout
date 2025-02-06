package boot

import (
	"flag"
	"fmt"
	"github.com/go-kratos/kratos-layout/internal/conf"
	knacos "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/ut-cloud/atlas-toolkit/utils"
)

var (
	flagconf string
	Name     string = "atlas-core"
)

func init() {
	// 初始化雪花算法
	if err := utils.Init("2024-01-01", 1); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	flag.StringVar(&flagconf, "conf", "configs", "config path, eg: -conf config.yaml")
	flag.Parse()
}

type Conf struct{}

func (c *Conf) Run() (*conf.Bootstrap, *conf.Registry) {
	rg := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	if err := rg.Load(); err != nil {
		panic(err)
	}
	var rc conf.Registry
	if err := rg.Scan(&rc); err != nil {
		panic(err)
	}
	var nacosConf utils.NacosAkConf
	if err := copier.Copy(&nacosConf, &rc.Nacos); err != nil {
		panic(err)
	}
	client, _ := utils.InitAkNacos(&nacosConf)
	var configSources []config.Source
	for _, configItem := range rc.Nacos.ConfigItems {
		configSources = append(configSources, knacos.NewConfigSource(
			client,
			knacos.WithGroup(configItem.GetGroupName()),
			knacos.WithDataID(configItem.GetDataId()),
		))
	}
	configParams := config.New(
		config.WithSource(configSources...),
	)
	defer configParams.Close()
	if err := configParams.Load(); err != nil {
		panic(err)
	}
	var bc conf.Bootstrap
	if err := configParams.Scan(&bc); err != nil {
		panic(err)
	}
	serverId, _ := uuid.NewUUID()
	bc.Server.Id = serverId.String()
	return &bc, &rc

}
