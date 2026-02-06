package alibaba

import (
	"os"
	"regexp"

	"github.com/rehiy/cloudgo/provider"
	"github.com/rehiy/cloudgo/setting"

	ac "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	au "github.com/alibabacloud-go/tea-utils/v2/service"
	tea "github.com/alibabacloud-go/tea/tea"
)

type Client struct {
	*provider.ReqeustParam
	config  *ac.Config
	runtime *au.RuntimeOptions
}

func NewClient(rq *provider.ReqeustParam) *Client {

	c := &Client{rq, nil, nil}

	c.NewConfig()
	c.NewRuntime()

	return c

}

func (c *Client) NewConfig() {

	if setting.Debug {
		os.Setenv("DEBUG", "tea")
	}

	config := &ac.Config{
		AccessKeyId:     tea.String(c.SecretId),
		AccessKeySecret: tea.String(c.SecretKey),
		RegionId:        tea.String(c.RegionId),
	}

	// 回传参数
	c.config = config

}

func (c *Client) NewRuntime() {

	runtime := &au.RuntimeOptions{
		// 自动重试机制
		Autoretry:   tea.Bool(true),
		MaxAttempts: tea.Int(2),

		// 超时配置（单位 ms）
		ConnectTimeout: tea.Int(5000),
		ReadTimeout:    tea.Int(10000),
	}

	// 回传参数
	c.runtime = runtime

}

// 处理错误

func (c *Client) Error(err any) *provider.ResponseError {

	if er, ok := err.(*tea.SDKError); ok {
		exp := regexp.MustCompile(`^code: \d+, (.+) request id.+$`)
		msg := exp.ReplaceAllString(*er.Message, "$1")
		return &provider.ResponseError{
			Code:    *er.Code,
			Message: msg,
		}
	}

	re := &provider.ResponseError{}
	re.Create(err)

	return re

}
