//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/qin8948050/kratos-layout/internal/biz"
	"github.com/qin8948050/kratos-layout/internal/conf"
	"github.com/qin8948050/kratos-layout/internal/data"
	"github.com/qin8948050/kratos-layout/internal/pkg/log"
	"github.com/qin8948050/kratos-layout/internal/pkg/trace"
	"github.com/qin8948050/kratos-layout/internal/server"
	"github.com/qin8948050/kratos-layout/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp, log.ProviderSet, trace.ProviderSet))
}
