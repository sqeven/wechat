package wechat

import (
	"github.com/valyala/fasthttp"
	"net/http"
	"sync"

	"github.com/sqeven/wechat/cache"
	"github.com/sqeven/wechat/context"
	"github.com/sqeven/wechat/js"
	"github.com/sqeven/wechat/material"
	"github.com/sqeven/wechat/menu"
	"github.com/sqeven/wechat/miniprogram"
	"github.com/sqeven/wechat/oauth"
	"github.com/sqeven/wechat/pay"
	"github.com/sqeven/wechat/qr"
	"github.com/sqeven/wechat/server"
	"github.com/sqeven/wechat/template"
	"github.com/sqeven/wechat/user"
)

// Wechat struct
type Wechat struct {
	Context *context.Context
}

// Config for user
type Config struct {
	AppID          string
	AppSecret      string
	Token          string
	TokenUrl       string
	EncodingAESKey string
	PayMchID       string //支付 - 商户 ID
	PayNotifyURL   string //支付 - 接受微信支付结果通知的接口地址
	PayKey         string //支付 - 商户后台设置的支付 key
	Cache          cache.Cache
}

// NewWechat init
func NewWechat(cfg *Config) *Wechat {
	context := new(context.Context)
	copyConfigToContext(cfg, context)
	return &Wechat{context}
}

func copyConfigToContext(cfg *Config, context *context.Context) {
	context.AppID = cfg.AppID
	context.AppSecret = cfg.AppSecret
	context.Token = cfg.Token
	context.TokenUrl = cfg.TokenUrl
	context.EncodingAESKey = cfg.EncodingAESKey
	context.PayMchID = cfg.PayMchID
	context.PayKey = cfg.PayKey
	context.PayNotifyURL = cfg.PayNotifyURL
	context.Cache = cfg.Cache
	context.SetAccessTokenLock(new(sync.RWMutex))
	context.SetJsAPITicketLock(new(sync.RWMutex))
}

// GetServer 消息管理
func (wc *Wechat) GetServer(req *http.Request, writer http.ResponseWriter) *server.Server {
	wc.Context.Request = req
	wc.Context.Writer = writer
	return server.NewServer(wc.Context)
}

// GetServer 消息管理
func (wc *Wechat) GetServerWithFastHttp(req *fasthttp.RequestCtx, writer fasthttp.Response) *server.Server {
	wc.Context.FastHttpCtx = req
	wc.Context.FastHttpWriter = writer
	wc.Context.HttType = "fasthttp"
	return server.NewServer(wc.Context)
}

//GetAccessToken 获取access_token
func (wc *Wechat) GetAccessToken() (string, error) {
	return wc.Context.GetAccessToken()
}

// GetOauth oauth2网页授权
func (wc *Wechat) GetOauth() *oauth.Oauth {
	return oauth.NewOauth(wc.Context)
}

// GetMaterial 素材管理
func (wc *Wechat) GetMaterial() *material.Material {
	return material.NewMaterial(wc.Context)
}

// GetJs js-sdk配置
func (wc *Wechat) GetJs() *js.Js {
	return js.NewJs(wc.Context)
}

// GetMenu 菜单管理接口
func (wc *Wechat) GetMenu() *menu.Menu {
	return menu.NewMenu(wc.Context)
}

// GetUser 用户管理接口
func (wc *Wechat) GetUser() *user.User {
	return user.NewUser(wc.Context)
}

// GetTemplate 模板消息接口
func (wc *Wechat) GetTemplate() *template.Template {
	return template.NewTemplate(wc.Context)
}

// GetPay 返回支付消息的实例
func (wc *Wechat) GetPay() *pay.Pay {
	return pay.NewPay(wc.Context)
}

// GetQR 返回二维码的实例
func (wc *Wechat) GetQR() *qr.QR {
	return qr.NewQR(wc.Context)
}

// GetMiniProgram 获取小程序的实例
func (wc *Wechat) GetMiniProgram() *miniprogram.MiniProgram {
	return miniprogram.NewMiniProgram(wc.Context)
}
