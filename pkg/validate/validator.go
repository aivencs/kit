package validate

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/go-playground/validator"
)

var once sync.Once
var inspector Inspector

type Inspector interface {
	Check(ctx context.Context, payload interface{}) (string, error)
}

type Validator struct {
	Instance *validator.Validate
}

func InitValidate(name string) {
	once.Do(func() {
		switch name {
		case "validator":
			inspector = NewValidator()
		default:
			inspector = NewValidator()
		}
	})
}

func NewValidator() Inspector {
	v := validator.New()
	// use label instead name
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("label")
		return name
	})
	return &Validator{
		Instance: v,
	}
}

func (c *Validator) Check(ctx context.Context, payload interface{}) (string, error) {
	message := ""
	err := c.Instance.Struct(payload)
	if err == nil {
		return "", err
	}
	for _, err := range err.(validator.ValidationErrors) {

		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("%s为必填项", err.Field())
		case "min":
			message = fmt.Sprintf("%s的长度不应小于%v", err.Field(), err.Param())
		case "max":
			message = fmt.Sprintf("%s的长度不应超过%v", err.Field(), err.Param())
		case "ne":
			message = fmt.Sprintf("%s的值不应为%v", err.Field(), err.Value())
		case "len":
			message = fmt.Sprintf("%s的长度必须为%v", err.Field(), err.Param())
		case "eq":
			message = fmt.Sprintf("%s的值必须为%v", err.Field(), err.Param())
		case "oneof":
			message = fmt.Sprintf("%s的值必须在[%v]其中", err.Field(), err.Param())
		case "gt":
			message = fmt.Sprintf("%s的值必须大于%v", err.Field(), err.Param())
		case "gte":
			message = fmt.Sprintf("%s的值必须大于或等于%v", err.Field(), err.Param())
		case "lt":
			message = fmt.Sprintf("%s的值必须小于%v", err.Field(), err.Param())
		case "lte":
			message = fmt.Sprintf("%s的值必须小于或等于%v", err.Field(), err.Param())
		case "eqfield":
			message = fmt.Sprintf("%s的值必须与%v的值相等", err.Field(), err.Param())
		case "numeric":
			message = fmt.Sprintf("%s的值必须为数字", err.Field())
		case "email":
			message = fmt.Sprintf("%s的值必须符合邮箱格式", err.Field())
		case "url":
			message = fmt.Sprintf("%s的值必须符合网址格式", err.Field())
		case "ip":
			message = fmt.Sprintf("%s的内容必须符合IP格式", err.Field())
		case "contains":
			message = fmt.Sprintf("%s的值必须包含%v", err.Field(), err.Param())
		case "excludes":
			message = fmt.Sprintf("%s的值不可包含%v", err.Field(), err.Param())
		case "containsany":
			message = fmt.Sprintf("%s的值必须包含[%v]其中任意一个", err.Field(), err.Param())
		case "excludesall":
			message = fmt.Sprintf("%s的值不可包含[%v]其中任意一个", err.Field(), err.Param())
		case "startswith":
			message = fmt.Sprintf("%s的值必须以[%v]为开头", err.Field(), err.Param())
		case "endswith":
			message = fmt.Sprintf("%s的值必须以[%v]为结尾", err.Field(), err.Param())
		default:
			message = fmt.Sprintf("%s的值未通过校验", err.Field())
		}
	}
	return message, err
}

func Check(ctx context.Context, payload interface{}) (string, error) {
	return inspector.Check(ctx, payload)
}
