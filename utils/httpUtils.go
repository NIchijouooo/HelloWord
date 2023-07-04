package utils

import (
	"bytes"
	"crypto/tls"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

// SendPostTls 发送HTTPS的POST请求
func SendPostTls(url string, body string, header map[string]string) (string, error) {

	// 创建一个自定义的TLS配置
	tlsConfig := &tls.Config{
		// 忽略证书验证
		InsecureSkipVerify: true,
	}

	// 创建一个自定义的HTTP客户端
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// 发送POST请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	if header != nil {
		for key := range header {
			req.Header.Set(key, header[key])
		}
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("请求错误")
	}

	// 读取响应体的内容
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(responseData), nil
}
