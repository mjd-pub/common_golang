package wechat

import (
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"github.com/mjd-pub/common_golang/utils"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	UNIFIED_ORDER     = "https://api.mch.weixin.qq.com/pay/unifiedorder"                      // 统一下单接口地址
	ORDER_QUERY       = "https://api.mch.weixin.qq.com/pay/orderquery"                        // 查询接口地址
	REFUND            = "https://api.mch.weixin.qq.com/secapi/pay/refund"                     // 退款接口地址
	REFEUN_QUERY      = "https://api.mch.weixin.qq.com/pay/refundquery"                       // 退款查询接口地址
	COMPANY_PAY       = "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers" // 企业支付下单
	COMPANY_PAY_QUERY = "https://api.mch.weixin.qq.com/mmpaymkttransfers/gettransferinfo"     // 企业支付查询
)

const (
	DEFAULT        = 0
	PAY_SUCCESS    = 1
	REFUND_PROCESS = 11
	REFUND_SUCCESS = 12
	REFUND_FAIL    = 13
)

//　WechatPay 微信支付基础结构体
type wechatPay struct {
	apiclientCert string
	apiclientKey  string
	appid         string
	key           string
	mchid         string
}

// RefundRequests 微信申请退款请求参数
type RefundRequests struct {
	Appid         string `json:"appid" xml:"appid" structs:"appid"`
	MchId         string `json:"mch_id" xml:"mch_id" structs:"mch_id"`
	TransactionId string `json:"transaction_id" xml:"transaction_id" structs:"transaction_id"`
	OutTradeNo    string `json:"out_trade_no" xml:"out_trade_no" structs:"out_trade_no"`
	NonceStr      string `json:"nonce_str" xml:"nonce_str" structs:"nonce_str"`
	SignType      string `json:"sign_type" xml:"sign_type" structs:"sign_type"`
	OutRefundNo   string `json:"out_refund_no" xml:"out_refund_no" structs:"out_refund_no"`
	TotalFee      int    `json:"total_fee" xml:"total_fee" structs:"total_fee"`
	RefundFee     int    `json:"refund_fee" xml:"refund_fee" structs:"refund_fee"`
	RefundFeeType string `json:"refund_fee_type" xml:"refund_fee_type" structs:"refund_fee_type"`
	RefundDesc    string `json:"refund_desc" xml:"refund_desc" structs:"refund_desc"`
	RefundAccount string `json:"refund_account" xml:"refund_account" structs:"refund_account"`
	NotifyUrl     string `json:"notify_url" xml:"notify_url" structs:"notify_url"`
}

