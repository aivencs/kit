package main

import (
	"context"
	"fmt"

	"github.com/aivencs/kit/pkg/validate"
)

type Users struct {
	Phone   string `form:"phone" json:"phone" label:"手机号" validate:"required,max=8,min=6"`
	Passwd  string `form:"passwd" json:"passwd" label:"密码" validate:"required,max=20,min=6"`
	Code    string `form:"code" json:"code" label:"验证码" validate:"required,len=6"`
	Text    string `json:"text" label:"文本" validate:"oneof=red green"`
	Id      string `json:"id" label:"编号" validate:"required,numeric"`
	Confirm string `json:"confirm" label:"校验密码" validate:"eqfield=Passwd"`
	Email   string `json:"email" label:"邮箱" validate:"email"`
	Content string `json:"content" label:"正文" validate:"html"`
}

func main() {
	users := &Users{
		Phone:   "109287222",
		Passwd:  "123098",
		Code:    "123456",
		Text:    "red",
		Confirm: "123098",
		Id:      "12",
		Email:   "abcfoxmail@foxmail.com",
		Content: "<a>aps<a>",
	}
	ctx := context.WithValue(context.Background(), "trace", "ctx-validate-001")
	validate.InitValidate("validator")
	message, err := validate.Check(ctx, users)
	fmt.Println("message: ", message, err) // output: 邮箱的内容必须符合邮箱格式
}
