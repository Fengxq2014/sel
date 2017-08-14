package apis

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/chanxuehong/rand"
	"github.com/chanxuehong/session"
	"github.com/chanxuehong/sid"

	"github.com/Fengxq2014/sel/models"
	"github.com/chanxuehong/wechat.v2/mp/core"
	"github.com/chanxuehong/wechat.v2/mp/menu"
	"github.com/chanxuehong/wechat.v2/mp/message/callback/request"
	"github.com/chanxuehong/wechat.v2/mp/message/callback/response"
	mpoauth2 "github.com/chanxuehong/wechat.v2/mp/oauth2"
	"github.com/chanxuehong/wechat.v2/oauth2"
	"github.com/gin-gonic/gin"
)

const (
	wxAppId     = "wx91579c743ffe0d71"
	wxAppSecret = "355d58755b0903fde1756cf842fef998"

	wxOriId         = "gh_f7cd4d03fd3b"
	wxToken         = "feng"
	wxEncodedAESKey = ""

	oauth2RedirectURI = "http://17255r43z1.51mypc.cn/oauth1" // 填上自己的参数
	oauth2Scope       = "snsapi_base"
)

var (
	// 下面两个变量不一定非要作为全局变量, 根据自己的场景来选择.
	msgHandler core.Handler
	msgServer  *core.Server

	sessionStorage                 = session.New(20*60, 60*60)
	oauth2Endpoint oauth2.Endpoint = mpoauth2.NewEndpoint(wxAppId, wxAppSecret)
)

func init() {
	mux := core.NewServeMux()
	mux.DefaultMsgHandleFunc(defaultMsgHandler)
	mux.DefaultEventHandleFunc(defaultEventHandler)
	mux.MsgHandleFunc(request.MsgTypeText, textMsgHandler)
	mux.EventHandleFunc(menu.EventTypeClick, menuClickEventHandler)

	msgHandler = mux
	msgServer = core.NewServer(wxOriId, wxAppId, wxToken, wxEncodedAESKey, msgHandler, nil)
}

func textMsgHandler(ctx *core.Context) {
	log.Printf("收到文本消息:\n%s\n", ctx.MsgPlaintext)

	msg := request.GetText(ctx.MixedMsg)
	resp := response.NewText(msg.FromUserName, msg.ToUserName, msg.CreateTime, msg.Content)
	ctx.RawResponse(resp) // 明文回复
	//ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}

func defaultMsgHandler(ctx *core.Context) {
	log.Printf("收到消息:\n%s\n", ctx.MsgPlaintext)
	ctx.NoneResponse()
}

func menuClickEventHandler(ctx *core.Context) {
	log.Printf("收到菜单 click 事件:\n%s\n", ctx.MsgPlaintext)

	event := menu.GetClickEvent(ctx.MixedMsg)
	resp := response.NewText(event.FromUserName, event.ToUserName, event.CreateTime, "收到 click 类型的事件")
	ctx.RawResponse(resp) // 明文回复
	// ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}

func defaultEventHandler(ctx *core.Context) {
	log.Printf("收到事件:\n%s\n", ctx.MsgPlaintext)
	ctx.NoneResponse()
}

func WeixinHandler(c *gin.Context) {
	msgServer.ServeHTTP(c.Writer, c.Request, nil)
}

// 建立必要的 session, 然后跳转到授权页面
func Page1Handler(c *gin.Context) {
	sid := sid.New()
	state := string(rand.NewHex())

	if err := sessionStorage.Add(sid, state); err != nil {
		io.WriteString(c.Writer, err.Error())
		log.Println(err)
		return
	}

	cookie := http.Cookie{
		Name:     "sid",
		Value:    sid,
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &cookie)

	AuthCodeURL := mpoauth2.AuthCodeURL(wxAppId, oauth2RedirectURI, oauth2Scope, "")
	log.Println("AuthCodeURL:", AuthCodeURL)

	http.Redirect(c.Writer, c.Request, AuthCodeURL, http.StatusFound)
}

// 授权后回调页面
func Page2Handler(c *gin.Context) {
	log.Println(c.Request.RequestURI)

	cookie, err := c.Cookie("sid")
	if err != nil {
		io.WriteString(c.Writer, err.Error())
		log.Println(err)
		return
	}

	session, err := sessionStorage.Get(cookie)
	if err != nil {
		io.WriteString(c.Writer, err.Error())
		log.Println(err)
		return
	}

	savedState := session.(string) // 一般是要序列化的, 这里保存在内存所以可以这么做

	code := c.Query("code")
	if code == "" {
		log.Println("用户禁止授权")
		return
	}

	queryState := c.Query("state")
	if queryState == "" {
		log.Println("state 参数为空")
		return
	}
	if savedState != queryState {
		str := fmt.Sprintf("state 不匹配, session 中的为 %q, url 传递过来的是 %q", savedState, queryState)
		io.WriteString(c.Writer, str)
		log.Println(str)
		return
	}

	oauth2Client := oauth2.Client{
		Endpoint: oauth2Endpoint,
	}
	token, err := oauth2Client.ExchangeToken(code)
	if err != nil {
		io.WriteString(c.Writer, err.Error())
		log.Println(err)
		return
	}
	log.Printf("token: %+v\r\n", token)

	userinfo, err := mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
	if err != nil {
		io.WriteString(c.Writer, err.Error())
		log.Println(err)
		return
	}
	AuthCodeURL := ""
	p := models.User{Openid: userinfo.OpenId}
	usercookie := http.Cookie{
		Name:     "openid",
		Value:    userinfo.OpenId,
		HttpOnly: true,
		MaxAge:   int(time.Hour / time.Second),
	}
	http.SetCookie(c.Writer, &usercookie)
	_, err = p.GetUserByOpenid()
	if err != nil {
		AuthCodeURL = "/login"
	} else {
		AuthCodeURL = "/list"
	}
	http.Redirect(c.Writer, c.Request, AuthCodeURL, http.StatusFound)
	// json.NewEncoder(c.Writer).Encode(userinfo)
	// log.Printf("userinfo: %+v\r\n", userinfo)
	return
}
