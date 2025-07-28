package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/joinmouse/EasySwapBase/logger/xzap"
	"go.uber.org/zap"

	"github.com/joinmouse/EasySwapBackend/src/config"
	"github.com/joinmouse/EasySwapBackend/src/service/svc"
)

type Platform struct {
	config    *config.Config
	router    *gin.Engine
	serverCtx *svc.ServerCtx
}

func NewPlatform(config *config.Config, router *gin.Engine, serverCtx *svc.ServerCtx) (*Platform, error) {
	return &Platform{
		config:    config,
		router:    router,
		serverCtx: serverCtx,
	}, nil
}

func (p *Platform) Start() {
	xzap.WithContext(context.Background()).Info("EasySwap-End run", zap.String("port", p.config.Api.Port))
	if err := p.router.Run(p.config.Api.Port); err != nil {
		panic(err)
	}
}
