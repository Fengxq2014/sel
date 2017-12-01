package apis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Fengxq2014/sel/models"

	"github.com/Fengxq2014/sel/tool"

	"github.com/chanxuehong/rand"
	"github.com/chanxuehong/session"
	"github.com/chanxuehong/sid"

	"github.com/Fengxq2014/sel/conf"
	"github.com/chanxuehong/wechat.v2/mp/core"
	"github.com/chanxuehong/wechat.v2/mp/jssdk"
	"github.com/chanxuehong/wechat.v2/mp/media"
	"github.com/chanxuehong/wechat.v2/mp/menu"
	"github.com/chanxuehong/wechat.v2/mp/message/callback/request"
	"github.com/chanxuehong/wechat.v2/mp/message/callback/response"
	"github.com/chanxuehong/wechat.v2/mp/message/template"
	mpoauth2 "github.com/chanxuehong/wechat.v2/mp/oauth2"
	"github.com/chanxuehong/wechat.v2/oauth2"
	"github.com/gin-gonic/gin"
)

var (
	wxAppId           = conf.Config.WXAppID
	wxAppSecret       = conf.Config.WXAppSecret
	wxOriId           = conf.Config.WXOriID
	wxToken           = conf.Config.WXToken
	wxEncodedAESKey   = ""
	oauth2RedirectURI = conf.Config.Oauth2RedirectURI
	oauth2Scope       = "snsapi_userinfo"
	// 下面两个变量不一定非要作为全局变量, 根据自己的场景来选择.
	msgHandler core.Handler
	msgServer  *core.Server

	sessionStorage                           = session.New(20*60, 60*60)
	oauth2Endpoint    oauth2.Endpoint        = mpoauth2.NewEndpoint(wxAppId, wxAppSecret)
	accessTokenServer core.AccessTokenServer = core.NewDefaultAccessTokenServer(wxAppId, wxAppSecret, nil)
	wechatClient      *core.Client           = core.NewClient(accessTokenServer, nil)
	ticketserver                             = jssdk.NewDefaultTicketServer(wechatClient)
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
		MaxAge:   int(time.Minute / time.Second),
	}
	http.SetCookie(c.Writer, &cookie)

	AuthCodeURL := mpoauth2.AuthCodeURL(wxAppId, oauth2RedirectURI+"?menuType="+c.Query("menuType"), oauth2Scope, state)

	http.Redirect(c.Writer, c.Request, AuthCodeURL, http.StatusFound)
}

// Page2Handler 授权后回调页面
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
		tool.Error(err)
		return
	}

	userinfo, err := mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
	if err != nil {
		io.WriteString(c.Writer, err.Error())
		tool.Error(err)
		return
	}
	usercookie1 := http.Cookie{
		Name:  "openid",
		Value: userinfo.OpenId,
	}
	usercookie2 := http.Cookie{
		Name:  "nickname",
		Value: url.QueryEscape(userinfo.Nickname),
	}
	usercookie3 := http.Cookie{
		Name:  "headimgurl",
		Value: userinfo.HeadImageURL,
	}
	accesstoken, err := ticketserver.Ticket()
	if err != nil {
		io.WriteString(c.Writer, err.Error())
		tool.Error(err)
		return
	}
	noncestr := string(rand.NewHex())
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sign := jssdk.WXConfigSign(accesstoken, noncestr, timestamp, conf.Config.Host+"/front/dist/?")
	ss := []string{conf.Config.WXAppID, timestamp, noncestr, sign}
	wxconfigCookie := http.Cookie{
		Name:  "wxconfig",
		Value: strings.Join(ss, "|"),
	}
	http.SetCookie(c.Writer, &usercookie1)
	http.SetCookie(c.Writer, &usercookie2)
	http.SetCookie(c.Writer, &usercookie3)
	http.SetCookie(c.Writer, &wxconfigCookie)
	AuthCodeURL := ""
	switch menuType := c.Query("menuType"); menuType {
	case "1":
		AuthCodeURL = "/front/dist/?#/appbase/assessment"
	case "2":
		AuthCodeURL = "/front/dist/?#/appbase/course"
	case "3":
		AuthCodeURL = "/front/dist/?#/appbase/mine"
	default:
		AuthCodeURL = "/front/appbase/mine"
	}
	http.Redirect(c.Writer, c.Request, AuthCodeURL, http.StatusFound)
	return
}

// DownloadMedia 通过mediaid下载媒体文件
func DownloadMedia(c *gin.Context) {
	mediaID := c.Query("mediaid")
	if mediaID == "" {
		c.Error(errors.New("参数为空"))
		return
	}
	fileName := getFileName(mediaID)
	if checkFileIsExist(fileName) {
		c.JSON(http.StatusOK, models.Result{Data: "/front/childimg/" + mediaID + ".jpg"})
		return
	}
	myfile, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		c.Error(err)
		return
	}
	_, err = media.DownloadToWriter(wechatClient, mediaID, myfile)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, models.Result{Data: "/front/childimg/" + mediaID + ".jpg"})
}

func getFileName(mediaID string) string {
	pwd, _ := os.Getwd()
	return filepath.Join(pwd, "front", "childimg", mediaID+".jpg")
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// TemplateMessage 发送模板消息
func TemplateMessage(openid, url, evaluationName, evaluationTime, nick_name, childName string) (err error) {
	type TemplateMessage struct {
		ToUser     string          `json:"touser"`        // 必须, 接受者OpenID
		TemplateId string          `json:"template_id"`   // 必须, 模版ID
		URL        string          `json:"url,omitempty"` // 可选, 用户点击后跳转的URL, 该URL必须处于开发者在公众平台网站中设置的域中
		Data       json.RawMessage `json:"data"`          // 必须, 模板数据, JSON 格式的 []byte, 满足特定的模板需求
	}
	json := `{
		"first": {
			"value":"您好，` + nick_name + `，您有一份完整测评报告已生成。",
			"color":"#89bd41"
		},
		"keyword1":{
			"value":"` + evaluationName + `",
			"color":"#89bd41"
		},
		"keyword2": {
			"value":"` + childName + `",
			"color":"#89bd41"
		},
		"keyword3": {
			"value":"` + evaluationTime + `",
			"color":"#89bd41"
		},
		"remark":{
			"value":"点击查看完整测评报告。您可以转发本消息，与家人分享报告。",
			"color":"#89bd41"
		}}`

	tool.Info("json:" + json)
	var jsonBlob = []byte(json)

	msg := TemplateMessage{ToUser: openid, TemplateId: conf.Config.Template_id, URL: url, Data: jsonBlob}
	msgid, err := template.Send(wechatClient, msg)
	id := strconv.FormatInt(msgid, 10)
	if err != nil && id != "" {
		return err
	}
	return err
}
