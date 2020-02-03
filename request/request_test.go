/**
* @Author: zhangjian@mioji.com
* @Date: 2019/9/10 下午4:06
 */
package request

import (
	"testing"
)

func TestHttpRequest(t *testing.T) {
	body := map[string]interface{}{
		"type": "c221",
		"qid":  "12345566",
		"query": map[string]interface{}{
			"cn": "铁函函",
		},
		"role":  "o",
		"uid":   "",
		"csuid": "zhangjian",
		"ptid":  "ptid",
		"tid":   "tid",
		"dconf": map[string]interface{}{},
	}
	header := map[string]string{
		"Content-Type": "application/json;charset=UTF-8",
	}
	HttpRequest("POST", "http://127.0.0.1:8081", body, header, 10)
	HttpRequest("GeT", "http://127.0.0.1:8081", body, nil, 10)
}
