package front

import (
	"ginana-blog/internal/config"
	"ginana-blog/internal/model"
	"ginana-blog/internal/service"
	"github.com/griffin702/service/tools"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"strings"
)

type CFront struct {
	Ctx          iris.Context
	Session      *sessions.Session
	Svc          service.Service
	JsonPlus     model.JsonPlus
	Pager        *model.Pager
	GetClientIP  model.GetClientIP
	GetOption    model.GetOption
	Links        *model.Links
	DisableRight bool
	Hm           service.HelperMap
	Valid        model.Validator
	Tool         *tools.Tool
	Config       *config.Config
}

func (c *CFront) IsLogin() bool {
	userId := c.Session.Get("userId")
	if userId != nil && userId.(int64) > 0 {
		return true
	}
	return false
}

func (c *CFront) setHeadMetas(params ...string) {
	c.Ctx.ViewData("isLogin", c.IsLogin())
	c.Ctx.ViewData("disableRight", c.DisableRight)
	titleBuf := make([]string, 0, 3)
	if len(params) == 0 && c.GetOption("sitename") != "" {
		titleBuf = append(titleBuf, c.GetOption("sitename"))
	}
	if len(params) > 0 {
		titleBuf = append(titleBuf, params[0])
	}
	titleBuf = append(titleBuf, c.GetOption("subtitle"))
	c.Ctx.ViewData("title", strings.Join(titleBuf, " - "))
	if len(params) > 1 {
		c.Ctx.ViewData("keywords", params[1])
	} else {
		c.Ctx.ViewData("keywords", c.GetOption("keywords"))
	}
	if len(params) > 2 {
		c.Ctx.ViewData("description", params[2])
	} else {
		c.Ctx.ViewData("description", c.GetOption("description"))
	}
}
