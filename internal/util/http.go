package util

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/proxy"
)

// NewClient
// 新建HTTP请求客户端
func NewClient() *Client {
	return &Client{
		Client: &http.Client{
			Transport: &http.Transport{
				DisableKeepAlives: true,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

type Client struct {
	*http.Client
}

// SetProxy 设置代理
func (this *Client) SetProxy(u string) error {
	if transport, ok := this.Client.Transport.(*http.Transport); ok {
		//为空表示取消代理
		if len(u) == 0 {
			transport.Proxy = nil
			transport.DialContext = nil
			return nil
		}
		proxyUrl, err := url.Parse(u)
		if err != nil {
			transport.Proxy = nil
			return err
		}
		switch proxyUrl.Scheme {
		case "socks5", "socks5h":
			dialer, err := proxy.FromURL(proxyUrl, this)
			if err != nil {
				return err
			}
			transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			}
		default: //"http", "https"
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
		return nil
	}
	return fmt.Errorf("http.Transport类型错误: 预期(*http.Transport),得到(%T)", this.Client.Transport)
}

// SetTimeout 设置请求超时时间
// 下载大文件的时候需要设置长的超时时间
func (this *Client) SetTimeout(t time.Duration) *Client {
	this.Client.Timeout = t
	return this
}

func (this *Client) Dial(network, addr string) (net.Conn, error) {
	d := &net.Dialer{
		Timeout:   this.Client.Timeout,
		KeepAlive: this.Client.Timeout,
	}
	return d.Dial(network, addr)
}

func (this *Client) GetToFile(url string, filename string) (int64, error) {
	resp, err := this.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	w, err := os.Create(filename + ".downloading")
	if err != nil {
		return 0, err
	}
	n, err := io.Copy(w, resp.Body)
	if err != nil {
		w.Close()
		return n, err
	}
	w.Close()
	<-time.After(time.Millisecond * 100)
	return n, os.Rename(filename+".downloading", filename)
}