// RefundRespones 微信申请退款请求返回参数
type RefundRespones struct {
	ReturnCode          string `json:"return_code,omitempty" xml:"return_code,omitempty" structs:"return_code"`
	ReturnMsg           string `json:"return_msg,omitempty" xml:"return_msg,omitempty" structs:"return_msg"`
	ResultCode          string `json:"result_code,omitempty" xml:"result_code,omitempty" structs:"result_code"`
	ErrCode             string `json:"err_code,omitempty" xml:"err_code,omitempty" structs:"err_code"`
	ErrCodeDes          string `json:"err_code_des,omitempty" xml:"err_code_des,omitempty" structs:"err_code_des"`
	Appid               string `json:"appid,omitempty" xml:"appid,omitempty" structs:"appid"`
	MchId               string `json:"mch_id,omitempty" xml:"mch_id,omitempty" structs:"mch_id"`
	NonceStr            string `json:"nonce_str,omitempty" xml:"nonce_str,omitempty" structs:"nonce_str"`
	Sign                string `json:"sign,omitempty" xml:"sign,omitempty" structs:"sign"`
	TransactionId       string `xml:"transaction_id,omitempty" json:"transaction_id,omitempty" structs:"transaction_id"`
	OutTradeNo          string `xml:"out_trade_no,omitempty" json:"out_trade_no,omitempty" structs:"out_trade_no"`
	OutRefundNo         string `xml:"out_refund_no,omitempty" json:"out_refund_no,omitempty" structs:"out_refund_no"`
	RefundId            string `json:"refund_id,omitempty" xml:"refund_id,omitempty" structs:"refund_id"`
	RefundFee           int    `json:"refund_fee,omitempty" xml:"refund_fee,omitempty" structs:"refund_fee"`
	SettlementTotalFree int    `json:"settlement_total_free,omitempty" xml:"settlement_total_free,omitempty" structs:"settlement_total_free"`
	FreeType            string `json:"free_type,omitempty" xml:"free_type,omitempty" structs:"free_type"`
	CashFee             int    `xml:"cash_fee,omitempty" json:"cash_fee,omitempty" structs:"cash_fee"`
	CashFeeType         string `xml:"cash_fee_type,omitempty" json:"cash_fee_type,omitempty" structs:"cash_fee_type"`
	CashRefundFee       int    `json:"cash_refund_fee,omitempty" xml:"cash_refund_fee,omitempty" structs:"cash_refund_fee"`
	CouponType0         string `json:"coupon_type_0,omitempty" xml:"coupon_type_0,omitempty" structs:"coupon_type_0"`
	CouponRefundFee     int    `json:"coupon_refund_fee,omitempty" xml:"coupon_refund_fee,omitempty" structs:"coupon_refund_fee"`
	CouponRefundFee0    int    `json:"coupon_refund_fee_0,omitempty" xml:"coupon_refund_fee_0,omitempty" structs:"coupon_refund_fee_0"`
	ConponRefundCount   int    `json:"conpon_refund_count,omitempty" xml:"conpon_refund_count,omitempty" structs:"conpon_refund_count"`
	ConponRefundId0     string `json:"conpon_refund_id_0,omitempty" xml:"conpon_refund_id_0,omitempty" structs:"conpon_refund_id_0"`
}

type RefundQueryRequests struct {
	AppId         string `json:"app_id" xml:"app_id" structs:"app_id"`
	MchId         string `json:"mch_id" xml:"mch_id" structs:"mch_id"`
	NonceStr      string `json:"nonce_str" xml:"nonce_str" structs:"nonce_str"`
	SignType      string `json:"sign_type" xml:"sign_type" structs:"sign_type"`
	TransactionId string `json:"transaction_id" xml:"transaction_id" structs:"transaction_id"`
	OutTradeNo    string `json:"out_trade_no" xml:"out_trade_no" structs:"out_trade_no"`
	OutRefundNo   string `json:"out_refund_no" xml:"out_refund_no" structs:"out_refund_no"`
	RefundId      string `json:"refund_id" xml:"refund_id" structs:"refund_id"`
	Offset        int    `json:"offset" xml:"offset" structs:"offset"`
}

