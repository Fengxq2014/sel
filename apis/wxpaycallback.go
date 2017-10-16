package apis

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Fengxq2014/sel/models"
	"github.com/gin-gonic/gin"
)

type PayOrderResponse struct {
	Return_code string `form:"return_code" xml:"return_code" binding:"required"` //返回状态码
	Return_msg  string `form:"return_msg" xml:"return_msg" binding:"required"`   //返回信息

	Appid                string `form:"return_code" xml:"appid"`                //公众账号ID
	Mch_id               string `form:"return_code" xml:"mch_id"`               //商户号
	Device_info          string `form:"return_code" xml:"device_info"`          //设备号
	Nonce_str            string `form:"return_code" xml:"nonce_str"`            //随机字符串
	Sign                 string `form:"return_code" xml:"sign"`                 //签名
	Sign_type            string `form:"return_code" xml:"sign_type"`            //签名类型
	Result_code          string `form:"return_code" xml:"result_code"`          //业务结果
	Err_code             string `form:"return_code" xml:"err_code"`             //错误代码
	Err_code_des         string `form:"return_code" xml:"err_code_des"`         //错误代码描述
	Openid               string `form:"return_code" xml:"openid"`               //用户标识
	Is_subscribe         string `form:"return_code" xml:"is_subscribe"`         //是否关注公众账号
	Trade_type           string `form:"return_code" xml:"trade_type"`           //交易类型
	Bank_type            string `form:"return_code" xml:"bank_type"`            //付款银行
	Total_fee            string `form:"return_code" xml:"total_fee"`            //订单金额
	Settlement_total_fee string `form:"return_code" xml:"settlement_total_fee"` //应结订单金额
	Fee_type             string `form:"return_code" xml:"fee_type"`             //货币种类
	Cash_fee             string `form:"return_code" xml:"cash_fee"`             //现金支付金额
	Cash_fee_type        string `form:"return_code" xml:"cash_fee_type"`        //现金支付金额
	Coupon_fee           string `form:"return_code" xml:"coupon_fee"`           //总代金券金额
	Coupon_count         string `form:"return_code" xml:"coupon_count"`         //代金券使用数量
	Coupon_type_nn       string `form:"return_code" xml:"coupon_type_$n"`       //代金券类型
	Coupon_id_nn         string `form:"return_code" xml:"coupon_id_$n"`         //代金券ID
	Coupon_fee_nn        string `form:"return_code" xml:"coupon_fee_$n"`        //单个代金券支付金额
	Transaction_id       string `form:"return_code" xml:"transaction_id"`       //微信支付订单号
	Out_trade_no         string `form:"return_code" xml:"out_trade_no"`         //商户订单号
	Attach               string `form:"return_code" xml:"attach"`               //商家数据包
	Time_end             string `form:"return_code" xml:"time_end"`             //支付完成时间
}

type PayOrderRequest struct {
	Return_code string `xml:"return_code"` //返回状态码
	Return_msg  string `xml:"return_msg"`  //返回信息
}

func WxPayCallBack(c *gin.Context) {
	res := models.Result{}

	respBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		res.Res = 1
		res.Msg = "解析返回body错误" + err.Error()
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	}

	Payres := PayOrderResponse{}
	err = xml.Unmarshal(respBytes, &Payres)

	if Payres.Return_code == "FAIL" || Payres.Result_code == "FAIL" {
		res.Res = 1
		res.Msg = "微信支付失败, 原因:" + Payres.Return_msg + Payres.Err_code + Payres.Err_code_des
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	}

	Req := PayOrderRequest{}
	Req.Return_code = "SUCCESS"
	Req.Return_msg = "OK"
	bytes_req, err := xml.Marshal(Req)
	if err != nil {
		res.Res = 1
		res.Msg = "解析XML报文失败, 原因:" + err.Error()
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	}

	str_req := string(bytes_req)
	//wxpay的unifiedorder接口需要http body中xmldoc的根节点是<xml></xml>这种，所以这里需要replace一下
	str_req = strings.Replace(str_req, "XUnifyOrderReq", "xml", -1)
	bytes_req = []byte(str_req)

	c.JSON(http.StatusOK, bytes.NewReader(bytes_req))
}
