package logger

var erc map[string]Erc
var defaultName = "success"

type Erc struct {
	Name  string
	Code  int
	Level string
	Label string
}

func InitErc() {
	erc = map[string]Erc{
		"success":      {Label: "操作成功", Code: 10000, Level: "info"},
		"check":        {Label: "操作失败", Code: 10001, Level: "check"},
		"over":         {Label: "超限", Code: 10002, Level: "error"},
		"timeout":      {Label: "超时", Code: 10003, Level: "error"},
		"supp":         {Label: "补充数据", Code: 10004, Level: "warn"},
		"abnormal":     {Label: "非常规状态码", Code: 10005, Level: "error"},
		"ede":          {Label: "编码/解码失败", Code: 10006, Level: "error"},
		"rpe":          {Label: "运行时参数错误", Code: 10007, Level: "error"},
		"pve":          {Label: "参数未通过校验", Code: 10008, Level: "error"},
		"dve":          {Label: "数据结果未通过校验", Code: 10009, Level: "error"},
		"rwa":          {Label: "运行时发生异常", Code: 10010, Level: "warn"},
		"rpw":          {Label: "运行时参数错误", Code: 10011, Level: "warn"},
		"timeout-call": {Label: "调用超时", Code: 20001, Level: "check"},
		"error-call":   {Label: "调用错误", Code: 20002, Level: "error"},
		"interrupt":    {Label: "组件中断", Code: 30001, Level: "fatal"},
	}
}

func GetErc(name string) Erc {
	value := erc[name]
	if value.Code == 0 {
		return GetDefaultErc()
	}
	value.Name = name
	return value
}

func GetDefaultErc() Erc {
	value := erc[defaultName]
	value.Name = defaultName
	return value
}
