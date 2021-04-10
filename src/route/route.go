package route

import (
	"fmt"
	"net/url"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/osins/osin-simple/simple"
	"github.com/osins/osin-simple/simple/config"
	simple_request "github.com/osins/osin-simple/simple/request"
	"github.com/osins/osin-simple/simple/response"
	"github.com/osins/osin-storage/storage/pg"
	"sso.humanrisk.cn/auth"
)

func New() Route {
	accessGen, err := auth.NewJwt()
	if err != nil {
		return nil
	}

	// ex.NewTestStorage implements the "osin.Storage" interface
	conf := config.NewServerConfig()
	conf.AllowClientSecretInParams = true
	conf.AccessExpiration = 1000000
	conf.AllowedAuthorizeTypes = simple_request.AllowedAuthorizeResponseType{
		simple_request.AUTHORIZE_RESPONSE_CODE,
		simple_request.AUTHORIZE_RESPONSE_LOGIN,
		simple_request.AUTHORIZE_RESPONSE_REGISTER,
	}
	conf.AllowAccessGrantType = simple_request.AllowedAccessGrantType{
		simple_request.ACCESS_GRANT_AUTHORIZATION_CODE,
		simple_request.ACCESS_GRANT_REFRESH_TOKEN,
		simple_request.ACCESS_GRANT_PASSWORD,
	}
	conf.AccessToken = accessGen
	conf.Storage.Client = pg.NewClientStorage()
	conf.Storage.User = pg.NewUserStorage()
	conf.Storage.Authorize = pg.NewAuthorizeStorage()
	conf.Storage.Access = pg.NewAccessStorage()

	return &route{
		Server: simple.NewSimpleServer(conf),
	}
}

type Route interface {
	Authorize(ctx *fiber.Ctx) error
	Token(ctx *fiber.Ctx) error
	Info(ctx *fiber.Ctx) error
}

type route struct {
	Server *simple.SimpleServer
}

func (r *route) ExceptionCatch(err interface{}) {
	if e := err; e != nil {
		fmt.Printf("\nexception catch: %s\n", e)
		debug.PrintStack()
	}
}

func (r *route) Authorize(ctx *fiber.Ctx) error {
	defer func() {
		r.ExceptionCatch(recover())
	}()

	fmt.Printf("authorize handle start:\n")
	fmt.Printf("method: %s\n", ctx.Route().Method)

	var req *simple_request.AuthorizeRequest
	if ctx.Route().Method == fiber.MethodPost {
		fmt.Printf("client_id: %s\n", ctx.FormValue("client_id"))
		req = &simple_request.AuthorizeRequest{
			ClientId:     ctx.FormValue("client_id"),
			ClientSecret: ctx.FormValue("client_secret"),
			ResponseType: simple_request.AuthorizeResponseType(ctx.FormValue("response_type")),
			RedirectUri:  ctx.FormValue("redirect_uri"),
			State:        ctx.FormValue("state"),
			Username:     ctx.FormValue("username"),
			Password:     ctx.FormValue("password"),
		}

		if req.ResponseType == simple_request.AUTHORIZE_RESPONSE_REGISTER {
			req.EMail = ctx.FormValue("email")
			req.Mobile = ctx.FormValue("mobile")
		}
	} else {
		fmt.Printf("client_id: %s\n", ctx.Query("client_id"))

		req = &simple_request.AuthorizeRequest{
			ClientId:     ctx.Query("client_id"),
			ClientSecret: ctx.Query("client_secret"),
			ResponseType: simple_request.AuthorizeResponseType(ctx.Query("response_type")),
			RedirectUri:  ctx.Query("redirect_uri"),
			State:        ctx.Query("state"),
		}
	}

	fmt.Printf("\nquerys: %s\n", req)

	authorize := simple.NewAuthorize(r.Server)
	var (
		res *response.AuthorizeResponse
		err error
	)

	switch req.ResponseType {
	case simple_request.AUTHORIZE_RESPONSE_CODE:
		res, err = authorize.Authorization(req)
		if err != nil {
			if res.NeedLogin {
				return ctx.Render("login", fiber.Map{
					"Title":     "Humanrisk Login",
					"authorize": req,
				})
			}
			fmt.Printf("authorize handle error:%s\n", err.Error())
			return err
		}
	case simple_request.AUTHORIZE_RESPONSE_LOGIN:
		res, err = authorize.Login(req)

		if err != nil {
			fmt.Printf("authorize handle error:%s\n", err.Error())
			return err
		}

		if res == nil {
			return nil
		}
	case simple_request.AUTHORIZE_RESPONSE_REGISTER:
		res, err = authorize.Register(req)
		if err != nil {
			fmt.Printf("authorize handle error:%s\n", err.Error())
			return err
		}
	}

	params := url.Values{
		"code":  {res.Code},
		"state": {res.State},
	}

	fmt.Printf("authorize handle complete.\n")

	res.RedirectUri = fmt.Sprintf("%s?%s", res.RedirectUri, params.Encode())

	if ctx.Route().Method == fiber.MethodPost &&
		(req.ResponseType == simple_request.AUTHORIZE_RESPONSE_LOGIN ||
			req.ResponseType == simple_request.AUTHORIZE_RESPONSE_REGISTER) {
		return ctx.JSON(res)
	}

	return ctx.Redirect(res.RedirectUri)
}

func (r *route) Token(ctx *fiber.Ctx) error {
	defer func() {
		r.ExceptionCatch(recover())
	}()

	req := &simple_request.AccessRequest{
		ClientId:           ctx.FormValue("client_id"),
		ClientSecret:       ctx.FormValue("client_secret"),
		GrantType:          simple_request.AccessGrantType(ctx.FormValue("grant_type")),
		Code:               ctx.FormValue("code"),
		Scope:              ctx.FormValue("scope"),
		State:              ctx.FormValue("state"),
		CodeVerifier:       ctx.FormValue("code_verifier"),
		CodeVerifierMethod: ctx.FormValue("code_verifier_method"),
		Expiration:         r.Server.Config.AccessExpiration,
		Authorized:         false,
	}

	if ctx.FormValue("grant_type") == "password" {
		req.Username = ctx.FormValue("username")
		req.Password = ctx.FormValue("password")
	}

	res, err := simple.NewAccess(r.Server).Access(req)
	if err != nil {
		return err
	}

	return ctx.JSON(res)
}

func (r *route) Info(ctx *fiber.Ctx) error {
	defer func() {
		r.ExceptionCatch(recover())
	}()

	code, err := simple.NewToken().AuthorizationToCode(ctx.Get("Authorization"))
	if err != nil {
		return err
	}

	req := &simple_request.InfoRequest{
		Code:  code,
		State: ctx.FormValue("state"),
	}

	if ad, err := simple.NewInfo(r.Server).Info(req); err != nil {
		return err
	} else {
		if err := ctx.JSON(ad); err != nil {
			return err
		}
	}

	return nil
}