type RefundQueryRespones struct {
	ReturnCode           string `json:"return_code,omitempty" xml:"return_code,omitempty" structs:"return_code"`
	ReturnMsg            string `json:"return_msg,omitempty" xml:"return_msg,omitempty" structs:"return_msg"`
	ResultCode           string `json:"result_code,omitempty" xml:"result_code,omitempty" structs:"result_code"`
	ErrCode              string `json:"err_code,omitempty" xml:"err_code,omitempty" structs:"err_code"`
	ErrCodeDes           string `json:"err_code_des,omitempty" xml:"err_code_des,omitempty" structs:"err_code_des"`
	Appid                string `json:"appid,omitempty" xml:"appid,omitempty" structs:"appid"`
	MchId                string `json:"mch_id,omitempty" xml:"mch_id,omitempty" structs:"mch_id"`
	NonceStr             string `json:"nonce_str,omitempty" xml:"nonce_str,omitempty" structs:"nonce_str"`
	Sign                 string `json:"sign,omitempty" xml:"sign,omitempty" structs:"sign"`
	TotalRefundCount     int    `json:"total_refund_count,omitempty" xml:"total_refund_count,omitempty" structs:"total_refund_count"`
	TransactionId        string `xml:"transaction_id,omitempty" json:"transaction_id,omitempty" structs:"transaction_id"`
	OutTradeNo           string `xml:"out_trade_no,omitempty" json:"out_trade_no,omitempty" structs:"out_trade_no"`
	TotalFee             int    `json:"total_fee,omitempty" xml:"total_fee,omitempty" structs:"total_fee"`
	SettlementTotalFree  int    `json:"settlement_total_free,omitempty" xml:"settlement_total_free,omitempty" structs:"settlement_total_free"`
	FreeType             string `json:"free_type,omitempty" xml:"free_type,omitempty" structs:"free_type"`
	CashFee              int    `xml:"cash_fee,omitempty" json:"cash_fee,omitempty" structs:"cash_fee"`
	RefundCount          int    `json:"refund_count,omitempty" xml:"refund_count,omitempty" structs:"refund_count"`
	OutRefundNo0         string `xml:"out_refund_no_0,omitempty" json:"out_refund_no_0,omitempty" structs:"out_refund_no_0"`
	RefundId0            string `json:"refund_id_0,omitempty" xml:"refund_id_0,omitempty" structs:"refund_id_0"`
	RefundFee0           int    `json:"refund_fee_0,omitempty" xml:"refund_fee_0,omitempty" structs:"refund_fee_0"`
	SettleMentRefundFee0 int    `json:"settle_ment_refund_fee_0" xml:"settle_ment_refund_fee_0" structs:"settle_ment_refund_fee_0"`
	CouponType00         string `json:"coupon_type_0_0" xml:"coupon_type_0_0" structs:"coupon_type_0_0"`
	ConponRefundFee0     int    `json:"conpon_refund_fee_0" xml:"conpon_refund_fee_0" structs:"conpon_refund_fee_0"`
	ConponRefundCount0   int    `json:"conpon_refund_count_0" xml:"conpon_refund_count_0" structs:"conpon_refund_count_0"`
	ConponRefundId00     string `json:"conpon_refund_id_0_0" xml:"conpon_refund_id_0_0" structs:"conpon_refund_id_0_0"`
	ConponRefundFee00    string `json:"conpon_refund_fee_0_0" xml:"conpon_refund_fee_0_0" structs:"conpon_refund_fee_0_0"`
	RefundStatus0        string `json:"refund_status_0" xml:"refund_status_0" structs:"refund_status_0"`
	RefundAccount0       string `json:"refund_account_0" xml:"refund_account_0" structs:"refund_account_0"`
	RefundRecvAccount0   string `json:"refund_recv_account_0" xml:"refund_recv_account_0" structs:"refund_recv_account_0"`
	RefundSuccessTime0   string `json:"refund_success_time_0" xml:"refund_success_time_0" structs:"refund_success_time_0"`
}

