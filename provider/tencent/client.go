package tencent

import (
	"regexp"
	"strings"

	"github.com/rehiy/cloudgo/provider"
	"github.com/rehiy/cloudgo/setting"

	tc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	te "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	th "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	tp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

type Client struct {
	*provider.ReqeustParam
	credential *tc.Credential
	profile    *tp.ClientProfile
}

func NewClient(rq *provider.ReqeustParam) *Client {

	c := &Client{rq, nil, nil}

	c.NewCredential()
	c.NewProfile()

	return c

}

func (c *Client) NewCredential() {

	// 初始化
	credential := tc.NewCredential(c.SecretId, c.SecretKey)

	// 回传参数
	c.credential = credential

}

func (c *Client) NewProfile() {

	// 初始化
	profile := tp.NewClientProfile()

	// 调试模式
	profile.Debug = setting.Debug

	// 网络错误重试
	profile.NetworkFailureMaxRetries = 2

	// API 限频重试
	profile.RateLimitExceededMaxRetries = 2

	// 地域容灾机制
	profile.DisableRegionBreaker = false
	profile.BackupEndpoint = "ap-hongkong." + th.RootDomain

	// 按地域设置接口
	if c.Endpoint != "" {
		profile.HttpProfile.Endpoint = c.Endpoint // 完整域名
	} else if c.RegionId != "" {
		if !strings.HasSuffix(c.RegionId, "-ec") {
			profile.HttpProfile.Endpoint = c.Service + "." + c.RegionId + "." + th.RootDomain
		}
	}

	// 回传参数
	c.profile = profile

}

// 处理错误

func (c *Client) Error(err any) *provider.ResponseError {

	if er, ok := err.(*te.TencentCloudSDKError); ok {
		exp := regexp.MustCompile(`\[request id:.+\]`)
		ret := strings.Split(exp.ReplaceAllString(er.Message, ""), "\n")[0]
		return &provider.ResponseError{
			Code:    er.Code,
			Message: ret,
		}
	}

	re := &provider.ResponseError{}
	re.Create(err)

	return re

}
