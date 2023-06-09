package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/qin8948050/kratos-layout/api/helloworld/v1"
	"github.com/qin8948050/kratos-layout/internal/biz"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer
	uc  *biz.GreeterUsecase
	log *log.Helper
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase, logger log.Logger) *GreeterService {
	return &GreeterService{uc: uc, log: log.NewHelper(logger)}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	/*	fmt.Println(tracing.TraceID()(ctx).(string))
		metadata.Metadata{}.Set("traceID", tracing.TraceID()(ctx).(string))
		a := metadata.Metadata{}.Get("traceID")
		fmt.Println(a)*/
	s.log.WithContext(ctx).Infow("logd test", "fefee")
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}
