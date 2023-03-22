package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	v1 "github.com/qin8948050/kratos-layout/api/helloworld/v1"
	"github.com/qin8948050/kratos-layout/internal/conf"
	"github.com/qin8948050/kratos-layout/internal/pkg/trace"
	"github.com/qin8948050/kratos-layout/internal/service"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger, provider *tracesdk.TracerProvider) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			trace.ReplyHeader(),
			logging.Server(logger),
		),
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
	//swagger ui
	srv.HandlePrefix("/q/", openapiv2.NewHandler())
	v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