// PayNotifyRequest 支付回调
type PayNotifyRequest struct {
	ReturnCode         string `json:"return_code" xml:"return_code" structs:"return_code"`
	ReturnMsg          string `json:"return_msg" xml:"return_msg" structs:"return_msg"`
	Appid              string `json:"appid" xml:"appid" structs:"appid"`
	MchId              string `json:"mch_id" xml:"mch_id"  structs:"mch_id"`
	DeviceInfo         string `json:"device_info" xml:"device_info" structs:"device_info"`
	NonceStr           string `json:"nonce_str" xml:"nonce_str" structs:"nonce_str"`
	Sign               string `json:"sign" xml:"sign" structs:"sign"`
	SignType           string `json:"sign_type" xml:"sign_type" structs:"sign_type"`
	ResultCode         string `json:"result_code" xml:"result_code" structs:"result_code, omitempty"`
	ErrCode            string `json:"err_code" xml:"err_code" structs:"err_code, omitempty"`
	ErrCodeDes         string `json:"err_code_des" xml:"err_code_des" structs:"err_code_des, omitempty"`
	Openid             string `json:"openid" xml:"openid" structs:"openid, omitempty"`
	IsSubscribe        string `json:"is_subscribe" xml:"is_subscribe" structs:"is_subscribe, omitempty"`
	TradeType          string `json:"trade_type" xml:"trade_type" structs:"trade_type, omitempty"`
	BankType           string `json:"bank_type" xml:"bank_type" structs:"bank_type, omitempty"`
	TotalFee           int    `json:"total_fee" xml:"total_fee" structs:"total_fee, omitempty"`
	SettlementTotalFee int    `json:"settlement_total_fee" xml:"settlement_total_fee" structs:"settlement_total_fee, omitempty"`
	FeeType            string `json:"fee_type" xml:"fee_type" structs:"fee_type, omitempty"`
	CashFee            int    `json:"cash_fee" xml:"cash_fee" structs:"cash_fee, omitempty"`
	CashFeeType        string `json:"cash_fee_type" xml:"cash_fee_type" structs:"cash_fee_type, omitempty"`
	CouponFee          string `json:"coupon_fee" xml:"coupon_fee" structs:"coupon_fee, omitempty"`
	CouponCount        string `json:"coupon_count" xml:"coupon_count" structs:"coupon_count, omitempty"`
	TransactionId      string `json:"transaction_id" xml:"transaction_id" structs:"transaction_id, omitempty"`
	OutTradeNo         string `json:"out_trade_no" xml:"out_trade_no" structs:"out_trade_no, omitempty"`
	Attach             string `json:"attach" xml:"attach" structs:"attach, omitempty"`
	TimeEnd            string `json:"time_end" xml:"time_end" structs:"time_end, omitempty"`
}

type RefundNotifyRequest struct {
	ReturnCode          string `xml:"return_code,omitempty" json:"return_code,omitempty"`
	ReturnMsg           string `xml:"return_msg,omitempty" json:"return_msg,omitempty"`
	Appid               string `xml:"appid,omitempty" json:"appid,omitempty"`
	MchId               string `xml:"mch_id,omitempty" json:"mch_id,omitempty"`
	NonceStr            string `xml:"nonce_str,omitempty" json:"nonce_str,omitempty"`
	ReqInfo             string `xml:"req_info,omitempty" json:"req_info,omitempty"`
	UnmarshalReqInfo    RefundReqInfo  `json:"unmarshal_req_info" xml:"unmarshal_req_info"`
}

type RefundReqInfo struct {
	TransactionId       string `xml:"transaction_id,omitempty" json:"transaction_id,omitempty"`
	OutTradeNo          string `xml:"out_trade_no,omitempty" json:"out_trade_no,omitempty"`
	RefundId            string `xml:"refund_id,omitempty" json:"refund_id,omitempty"`
	OutRefundNo         string `xml:"out_refund_no,omitempty" json:"out_refund_no,omitempty"`
	TotalFee            int    `xml:"total_fee,omitempty" json:"total_fee,omitempty"`
	SettlementTotalFee  int    `xml:"settlement_total_fee,omitempty" json:"settlement_total_fee,omitempty"`
	RefundFee           int    `xml:"refund_fee,omitempty" json:"refund_fee,omitempty"`
	SettlementRefundFee int    `xml:"settlement_refund_fee,omitempty" json:"settlement_refund_fee,omitempty"`
	RefundStatus        string `xml:"refund_status,omitempty" json:"refund_status,omitempty"`
	SuccessTime         string `xml:"success_time,omitempty" json:"success_time,omitempty"`
	RefundRecvAccout    string `xml:"refund_recv_accout,omitempty" json:"refund_recv_accout,omitempty"`
	RefundAccount       string `xml:"refund_account,omitempty" json:"refund_account,omitempty"`
	RefundRequestSource string `xml:"refund_request_source,omitempty" json:"refund_request_source,omitempty"`
}

// ServiceNotifyResponse 服务器主动回复微信
type ServiceNotifyResponse struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
}

/**
 * newWechatPay 微信支付初始化
 * @params appid 商户号绑定的appid
 * @params mchid 商户号
 * @params key   支付密钥
 * @params apiclientKey
 * @params apiclientCert
 */
