package server

import (
	"context"
	v1 "github.com/go-kratos/kratos-layout/api/helloworld/v1"
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/handlers"
	casbinM "github.com/ut-cloud/atlas-toolkit/casbin"
	middle "github.com/ut-cloud/atlas-toolkit/middleware"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, rPool *redis.Pool, greeter *service.GreeterService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		NewMiddleware(rPool, logger),
		http.Filter(handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
		)),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}

// NewMiddleware 创建中间件
func NewMiddleware(rPool *redis.Pool, logger log.Logger) http.ServerOption {
	//菜单权限
	perms := make([]*casbinM.RoleMenuPerm, 0)
	// 菜单
	roleKeys := make([]string, 0)
	// 需要授权的接口
	operaPolicies := make([]string, 0)
	m, a := casbinM.InitCasbin(rPool, perms, roleKeys, operaPolicies)
	return http.Middleware(
		validate.Validator(),
		recovery.Recovery(),
		selector.Server(
			middle.Auth(),
			casbinM.Server(
				casbinM.WithCasbinModel(m),
				casbinM.WithCasbinPolicy(a),
				casbinM.WithSecurityUserCreator(middle.NewSecurityUser),
				casbinM.WithWhiteList(NewWithList()),
			),
		).Match(NewWhiteListMatcher()).Build(),
	)
}

// NewWhiteListMatcher
//
//	TODO
//	@Description: 设置白名单，不需要 token 验证的接口
//	@return selector.MatchFunc
func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})

	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// NewWithList
//
//	TODO
//	@Description: 设置白名单， 不需要casbin鉴权的接口
//	@return []string
func NewWithList() []string {
	return []string{}
}
