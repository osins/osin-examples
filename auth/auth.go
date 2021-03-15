package auth

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/osins/osin-simple/simple"
	"github.com/osins/osin-simple/simple/config"
)

func NewAuth(conf *config.SimpleConfig, server *simple.SimpleServer) Auth {
	return &auth{}
}

type Auth interface {
}

type auth struct {
	conf *config.SimpleConfig
}

func (s *auth) ClientExists(ctx *fiber.Ctx) (bool, error) {
	if f, err := s.conf.Storage.Client.Get(ctx.Query("client_id")); err != nil {
		return false, err
	} else if f == nil {
		return false, fmt.Errorf("client not exists.")
	}

	return true, nil
}
