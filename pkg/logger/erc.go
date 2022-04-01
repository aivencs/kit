package logger

import "unicode/utf8"

var erc map[string]Erc
var defaultName = "SUCCESS"

type Erc struct {
	Name  string
	Code  int
	Level string
	Label string
}

func InitErc() {
	erc = map[string]Erc{
		"SUCCESS":      {Label: "操作成功", Code: 10000, Level: "info"},
		"CHECK":        {Label: "请检查", Code: 10001, Level: "check"},
		"OVERLOAD":     {Label: "超限", Code: 10002, Level: "error"},
		"TIMEOUT":      {Label: "超时", Code: 10003, Level: "error"},
		"SUPP":         {Label: "补充数据", Code: 10004, Level: "warn"},
		"ABNORMAL":     {Label: "非常规状态码", Code: 10005, Level: "error"},
		"EDE":          {Label: "编码/解码失败", Code: 10006, Level: "error"},
		"RPE":          {Label: "运行时参数错误", Code: 10007, Level: "error"},
		"PVE":          {Label: "参数未通过校验", Code: 10008, Level: "error"},
		"DVE":          {Label: "数据结果未通过校验", Code: 10009, Level: "error"},
		"RWA":          {Label: "运行时发生异常", Code: 10010, Level: "warn"},
		"RPW":          {Label: "运行时参数错误", Code: 10011, Level: "warn"},
		"CALL-TIMEOUT": {Label: "调用超时", Code: 20001, Level: "check"},
		"CALL_ERROR":   {Label: "调用错误", Code: 20002, Level: "error"},
		"INTERRUPT":    {Label: "组件中断", Code: 30001, Level: "fatal"},
	}
}

func GetErc(name string, label string) Erc {
	value := erc[name]
	if value.Code == 0 {
		return GetDefaultErc()
	}
	value.Name = name
	if utf8.RuneCountInString(label) > 1 {
		value.Label = label
	}
	return value
}

func GetDefaultErc() Erc {
	value := erc[defaultName]
	value.Name = defaultName
	return value
}
