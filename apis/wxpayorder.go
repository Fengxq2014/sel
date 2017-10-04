package apis

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Fengxq2014/sel/conf"
	"github.com/Fengxq2014/sel/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UnifiedOrderRequest struct {
	XMLName struct{} `xml:"xml" json:"-"`

	// 必选参数
	AppID          string  `xml:"appid"`            // 微信支付分配的公众账号ID（企业号corpid即为此appId）
	Mch_id         string  `xml:"mch_id"`           // 微信支付分配的商户号
	Body           string  `xml:"body"`             // 商品或支付单简要描述
	OutTradeNo     string  `xml:"out_trade_no"`     // 商户系统内部的订单号,32个字符内、可包含字母, 其他说明见商户订单号
	TotalFee       float64 `xml:"total_fee"`        // 订单总金额，单位为分，详见支付金额
	SpbillCreateIP string  `xml:"spbill_create_ip"` // APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP。
	NotifyURL      string  `xml:"notify_url"`       // 接收微信支付异步通知回调地址，通知url必须为直接可访问的url，不能携带参数。
	TradeType      string  `xml:"trade_type"`       // 取值如下：JSAPI，NATIVE，APP，详细说明见参数规定
	Sign           string  `xml:"sign"`             // 通过签名算法计算得出的签名值，详见签名生成算法

	// 可选参数
	DeviceInfo string `xml:"device_info"` // 终端设备号(门店号或收银设备ID)，注意：PC网页或公众号内支付请传"WEB"
	NonceStr   string `xml:"nonce_str"`   // 随机字符串，不长于32位。NOTE: 如果为空则系统会自动生成一个随机字符串。
	SignType   string `xml:"sign_type"`   // 签名类型，默认为MD5，支持HMAC-SHA256和MD5。
	Detail     string `xml:"detail"`      // 商品名称明细列表
	Attach     string `xml:"attach"`      // 附加数据，在查询API和支付通知中原样返回，该字段主要用于商户携带订单的自定义数据
	FeeType    string `xml:"fee_type"`    // 符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	TimeStart  string `xml:"time_start"`  // 订单生成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010。其他详见时间规则
	TimeExpire string `xml:"time_expire"` // 订单失效时间，格式为yyyyMMddHHmmss，如2009年12月27日9点10分10秒表示为20091227091010。其他详见时间规则
	GoodsTag   string `xml:"goods_tag"`   // 商品标记，代金券或立减优惠功能的参数，说明详见代金券或立减优惠
	ProductId  string `xml:"product_id"`  // trade_type=NATIVE，此参数必传。此id为二维码中包含的商品ID，商户自行定义。
	LimitPay   string `xml:"limit_pay"`   // no_credit--指定不能使用信用卡支付
	OpenId     string `xml:"openid"`      // rade_type=JSAPI，此参数必传，用户在商户appid下的唯一标识。
	SubOpenId  string `xml:"sub_openid"`  // trade_type=JSAPI，此参数必传，用户在子商户appid下的唯一标识。openid和sub_openid可以选传其中之一，如果选择传sub_openid,则必须传sub_appid。
	SceneInfo  string `xml:"scene_info"`  // 该字段用于上报支付的场景信息,针对H5支付有以下三种场景,请根据对应场景上报,H5支付不建议在APP端使用，针对场景1，2请接入APP支付，不然可能会出现兼容性问题
}

type UnifiedOrderResponse struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
	Appid       string `xml:"appid"`
	Mch_id      string `xml:"mch_id"`
	Nonce_str   string `xml:"nonce_str"`
	Sign        string `xml:"sign"`
	Result_code string `xml:"result_code"`
	Prepay_id   string `xml:"prepay_id"`
	Trade_type  string `xml:"trade_type"`
	TimeStamp   string `xml:"timeStamp"`
}