func NewWechatPay(appid, mchid, key, apiclientKey, apiclientCert string) *wechatPay {
	return &wechatPay{
		apiclientCert: apiclientCert,
		apiclientKey:  apiclientKey,
		appid:         appid,
		key:           key,
		mchid:         mchid,
	}
}

/**
 * Request 微信标准请求输出
 * @params uri 请求uri
 * @params requestData 请求参数为固定结构体
 *
 * @return resp 标准http请求 err 标准错误输出
 */
func (wechat *wechatPay) Request(uri string, requestData interface{}) (resp *http.Response, err error) {
	//2.签名和请求参数
	respXml, err := wechat.dealXmlRequest(requestData)
	if err != nil {
		return
	}
	//3. 证书和公钥文件
	cert := wechat.apiclientCert  //证书
	sslKey := wechat.apiclientKey //公钥
	//4. 带证书发送http请求
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(cert))
	//5.tls.X509KeyPair 直接读字符串
	cliCrt, err := tls.X509KeyPair([]byte(cert), []byte(sslKey))
	if err != nil {
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			//RootCAs:      pool,	不能添加这一项
			Certificates: []tls.Certificate{cliCrt},
			ClientAuth:   tls.RequireAndVerifyClientCert,
		},
	}
	client := &http.Client{Transport: tr}
	//请求微信接口
	resp, err = client.Post(uri, "text/xml; charset=UTF8", strings.NewReader(respXml))
	return
}

func (wechat *wechatPay) NewRefundRequests(outRefundNo, transactionId, businessId, notifyUrl string, totalFee, refundFee int) (request RefundRequests) {
	return RefundRequests{
		Appid:         wechat.appid,
		MchId:         wechat.mchid,
		TransactionId: transactionId,
		OutTradeNo:    businessId,
		NonceStr:      utils.GetNonceStr(),
		SignType:      "MD5",
		OutRefundNo:   outRefundNo,
		TotalFee:      totalFee,
		RefundFee:     refundFee,
		RefundFeeType: "CNY",
		RefundDesc:    "",
		RefundAccount: "",
		NotifyUrl:     notifyUrl,
	}
}

/**
 * Refund 小程序申请退款
 *
 * @params request AppletPayRefundRequests
 * @return AppletPayRefundRespones err
 */
