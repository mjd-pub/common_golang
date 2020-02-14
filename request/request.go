package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mjd-pub/common_golang/utils"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Response map[string]interface{}
type Body map[string]interface{}

// HttpRequest 实现网络请求
func HttpRequest(method string, uri string, body Body, header map[string]string, timeOut int64) Response {
	var err error
	var result Response     //返回结果
	var resp *http.Response //http返回
	var reqUrl string
	start := time.Now()
	client := &http.Client{
		Timeout: time.Second * time.Duration(timeOut),
	}
	method = strings.ToUpper(method)
	switch method {
	case "GET":
		urlParams := url.Values{}
		for k, v := range body {
			if k == "query" || k == "dconf" {
				v, _ = json.Marshal(v)
				v = string(v.([]byte))
			}
			val := fmt.Sprintf("%s", v)
			urlParams.Add(k, val)
		}
		reqUrl = uri + "?" + urlParams.Encode()
		request, err := http.NewRequest("GET", reqUrl, nil)
		if err != nil {
			panic(err)
		}
		resp, err = client.Do(request)
		if err != nil {
			panic(err)
		}
	case "POST":
		dconf, _ := json.Marshal(body["dconf"])
		apiType, _ := body["type"].(string)
		qid, _ := body["qid"].(string)
		uid, _ := body["uid"].(string)
		csuid, _ := body["csuid"].(string)
		ptid, _ := body["ptid"].(string)
		bodyData, _ := json.Marshal(body) //格式化body
		reader := bytes.NewReader(bodyData)
		reqUrl = uri + fmt.Sprintf("?type=%s&qid=%s&uid=%s&dconf=%s&csuid=%s&ptid=%s",
			apiType, qid, uid, dconf, csuid, ptid)
		request, err := http.NewRequest("POST", reqUrl, reader)
		if err != nil {
			panic(err)
		}
		if len(header) <= 0 {
			request.Header.Set("Content-Type", "application/json;charset=UTF-8")
		} else {
			for k, v := range header {
				request.Header.Set(k, v)
			}
		}
		resp, err = client.Do(request)
		if err != nil {
			panic(err)
		}
	default:
		panic(errors.New("inner unsupported "+method))
	}
	defer func(err error) {
		var (
			ptid, csuid, role, uid, tid, oriIp, destIp, response, reqBody string
			httpCode                                                      int
		)
		ptid, ok := body["ptid"].(string)
		if !ok {
			ptid = ""
		}
		csuid, ok = body["csuid"].(string)
		if !ok {
			csuid = ""
		}
		role, ok = body["role"].(string)
		if !ok {
			role = ""
		}
		uid, ok = body["uid"].(string)
		if !ok {
			uid = ""
		}
		tid, ok = body["tid"].(string)
		if !ok {
			uid = ""
		}
		oriIp, _ = utils.GetLocalHostIp()
		destIp = resp.Request.Host
		if err != nil && result == nil {
			response = err.Error()
			httpCode = 404
		} else {
			httpCode = resp.StatusCode
			respByte, _ := json.Marshal(result)
			response = string(respByte)
		}
		if method == "POST" {
			resquest, _ := json.Marshal(body["query"])
			reqBody = string(resquest)
		}
		log := map[string]interface{}{
			"key":      "MJ_MD_INNER_NGINX_LOG",
			"role":     role,
			"code":     httpCode,
			"qid":      body["qid"],
			"tid":      tid,
			"type":     body["type"],
			"uid":      uid,
			"csuid":    csuid,
			"ptid":     ptid,
			"ori_ip":   oriIp,
			"dest_ip":  destIp,
			"t":        time.Now().Format("2006-01-02T15:04:05+08:00"),
			"req":      reqUrl,
			"resp":     response,
			"req_body": reqBody,
			"cost":     time.Since(start).Seconds(),
		}
		logByte, _ := json.Marshal(log)
		fmt.Printf("\n%s\n", string(logByte))
		resp.Body.Close()
	}(err)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic(errors.New("状态码错误："+strconv.Itoa(resp.StatusCode)))
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		panic(err)
	}

	return result
}


