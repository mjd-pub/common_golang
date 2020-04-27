package wechat

import (
	"encoding/xml"
	"errors"
	"github.com/mjd-pub/common_golang/utils"
	"io/ioutil"
	"strconv"
)

// CompanyPay 企业支付
type CompanyPay struct {
	wechatPay *wechatPay
}

// CompanyPayRequest 企业支付请求
type CompanyPayRequest struct {
	MchAppid       string `json:"mch_appid" xml:"mch_appid" structs:"mch_appid"`
	Mchid          string `json:"mchid" xml:"mchid" structs:"mchid"`
	NonceStr       string `json:"nonce_str" xml:"nonce_str" structs:"nonce_str"`
	PartnerTradeNo string `json:"partner_trade_no" xml:"partner_trade_no" structs:"partner_trade_no"`
	Openid         string `json:"openid" xml:"openid" structs:"openid"`
	CheckName      string `json:"check_name" xml:"check_name" structs:"check_name"`
	Amount         int    `json:"amount" xml:"amount" structs:"amount"`
	Desc           string `json:"desc" xml:"amount" structs:"desc"`
}

// CompanyPayRequest 企业支付查询请求
type CompanyPayQueryRequest struct {
	NonceStr       string `json:"nonce_str" xml:"nonce_str" structs:"nonce_str"`
	Sign           string `json:"sign" xml:"sign" structs:"sign"`
	PartnerTradeNo string `json:"partner_trade_no" xml:"partner_trade_no" structs:"partner_trade_no"`
	MchId          string `json:"mch_id" xml:"mch_id" structs:"mch_id"`
	Appid          string `json:"appid" xml:"appid" structs:"appid"`
}

// CompanyPayRequest 企业支付查询返回
type CompanyPayQueryResponse struct {
	ReturnCode     string `json:"return_code" xml:"return_code" structs:"return_code"`
	ReturnMsg      string `json:"return_msg" xml:"return_msg" structs:"return_msg"`
	ResultCode     string `json:"result_code" xml:"result_code" structs:"result_code"`
	ErrCode        string `json:"err_code" xml:"err_code" structs:"err_code"`
	ErrCodeDes     string `json:"err_code_des" xml:"err_code_des" structs:"err_code_des"`
	PartnerTradeNo string `json:"partner_trade_no" xml:"partner_trade_no" structs:"partner_trade_no"`
	Appid          string `json:"appid" xml:"appid" structs:"appid"`
	MchId          string `json:"mch_id" xml:"mch_id" structs:"mch_id"`
	DetailId       string `json:"detail_id" xml:"detail_id" structs:"detail_id"`
	Status         string `json:"status" xml:"status" structs:"status"`
	Reason         string `json:"reason" xml:"reason" structs:"reason"`
	Openid         string `json:"openid" xml:"openid" structs:"openid"`
	TransferName   string `json:"transfer_name" xml:"transfer_name" structs:"transfer_name"`
	PaymentAmount  int    `json:"payment_amount" xml:"payment_amount" structs:"payment_amount"`
	TransferTime   string `json:"transfer_time" xml:"transfer_time" structs:"transfer_time"`
	PaymentTime    string `json:"payment_time" xml:"payment_time" structs:"payment_time"`
	Desc           string `json:"desc" xml:"desc" structs:"desc" structs:"desc"`
}

// CompanyPayResponse 企业支付查询返回
type CompanyPayResponse struct {
	ReturnCode     string `json:"return_code" xml:"return_code" structs:"return_code"`
	ReturnMsg      string `json:"return_msg" xml:"return_msg" structs:"return_msg"`
	MchAppid       string `json:"mch_appid" xml:"mch_appid" structs:"mch_appid"`
	Mchid          string `json:"mchid" xml:"mchid" structs:"mchid"`
	DeviceInfo     string `json:"device_info" xml:"device_info" structs:"device_info"`
	NonceStr       string `json:"nonce_str" xml:"nonce_str" structs:"nonce_str"`
	ResultCode     string `json:"result_code" xml:"result_code" structs:"result_code"`
	ErrCode        string `json:"err_code" xml:"err_code" structs:"err_code"`
	ErrCodeDes     string `json:"err_code_des" xml:"err_code_des" structs:"err_code_des"`
	PartnerTradeNo string `json:"partner_trade_no" xml:"partner_trade_no" structs:"partner_trade_no"`
	PaymentNo      string `json:"payment_no" xml:"payment_no" structs:"payment_no"`
	PaymentTime    string `json:"payment_time" xml:"payment_time" structs:"payment_time"`
}

// NewCompanyPayClient 构造基础连接
func NewCompanyPayClient(appid, mchid, key, apiclientKey, apiclientCert string) (companyPay *CompanyPay) {
	wechatPay := NewWechatPay(appid, mchid, key, apiclientKey, apiclientCert)
	companyPay = &CompanyPay{
		wechatPay: wechatPay,
	}
	return
}

/**
 * NewCompanyPayRequest 构造下单请求
 * @params companyPay
 * @params businessId 业务订单号
 * @params amount:int 金额(单位:分)
 * @return CompanyPayRequest
 */
func (companyPay *CompanyPay) NewCompanyPayRequest(businessId string, amount int, openid string, desc string) CompanyPayRequest {
	return CompanyPayRequest{
		MchAppid:       companyPay.wechatPay.appid,
		Mchid:          companyPay.wechatPay.mchid,
		NonceStr:       utils.GetNonceStr(),
		PartnerTradeNo: businessId,
		Openid:         openid,
		CheckName:      "NO_CHECK", // 不校验姓名
		Amount:         amount,
		Desc:           desc,
	}
}

/**
 * NewCompanyPayQueryRequest 构造查询请求
 * @params companyPay
 * @params businessId 业务订单号
 * @return CompanyPayQueryRequest
 */
func (companyPay *CompanyPay) NewCompanyPayQueryRequest(businessId string) CompanyPayQueryRequest {
	return CompanyPayQueryRequest{
		NonceStr:       utils.GetNonceStr(),
		PartnerTradeNo: businessId,
		MchId:          companyPay.wechatPay.mchid,
		Appid:          companyPay.wechatPay.appid,
	}
}

// Pay 发起支付
func (c *CompanyPay) Pay(request CompanyPayRequest) (queryResponse *CompanyPayResponse, err error) {
	// 向微信发送请求
	resp, err := c.wechatPay.Request(COMPANY_PAY, request)
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

// Query 企业支付查询
func (c *CompanyPay) Query(request CompanyPayQueryRequest) (queryResponse *CompanyPayQueryResponse, err error) {
	// 向微信发送请求
	resp, err := c.wechatPay.Request(COMPANY_PAY_QUERY, request)
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
