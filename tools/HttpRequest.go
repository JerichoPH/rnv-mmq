package tools

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"rnv-mmq/wrongs"
)

const (
	HttpRequestDataTypeGet       = ""
	HttpRequestDataTypeUrlEncode = "HTTP-REQUEST-DATA-TYPE-URL-ENCODE"
	HttpRequestDataTypeJson      = "HTTP-REQUEST-DATA-TYPE-JSON"
)

type (
	// HttpRequest http请求工具
	HttpRequest struct {
		Method, Url, DataType string
		Body                  io.Reader
		ErrorHandle           func(resp *http.Response, err error)
	}
)

func (receiver HttpRequest) Send() (int, string) {
	var (
		err  error
		req  *http.Request
		resp *http.Response
	)

	switch receiver.Method {
	case http.MethodGet:
		if req, err = http.NewRequest(receiver.Method, receiver.Url, receiver.Body); err != nil {
			wrongs.ThrowForbidden("初始化请求错误：%s", err.Error())
		}
		if receiver.DataType == HttpRequestDataTypeJson {
			req.Header.Add("Accept", "application/json")
		}
	case http.MethodPost:
	case http.MethodPut:
	case http.MethodDelete:
		switch receiver.DataType {
		case HttpRequestDataTypeUrlEncode:
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		case HttpRequestDataTypeJson:
			req.Header.Add("Content-Type", "application/json")
		}
	}

	if resp, err = http.DefaultClient.Do(req); err != nil {
		// 网络层面错误
		receiver.ErrorHandle(resp, err)
		return 0, ""
	} else {
		defer func(Body io.ReadCloser) {
			if err = Body.Close(); err != nil {
				wrongs.ThrowForbidden("读取响应结果错误：" + err.Error())
			}
		}(resp.Body)
		if result, err := ioutil.ReadAll(resp.Body); err != nil {
			wrongs.ThrowForbidden("解析响应结果错误：" + err.Error())
		} else {
			return resp.StatusCode, string(result)
		}
	}
	return 0, ""
}

// GetTargetUrl 获取新 url
func GetTargetUrl(ctx *gin.Context) string {
	var err error

	u := url.URL{}

	// 解析url参数
	err = ctx.Request.ParseForm()
	if err != nil {
		println(fmt.Sprintf("解析参数错误错误: %s", err.Error()))
	}

	// 重新拼装参数(去掉 target)
	urlValues := url.Values{}
	for k, values := range ctx.Request.Form {
		if k != "target" {
			for _, value := range values {
				urlValues.Add(k, value)
			}
		}
	}
	u.RawQuery = urlValues.Encode()

	// 拼装转发目标 url
	return fmt.Sprintf("%s%s", ctx.Query("target"), u.String())
}

// GetBody 解析 body
func GetBody(ctx *gin.Context) string {
	formData, _ := ioutil.ReadAll(ctx.Request.Body)
	return string(formData)
}
