package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	v1 "github.com/qin8948050/kratos-layout/api/helloworld/v1"
	"github.com/qin8948050/kratos-layout/internal/conf"
	"github.com/qin8948050/kratos-layout/internal/pkg/trace"
	"github.com/qin8948050/kratos-layout/internal/service"
	http2 "net/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *http.Server {
	err := trace.InitTracer("http://10.8.132.173:14268/api/traces")
	if err != nil {
		panic(err)
	}
	var opts = []http.ServerOption{
		http.Middleware(
			tracing.Server(),
			metadata.Server(),
			recovery.Recovery(),
			logging.Server(logger),
		),
		http.ResponseEncoder(responseEncoder),
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

func responseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if v == nil {
		return nil
	}
	if rd, ok := v.(http.Redirector); ok {
		url, code := rd.Redirect()
		http2.Redirect(w, r, url, code)
		return nil
	}
	codec, _ := http.CodecForRequest(r, "Accept")
	w.Header().Set("TraceID", "FEFE")
	data, err := codec.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}
