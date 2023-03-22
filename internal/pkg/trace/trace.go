package trace

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(GetUrl, NewTracer)

func ReplyHeader() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			traceID := tracing.TraceID()(ctx).(string)
			if tr, ok := transport.FromServerContext(ctx); ok {
				tr.ReplyHeader().Set("traceID", traceID)
			}
			return handler(ctx, req)
		}
	}
}
