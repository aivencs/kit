package request

import (
	"context"
	"crypto/tls"
	"errors"
	"net/url"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/go-resty/resty/v2"
)

var once sync.Once
var request Request

type Request interface {
	Get(ctx context.Context, link string, opt RequestOption) (string, error)
	Post(ctx context.Context, link string, opt RequestOption) (string, error)
}

type RequestOption struct {
	Trace            string
	Timeout          int
	Proxy            string
	EnableSkipVerify bool
	EnableHeader     bool
	Payload          []byte
}

type Resty struct{}

func InitRequest(name string) {
	once.Do(func() {
		switch name {
		case "resty":
			request = NewResty()
		default:
			request = NewResty()
		}
	})
}

func NewResty() Request {
	return &Resty{}
}

func (c *Resty) Get(ctx context.Context, link string, opt RequestOption) (string, error) {
	return c.work(ctx, "GET", link, opt)
}

func (c *Resty) Post(ctx context.Context, link string, opt RequestOption) (string, error) {
	return c.work(ctx, "POST", link, opt)
}

func (c *Resty) work(ctx context.Context, method string, link string, opt RequestOption) (string, error) {
	var response *resty.Response
	var err error
	serviceSafeString, _ := url.Parse(link)
	client := resty.New()
	client.SetHeaders(map[string]string{"Trace-ID": opt.Trace}) // set trace
	// apply option
	if opt.Timeout > 0 {
		client.SetTimeout(time.Duration(opt.Timeout) * time.Second)
	}
	if opt.EnableSkipVerify {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	if opt.EnableHeader {
		client.SetHeaders(map[string]string{
			"Trace-ID":   opt.Trace,
			"Host":       serviceSafeString.Host,
			"Referer":    serviceSafeString.Host,
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36",
		})
	}
	// set proxy
	if utf8.RuneCountInString(opt.Proxy) > 6 {
		client.SetProxy(opt.Proxy)
	}
	// set body

	switch method {
	case "GET":
		response, err = client.R().SetBody(opt.Payload).Get(link)
	case "POST":
		client.SetHeaders(map[string]string{
			"Content-Type": "application/json",
		})
		response, err = client.R().SetBody(opt.Payload).Post(link)
	default:
		response, err = client.R().SetBody(opt.Payload).Get(link)
	}
	if err != nil {
		return "", err
	}
	// status code
	if response.RawResponse.StatusCode < 299 {
		switch response.RawResponse.StatusCode {
		case 429:
			err = errors.New("并发超限")
		case 404:
			err = errors.New("资源不存在")
		case 200:
			err = nil
		case 201:
			err = nil
		default:
			err = errors.New("非正常状态码")
		}
	}
	return response.String(), err
}

func Get(ctx context.Context, link string, opt RequestOption) (string, error) {
	return request.Get(ctx, link, opt)
}
func Post(ctx context.Context, link string, opt RequestOption) (string, error) {
	return request.Post(ctx, link, opt)
}
