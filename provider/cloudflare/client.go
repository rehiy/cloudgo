package cloudflare

import (
	"context"
	"strconv"

	"github.com/rehiy/cloudgo/provider"

	cf "github.com/cloudflare/cloudflare-go"
)

type Client struct {
	*provider.ReqeustParam
	Ctx context.Context
}

func NewClient(rq *provider.ReqeustParam) *Client {

	c := &Client{rq, context.Background()}

	c.NewApi()

	return c

}

func (c *Client) NewApi() (*cf.API, error) {

	return cf.NewWithAPIToken(c.SecretKey)

}

// 处理错误

func (c *Client) Error(err any) *provider.ResponseError {

	if er, ok := err.(*cf.Error); ok {
		code := er.Messages[0].Code
		return &provider.ResponseError{
			Code:    strconv.Itoa(code),
			Message: er.Messages[0].Message,
		}
	}

	re := &provider.ResponseError{}
	re.Create(err)

	return re

}
