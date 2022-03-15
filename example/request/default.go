package main

import (
	"context"
	"fmt"

	"github.com/aivencs/kit/pkg/request"
)

func main() {
	request.InitRequest("resty")
	link := "https://www.taobao.com/help/getip.php"
	ctx := context.WithValue(context.Background(), "trace", "ctx-request-001")
	r, err := request.Get(ctx, link, request.RequestOption{
		Trace:            "ctx-request-001",
		Timeout:          10,
		Proxy:            "",
		EnableSkipVerify: true,
		EnableHeader:     true,
	})
	fmt.Println(r, err)
}
