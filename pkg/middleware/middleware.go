package middleware

import (
	"github.com/magicLian/gostarter/pkg/log"
	"github.com/magicLian/gostarter/pkg/setting"
)

var (
	RequireSignedIn = Auth(&AuthOptions{RequireSignedIn: true})
	RequireAdmin    = Auth(&AuthOptions{RequireSignedIn: true, RequireAdmin: true})
)

type MiddleWare struct {
	log log.Logger
	cfg *setting.Cfg
}

func ProvideMiddleWare(cfg *setting.Cfg) *MiddleWare {
	return &MiddleWare{
		log: log.New("middleware"),
		cfg: cfg,
	}
}
