package context

import (
	"github.com/valyala/fasthttp"
	"net/http"
	"sync"

	"github.com/sqeven/wechat/cache"
)

// Context struct
type Context struct {
	AppID          string
	AppSecret      string
	Token          string
	EncodingAESKey string
	PayMchID       string
	PayNotifyURL   string
	PayKey         string

	Cache cache.Cache

	HttType string // http/fasthttp

	// 源生http处理
	Writer  http.ResponseWriter
	Request *http.Request

	// fasthttp处理
	FastHttpCtx    *fasthttp.RequestCtx
	FastHttpWriter fasthttp.Response

	//accessTokenLock 读写锁 同一个AppID一个
	accessTokenLock *sync.RWMutex

	//jsAPITicket 读写锁 同一个AppID一个
	jsAPITicketLock *sync.RWMutex

	//accessTokenFunc 自定义获取 access token 的方法
	accessTokenFunc GetAccessTokenFunc
}

// Query returns the keyed url query value if it exists
func (ctx *Context) Query(key string) string {
	value, _ := ctx.GetQuery(key)
	return value
}

// GetQuery is like Query(), it returns the keyed url query value
func (ctx *Context) GetQuery(key string) (string, bool) {
	if "fasthttp" == ctx.HttType {
		val := string(ctx.FastHttpCtx.QueryArgs().Peek(key))
		if len(val) > 0 {
			return val, true
		}
	} else {
		req := ctx.Request
		if values, ok := req.URL.Query()[key]; ok && len(values) > 0 {
			return values[0], true
		}
	}
	return "", false
}

// SetJsAPITicketLock 设置jsAPITicket的lock
func (ctx *Context) SetJsAPITicketLock(lock *sync.RWMutex) {
	ctx.jsAPITicketLock = lock
}

// GetJsAPITicketLock 获取jsAPITicket 的lock
func (ctx *Context) GetJsAPITicketLock() *sync.RWMutex {
	return ctx.jsAPITicketLock
}
