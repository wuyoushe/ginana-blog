package api

import (
	"ginana-blog/internal/controller"
	"ginana-blog/internal/model"
	"github.com/griffin702/ginana/library/log"
	"github.com/griffin702/service/jwt-iris"
	"github.com/kataras/iris/v12"
	"time"
)

type CApiLogin struct {
	controller.BaseController
}

// GetLoginCaptcha godoc
// @Description 获取验证码
// @Tags Login
// @Accept  json
// @Produce  json
// @Success 200 {string} string "data:image/png;base64,U3dhZ2dlciByb2Nrcw=="
// @Failure 500 {object} model.JSON
// @Router /login/captcha [get]
func (c *CApiLogin) GetLoginCaptcha() {
	captcha, err := c.Svc.GetCaptcha()
	if err != nil {
		c.Ctx.JSON(c.JsonPlus(nil, err))
		return
	}
	c.Session.Set(c.Hm.GetCacheKey(8), captcha.Code)
	c.Ctx.ContentType("image/png")
	c.Ctx.Write(captcha.Image)
}

// PostLoginCaptchaCheck godoc
// @Description 提前检查验证码
// @Tags Login
// @Accept  json
// @Produce  json
// @Param captcha body model.Captcha true "Check Captcha"
// @Success 200 {object} model.JSON{data=bool}
// @Failure 500 {object} model.JSON
// @Router /login/captcha/check [post]
func (c *CApiLogin) PostLoginCaptchaCheck() {
	captcha := new(model.Captcha)
	if err := c.Ctx.ReadJSON(&captcha); err != nil {
		c.Ctx.JSON(c.JsonPlus(nil, err))
		return
	}
	code := c.Session.Get(c.Hm.GetCacheKey(8))
	if code == nil {
		err := c.Hm.GetMessage(1006)
		c.Ctx.JSON(c.JsonPlus(nil, err))
		return
	}
	if captcha.Code == code {
		c.Ctx.JSON(c.JsonPlus(true, nil))
		return
	}
	c.Ctx.JSON(c.JsonPlus(false, nil))
	return
}

// PostLogin godoc
// @Description 登陆接口
// @Tags Login
// @Accept  json
// @Produce  json
// @Param user body model.UserLoginReq true "User Login"
// @Success 200 {object} model.JSON{data=bool}
// @Failure 500 {object} model.JSON
// @Router /login [post]
func (c *CApiLogin) PostLogin() {
	req := new(model.UserLoginReq)
	if err := c.Ctx.ReadJSON(req); err != nil {
		c.Ctx.JSON(c.JsonPlus(nil, err))
		return
	}
	if err := c.Valid(req); err != nil {
		c.Ctx.JSON(c.JsonPlus(nil, err))
		return
	}
	code := c.Session.Get(c.Hm.GetCacheKey(8))
	if code == nil {
		err := c.Hm.GetMessage(1006)
		c.Ctx.JSON(c.JsonPlus(nil, err))
		return
	}
	if code != req.Captcha {
		err := c.Hm.GetMessage(1007)
		c.Ctx.JSON(c.JsonPlus(nil, err))
		return
	}
	req.LoginIP = c.GetClientIP()
	user, err := c.Svc.PostLogin(req)
	if err != nil {
		c.Ctx.JSON(c.JsonPlus(nil, err))
		return
	}
	c.setToken(user)
	log.Infof("userId: %d, username: %s, 登录成功", user.ID, user.Username)
	c.Ctx.JSON(c.JsonPlus(true, c.Hm.GetMessage(0, "登陆成功")))
}

// PostRegister godoc
// @Description 注册接口
// @Tags Login
// @Accept  json
// @Produce  json
// @Param user body model.UserRegisterReq true "User Register"
// @Success 200 {object} model.JSON{data=bool}
// @Failure 500 {object} model.JSON
// @Router /register [post]
func (c *CApiLogin) PostRegister() {
	req := new(model.UserRegisterReq)
	if err := c.Ctx.ReadJSON(req); err != nil {
		c.Ctx.JSON(c.JsonPlus(nil, err))
		return
	}
	if err := c.Valid(req); err != nil {
		c.Ctx.JSON(c.JsonPlus(nil, err))
		return
	}
	code := c.Session.Get(c.Hm.GetCacheKey(8))
	if code == nil {
		err := c.Hm.GetMessage(1006)
		c.Ctx.JSON(c.JsonPlus(nil, err))
		return
	}
	if code != req.Captcha {
		err := c.Hm.GetMessage(1007)
		c.Ctx.JSON(c.JsonPlus(nil, err))
		return
	}
	req.LoginIP = c.GetClientIP()
	user, err := c.Svc.PostRegister(req)
	if err != nil {
		c.Ctx.JSON(c.JsonPlus(nil, err))
		return
	}
	c.setToken(user)
	log.Infof("userId: %d, username: %s, 登录成功", user.ID, user.Username)
	c.Ctx.JSON(c.JsonPlus(true, c.Hm.GetMessage(0, "注册并登陆成功")))
}

// PostLogout godoc
// @Description 安全登出
// @Tags Login
// @Accept  json
// @Produce  json
// @Success 200 {object} model.JSON{data=bool}
// @Failure 500 {object} model.JSON
// @Router /logout [get]
func (c *CApiLogin) GetLogout() {
	c.Session.Destroy()
	c.Ctx.RemoveCookie("token")
	c.ShowMsg("已安全登出，期待与您下一次遇见...")
}

func (c *CApiLogin) setToken(user *model.User) {
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Duration(c.Config.SessionAndCookieExpire)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["userId"] = user.ID
	claims["username"] = user.Username
	token := c.Tool.JwtGenerate(claims, c.Config.JwtSecret)
	c.Ctx.SetCookieKV("token", token,
		iris.CookieExpires(time.Duration(c.Config.SessionAndCookieExpire)),
	)
	c.Session.Set("token", token)
}