func (wechatPay *wechatPay) Refund(request RefundRequests) (queryResponse *RefundRespones, err error) {
	queryResponse = new(RefundRespones)
	// 向微信发送请求
	resp, err := wechatPay.Request(REFUND, request)
	if err != nil {
		return nil, errors.New("请求异常:" + err.Error())
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("httpCode Err:" + strconv.Itoa(resp.StatusCode))
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//xml解码
	err = xml.Unmarshal(respData, queryResponse)
	if err != nil {
		return nil, err
	}
	return
}

/**
 * RefundQuery 小程序退款查询
 *
 * @params request AppletPayRefundQueryRequests
 * @return AppletPayRefundQueryRespones err
 */
func (wechatPay *wechatPay) RefundQuery(request RefundQueryRequests) (queryResponse *RefundQueryRespones, err error) {
	// 向微信发送请求
	resp, err := wechatPay.Request(REFEUN_QUERY, request)
	if err != nil {
		return nil, errors.New("请求异常:" + err.Error())
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("httpCode Err:" + strconv.Itoa(resp.StatusCode))
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//xml解码
	err = xml.Unmarshal(respData, queryResponse)
	if err != nil {
		return nil, err
	}
	return
}

// DealXmlRequest 处理初xml请求body
func (wechat *wechatPay) dealXmlRequest(params interface{}) (xmlRequest string, err error) {
	paramsMap := structs.Map(params)
	// 对数据进行签名
	sign, err := wechat.signData(paramsMap)
	if err != nil {
		return
	}
	paramsMap["sign"] = sign
	//4。转化为xml格式
	xmlRequest = utils.ToXml(paramsMap)
	return
}

// SignData 对数据进行签名
func (wechat *wechatPay) signData(data map[string]interface{}) (sign string, err error) {
	strs := utils.Ksort(data)
	//1.2 使用URL键值对的形式生成字符串
	str := utils.ToUrlParams(data, strs)
	//1.3 在str后加入KEY
	str = str + "&key=" + wechat.key
	//2. 将得到的数据做MD5结算得到signValue
	m := md5.New()
	_, err = io.WriteString(m, str)
	if err != nil {
		return "", errors.New("签名错误:" + err.Error())
	}
	//3。将签名转化为大写
	sign = strings.ToUpper(fmt.Sprintf("%x", m.Sum(nil)))
	return
}

func (wechat *wechatPay) ParsePayNotifyRequest(request *http.Request) (notifyReq *PayNotifyRequest, err error) {
	defer request.Body.Close()
	notifyReq = new(PayNotifyRequest)
	err = xml.NewDecoder(request.Body).Decode(notifyReq)
	SingData := map[string]interface{}{
		"return_code":          notifyReq.ReturnCode,
		"return_msg":           notifyReq.ReturnMsg,
		"appid":                notifyReq.Appid,
		"mch_id":               notifyReq.MchId,
		"device_info":          notifyReq.DeviceInfo,
		"nonce_str":            notifyReq.NonceStr,
		"sign":                 notifyReq.Sign,
		"sign_type":            notifyReq.SignType,
		"result_code":          notifyReq.ResultCode,
		"err_code":             notifyReq.ErrCode,
		"err_code_des":         notifyReq.ErrCodeDes,
		"openid":               notifyReq.Openid,
		"is_subscribe":         notifyReq.IsSubscribe,
		"trade_type":           notifyReq.TradeType,
		"bank_type":            notifyReq.BankType,
		"total_fee":            notifyReq.TotalFee,
		"settlement_total_fee": notifyReq.SettlementTotalFee,
		"fee_type":             notifyReq.FeeType,
		"cash_fee":             notifyReq.CashFee,
		"cash_fee_type":        notifyReq.CashFeeType,
		"coupon_fee":           notifyReq.CouponFee,
		"coupon_count":         notifyReq.CouponCount,
		"transaction_id":       notifyReq.TransactionId,
		"out_trade_no":         notifyReq.OutTradeNo,
		"attach":               notifyReq.Attach,
		"time_end":             notifyReq.TimeEnd,
	}
	sign, err := wechat.signData(SingData)
	if err != nil {
		return
	}
	if sign != notifyReq.Sign {
		return notifyReq, errors.New("验签失败:签名不一致")
	}
	return
}

func (wechat *wechatPay) ParseRefundRequest(request *http.Request) (notifyReq *RefundNotifyRequest, err error) {
	defer request.Body.Close()
	notifyReq = new(RefundNotifyRequest)
	err = xml.NewDecoder(request.Body).Decode(notifyReq)
	if err != nil {
		return
	}
	//解码数据
	respData, err := DecodeNotifyData(notifyReq.ReqInfo, wechat.key)
	if err != nil {
		return
	}
	reqInfo := new(RefundReqInfo)
	//解码数据
	err = xml.Unmarshal(respData, reqInfo)
	if err != nil {
		return
	}
	notifyReq.UnmarshalReqInfo = *reqInfo
}

func DecodeNotifyData(reqInfo string, paykey string) ([]byte, error) {
	//1.对加密串A做base64解码，得到加密串B
	b, err := base64.StdEncoding.DecodeString(reqInfo)
	if err != nil {
		return nil, err
	}
	//2.对商户key做md5，得到32位小写key
	err = utils.SetAesKey(strings.ToLower(utils.Md5(paykey)))
	if err != nil {
		return nil, err
	}

	//3.用key*对加密串B做AES-256-ECB解密（PKCS7Padding）
	plaintext, err := utils.AesECBDecrypt(b)

	return plaintext, nil
}
