package context

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ddliu/go-httpclient"
	"sync"
	"time"

	"github.com/sqeven/wechat/util"
)

const (
	//AccessTokenURL 获取access_token的接口
	AccessTokenURL = "https://api.weixin.qq.com/cgi-bin/token"
)

//ResAccessToken struct
type ResAccessToken struct {
	util.CommonError

	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

//GetAccessTokenFunc 获取 access token 的函数签名
type GetAccessTokenFunc func(ctx *Context) (accessToken string, err error)

//SetAccessTokenLock 设置读写锁（一个appID一个读写锁）
func (ctx *Context) SetAccessTokenLock(l *sync.RWMutex) {
	ctx.accessTokenLock = l
}

//SetGetAccessTokenFunc 设置自定义获取accessToken的方式, 需要自己实现缓存
func (ctx *Context) SetGetAccessTokenFunc(f GetAccessTokenFunc) {
	ctx.accessTokenFunc = f
}

//GetAccessToken 获取access_token
func (ctx *Context) GetAccessToken() (accessToken string, err error) {
	/*ctx.accessTokenLock.Lock()
	defer ctx.accessTokenLock.Unlock()

	if ctx.accessTokenFunc != nil {
		return ctx.accessTokenFunc(ctx)
	}
	accessTokenCacheKey := fmt.Sprintf("access_token_%s", ctx.AppID)
	if ctx.Cache != nil {
		val := ctx.Cache.Get(accessTokenCacheKey)
		if val != nil {
			accessToken = val.(string)
			return
		}
	}

	//从微信服务器获取
	var resAccessToken ResAccessToken
	resAccessToken, err = ctx.GetAccessTokenFromServer()
	if err != nil {
		return
	}

	accessToken = resAccessToken.AccessToken
	return*/
	// 中心TOKEN
	if ctx.TokenUrl == "" {
		return "", errors.New("无效地址")
	}
	tokenRes, err := httpclient.Get(ctx.TokenUrl, map[string]string{
		"appid": ctx.AppID,
	})

	if err != nil {
		return "", err
	} else {
		bodyString, _ := tokenRes.ToString()
		type gatapirsp struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
			Data string `json:"data"`
		}
		var rsp gatapirsp
		bodybyte := bytes.TrimPrefix([]byte(bodyString), []byte("\xef\xbb\xbf"))
		if err := json.Unmarshal(bodybyte, &rsp); err != nil {
			return "", err
		}
		if rsp.Code != 0 {
			return "", errors.New(rsp.Msg)
		} else {
			return rsp.Data, nil
		}
	}
}

//GetAccessTokenFromServer 强制从微信服务器获取token
func (ctx *Context) GetAccessTokenFromServer() (resAccessToken ResAccessToken, err error) {
	url := fmt.Sprintf("%s?grant_type=client_credential&appid=%s&secret=%s", AccessTokenURL, ctx.AppID, ctx.AppSecret)
	var body []byte
	body, err = util.HTTPGet(url)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resAccessToken)
	if err != nil {
		return
	}
	if resAccessToken.ErrMsg != "" {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v", resAccessToken.ErrCode, resAccessToken.ErrMsg)
		return
	}

	accessTokenCacheKey := fmt.Sprintf("access_token_%s", ctx.AppID)
	expires := resAccessToken.ExpiresIn - 1500
	if ctx.Cache != nil {
		err = ctx.Cache.Set(accessTokenCacheKey, resAccessToken.AccessToken, time.Duration(expires)*time.Second)
		return
	}
	return
}