// WxPayOrder 生成支付订单
func WxPayOrder(c *gin.Context) {
	type param struct {
		Name     string  `form:"name" binding:"required"`      //课程名称
		CourseId string  `form:"course_id" binding:"required"` //课程ID
		Price    float64 `form:"price" binding:"required"`     //价格
		OpenId   string  `form:"openid" binding:"required"`    //用户openid
		Uid      int     `form:"user_id" binding:"required"`   //用户ID
		Cid      int     `form:"child_id"`                     //儿童ID
	}

	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	order := UnifiedOrderRequest{}
	order.AppID = conf.Config.WXAppID
	order.Mch_id = conf.Config.Mch_id
	order.Body = queryStr.Name
	order.OutTradeNo = "sel" + time.Now().Format("20060102150405") + queryStr.CourseId
	order.TotalFee = queryStr.Price
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	order.SpbillCreateIP = ip
	// order.SpbillCreateIP = "192.168.10.12"
	if err != nil {
		c.Error(errors.New("解析客户端地址失败"))
		return
	}
	order.NotifyURL = conf.Config.CallBack
	order.TradeType = "JSAPI"
	order.OpenId = queryStr.OpenId

	order.DeviceInfo = ""
	order.NonceStr = time.Now().Format("20060102150405") + randSeq(10)
	// log.Printf("NonceStr:" + order.NonceStr)
	order.SignType = "MD5"
	order.Detail = ""
	println("Cid:" + strconv.Itoa(queryStr.Cid))
	if queryStr.Cid != 0 {
		order.Attach = queryStr.CourseId + "|" + strconv.Itoa(queryStr.Uid) + "|" + strconv.Itoa(queryStr.Cid)
	} else {
		order.Attach = queryStr.CourseId + "|" + strconv.Itoa(queryStr.Uid)
	}
	order.FeeType = "CNY"

	// local, err := time.LoadLocation("PRC") //服务器设置的时区
	// if err != nil {
	// 	fmt.Println(err)
	// }
	order.TimeStart = time.Now().Format("20060102150405")
	order.TimeExpire = time.Now().Add(time.Minute * 5).Format("20060102150405")
	order.GoodsTag = ""
	order.ProductId = ""
	order.LimitPay = ""

	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["appid"] = conf.Config.WXAppID
	m["mch_id"] = conf.Config.Mch_id
	m["body"] = order.Body
	m["out_trade_no"] = order.OutTradeNo
	m["total_fee"] = order.TotalFee
	m["spbill_create_ip"] = order.SpbillCreateIP
	m["notify_url"] = order.NotifyURL
	m["trade_type"] = order.TradeType
	m["openid"] = order.OpenId

	m["device_info"] = order.DeviceInfo
	m["nonce_str"] = order.NonceStr
	m["sign_type"] = order.SignType
	m["detail"] = order.Detail
	m["attach"] = order.Attach
	m["fee_type"] = order.FeeType
	m["time_start"] = order.TimeStart
	m["time_expire"] = order.TimeExpire
	m["goods_tag"] = order.GoodsTag
	m["product_id"] = order.ProductId
	m["limit_pay"] = order.LimitPay
	order.Sign = wxpayCalcSign(m, conf.Config.Key)

	res := models.Result{}
	bytes_req, err := xml.Marshal(order)
	if err != nil {
		res.Res = 1
		res.Msg = "以xml形式编码发送错误, 原因:" + err.Error()
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	}

	str_req := string(bytes_req)
	//wxpay的unifiedorder接口需要http body中xmldoc的根节点是<xml></xml>这种，所以这里需要replace一下
	str_req = strings.Replace(str_req, "XUnifyOrderReq", "xml", -1)
	bytes_req = []byte(str_req)

	//发送unified order请求.
	req, err := http.NewRequest("POST", "https://api.mch.weixin.qq.com/pay/unifiedorder", bytes.NewReader(bytes_req))
	if err != nil {
		res.Res = 1
		res.Msg = "New Http Request发生错误，原因:" + err.Error()
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	}
	req.Header.Set("Accept", "application/xml")
	//这里的http header的设置是必须设置的.
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")

	cc := http.Client{}
	resp, _err := cc.Do(req)
	if _err != nil {
		res.Res = 1
		res.Msg = "请求微信支付统一下单接口发送错误, 原因:" + _err.Error()
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	}
	//------------------到这里统一下单接口就已经执行完成了-------------------

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		res.Res = 1
		res.Msg = "解析返回body错误" + err.Error()
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	}
	xmlResp := UnifiedOrderResponse{}
	_err = xml.Unmarshal(respBytes, &xmlResp)
	//处理return code.
	if xmlResp.Return_code == "FAIL" {
		res.Res = 1
		res.Msg = "微信支付统一下单不成功，原因:" + xmlResp.Return_msg
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	}

	xmlResp.TimeStamp = strconv.FormatInt(time.Now().Unix(), 10)
	xmlResp.Nonce_str = order.NonceStr

	var mm map[string]interface{}
	mm = make(map[string]interface{}, 0)
	mm["appId"] = conf.Config.WXAppID
	mm["timeStamp"] = xmlResp.TimeStamp
	mm["nonceStr"] = order.NonceStr
	mm["package"] = "prepay_id=" + xmlResp.Prepay_id
	mm["signType"] = "MD5"

	xmlResp.Sign = wxpayCalcSign(mm, conf.Config.Key)

	c.JSON(http.StatusOK, models.Result{Data: xmlResp})
}

// wxpayCalcSign 微信支付 下单签名
func wxpayCalcSign(mReq map[string]interface{}, key string) (sign string) {
	fmt.Println("微信支付签名计算, API KEY:", key)
	//STEP 1, 对key进行升序排序.
	sorted_keys := make([]string, 0)
	for k, _ := range mReq {
		sorted_keys = append(sorted_keys, k)
	}

	sort.Strings(sorted_keys)

	//STEP2, 对key=value的键值对用&连接起来，略过空值
	var signStrings string
	for _, k := range sorted_keys {
		fmt.Printf("k=%v, v=%v\n", k, mReq[k])
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}

	//STEP3, 在键值对的最后加上key=API_KEY
	if key != "" {
		signStrings = signStrings + "key=" + key
	}

	//STEP4, 进行MD5签名并且将所有字符转为大写.
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	return upperSign
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//计算支付签名 跟下单签名不同的地方在于 最后一个字符串连接没有&
func wxpaySign(mReq map[string]interface{}, key string) string {
	//STEP 1, 对key进行升序排序.
	sorted_keys := make([]string, 0)
	for k, _ := range mReq {
		sorted_keys = append(sorted_keys, k)
	}

	sort.Strings(sorted_keys)

	//STEP2, 对key=value的键值对用&连接起来，略过空值
	var signStrings string
	for i, k := range sorted_keys {
		//fmt.Printf("k=%v, v=%v\n", k, mReq[k])
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" {
			if i != (len(sorted_keys) - 1) {
				signStrings = signStrings + k + "=" + value + "&"
			} else {
				signStrings = signStrings + k + "=" + value //最后一个不加此符号
			}
		}
	}
	//fmt.Println("=====键值对==============", signStrings)

	//STEP3, 在键值对的最后加上key=API_KEY
	if key != "" {
		signStrings = signStrings + "&key=" + key
	}
	fmt.Println("=====wxpaySign 键值对加key==============", signStrings)

	//STEP4, 进行MD5签名并且将所有字符转为大写.
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))

	fmt.Println("=====进行MD5签名并且将所有字符转为大写 ==============", upperSign)
	return upperSign
}
